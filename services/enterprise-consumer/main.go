package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/go-amqp"
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
		Help: "Total messages consumed from IBM MQ via AMQP 1.0",
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
	mqHost := getenv("MQ_HOST", "ibm-mq")
	mqAMQPPort := getenv("MQ_AMQP_PORT", "5672")
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
			fmt.Fprintf(w, `{"status":"ok","qmgr":%q,"queue":%q,"protocol":"amqp-1.0"}`, qmgrName, queueName)
		})
		log.Printf("Metrics on %s", metricsAddr)
		log.Fatal(http.ListenAndServe(metricsAddr, mux))
	}()

	addr := fmt.Sprintf("amqp://%s:%s", mqHost, mqAMQPPort)
	log.Printf("Connecting to IBM MQ via AMQP 1.0 at %s queue=%s", addr, queueName)

	for {
		if err := consume(ctx, addr, appUser, appPassword, queueName, qmgrName, sqsClient, sqsQueueURL); err != nil {
			log.Printf("AMQP session ended: %v — reconnecting in 5s", err)
			consumerErrors.WithLabelValues("amqp_connect").Inc()
			time.Sleep(5 * time.Second)
		}
	}
}

func consume(ctx context.Context, addr, user, password, queueName, qmgrName string, sqsClient *sqs.Client, sqsQueueURL string) error {
	opts := &amqp.ConnOptions{}
	if user != "" {
		opts.SASLType = amqp.SASLTypePlain(user, password)
	}

	conn, err := amqp.Dial(ctx, addr, opts)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer conn.Close()

	session, err := conn.NewSession(ctx, nil)
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	receiver, err := session.NewReceiver(ctx, queueName, nil)
	if err != nil {
		return fmt.Errorf("new receiver for %s: %w", queueName, err)
	}
	defer receiver.Close(ctx)

	log.Printf("AMQP receiver ready — queue=%s qmgr=%s", queueName, qmgrName)

	for {
		msg, err := receiver.Receive(ctx, nil)
		if err != nil {
			return fmt.Errorf("receive: %w", err)
		}

		start := time.Now()
		mqConsumed.Inc()

		var payload []byte
		switch v := msg.Value.(type) {
		case []byte:
			payload = v
		case string:
			payload = []byte(v)
		default:
			payload, _ = json.Marshal(msg.Value)
		}

		log.Printf("Message received (%d bytes): %.120s", len(payload), payload)

		envelope := map[string]interface{}{
			"source":      "ibm-mq",
			"qmgr":        qmgrName,
			"queue":       queueName,
			"payload":     string(payload),
			"received_at": start.UTC(),
			"protocol":    "amqp-1.0",
		}
		body, _ := json.Marshal(envelope)

		_, sqsErr := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(sqsQueueURL),
			MessageBody: aws.String(string(body)),
		})
		if sqsErr != nil {
			log.Printf("SQS send error: %v", sqsErr)
			consumerErrors.WithLabelValues("sqs_send").Inc()
			if err := receiver.RejectMessage(ctx, msg, nil); err != nil {
				return fmt.Errorf("reject: %w", err)
			}
			continue
		}

		if err := receiver.AcceptMessage(ctx, msg); err != nil {
			return fmt.Errorf("accept: %w", err)
		}
		sqsSent.Inc()
		processLatency.Observe(time.Since(start).Seconds())
		log.Printf("Forwarded to SQS in %.3fms", time.Since(start).Seconds()*1000)
	}
}
