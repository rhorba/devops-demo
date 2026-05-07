package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mqConsumed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "enterprise_mq_messages_consumed_total",
		Help: "Total messages consumed from IBM MQ via REST API",
	})
	sqsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "enterprise_sqs_messages_sent_total",
		Help: "Total messages forwarded to SQS",
	})
	consumerErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "enterprise_errors_total",
		Help: "Total errors by type",
	}, []string{"type"})
	processLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "enterprise_process_duration_seconds",
		Help:    "MQ receive → SQS send latency",
		Buckets: prometheus.DefBuckets,
	})
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newSQSClient(ctx context.Context, endpoint, region string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, reg string, opts ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("AWS config: %v", err)
	}
	return sqs.NewFromConfig(cfg)
}

func main() {
	qmgrName := getenv("MQ_QMGR", "QM1")
	queueName := getenv("MQ_QUEUE", "DEV.QUEUE.1")
	mqHost := getenv("MQ_HOST", "ibm-mq.messaging.svc.cluster.local")
	mqPort := getenv("MQ_REST_PORT", "9443")
	appUser := getenv("MQ_APP_USER", "app")
	appPassword := getenv("MQ_APP_PASSWORD", "")
	sqsQueueURL := getenv("SQS_QUEUE_URL", "http://localstack:4566/000000000000/enterprise-events")
	awsEndpoint := getenv("AWS_ENDPOINT_URL", "http://localstack:4566")
	awsRegion := getenv("AWS_REGION", "us-east-1")
	metricsAddr := getenv("METRICS_ADDR", ":8080")

	ctx := context.Background()
	sqsClient := newSQSClient(ctx, awsEndpoint, awsRegion)

	_, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String("enterprise-events"),
	})
	if err != nil {
		log.Printf("SQS queue may already exist: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"ok","qmgr":%q,"queue":%q,"protocol":"rest-api"}`, qmgrName, queueName)
		})
		log.Printf("Metrics on %s", metricsAddr)
		log.Fatal(http.ListenAndServe(metricsAddr, mux))
	}()

	mqURL := fmt.Sprintf("https://%s:%s/ibmmq/rest/v2/messaging/qmgr/%s/queue/%s/message", mqHost, mqPort, qmgrName, queueName)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(appUser + ":" + appPassword))
	log.Printf("Polling IBM MQ REST API at %s", mqURL)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for {
		if err := poll(ctx, httpClient, mqURL, basicAuth, qmgrName, queueName, sqsClient, sqsQueueURL); err != nil {
			log.Printf("poll error: %v — retrying in 2s", err)
			consumerErrors.WithLabelValues("rest_poll").Inc()
			time.Sleep(2 * time.Second)
		}
	}
}

func poll(ctx context.Context, client *http.Client, mqURL, basicAuth, qmgrName, queueName string, sqsClient *sqs.Client, sqsQueueURL string) error {
	// Destructive GET (DELETE verb in IBM MQ REST API v2)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, mqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+basicAuth)
	req.Header.Set("ibm-mq-rest-csrf-token", "x")
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		// 204 = queue empty, normal — poll again shortly
		time.Sleep(500 * time.Millisecond)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	mqConsumed.Inc()
	log.Printf("Message received (%d bytes): %.120s", len(payload), payload)

	envelope := map[string]interface{}{
		"source":      "ibm-mq",
		"qmgr":        qmgrName,
		"queue":       queueName,
		"payload":     string(payload),
		"received_at": start.UTC(),
		"protocol":    "rest-api",
	}
	body, _ := json.Marshal(envelope)

	_, sqsErr := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(sqsQueueURL),
		MessageBody: aws.String(string(body)),
	})
	if sqsErr != nil {
		log.Printf("SQS send error: %v", sqsErr)
		consumerErrors.WithLabelValues("sqs_send").Inc()
		return nil
	}

	sqsSent.Inc()
	processLatency.Observe(time.Since(start).Seconds())
	log.Printf("Forwarded to SQS in %.3fms", time.Since(start).Seconds()*1000)
	return nil
}
