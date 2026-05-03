package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	recordsSinked = promauto.NewCounter(prometheus.CounterOpts{
		Name: "datasink_records_sinked_total",
		Help: "Total records successfully written to S3",
	})

	s3Errors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "datasink_s3_errors_total",
		Help: "Total S3 put errors",
	})

	kafkaLag = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "datasink_kafka_consumer_lag",
		Help: "Estimated Kafka consumer lag",
	})

	batchSize = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "datasink_batch_size",
		Help:    "Number of records per S3 batch",
		Buckets: []float64{1, 5, 10, 25, 50, 100},
	})
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newS3Client(ctx context.Context, endpoint, region string) *s3.Client {
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
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func s3Key(t time.Time) string {
	return fmt.Sprintf("%s/%02d/%s.json",
		t.Format("2006-01-02"),
		t.Hour(),
		uuid.New().String(),
	)
}

func main() {
	brokers := []string{getenv("KAFKA_BROKERS", "redpanda:9092")}
	topic := getenv("KAFKA_TOPIC", "raw.telemetry")
	group := getenv("KAFKA_GROUP", "data-sink")
	bucket := getenv("S3_BUCKET", "telemetry-archive")
	s3Endpoint := getenv("AWS_ENDPOINT_URL", "http://localstack:4566")
	s3Region := getenv("AWS_REGION", "us-east-1")
	metricsAddr := getenv("METRICS_ADDR", ":8080")
	flushInterval := 5 * time.Second

	ctx := context.Background()
	s3Client := newS3Client(ctx, s3Endpoint, s3Region)

	// Ensure bucket exists
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		log.Printf("Bucket %s may already exist: %v", bucket, err)
	}

	// Kafka consumer via franz-go (pure Go, no CGO)
	kgoClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup(group),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.WithLogger(kgo.BasicLogger(log.Writer(), kgo.LogLevelInfo, nil)),
	)
	if err != nil {
		log.Fatalf("Kafka client: %v", err)
	}
	defer kgoClient.Close()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"ok","bucket":%q,"topic":%q}`, bucket, topic)
		})
		log.Printf("Metrics on %s", metricsAddr)
		log.Fatal(http.ListenAndServe(metricsAddr, mux))
	}()

	log.Printf("data-sink starting: brokers=%v topic=%s group=%s bucket=%s", brokers, topic, group, bucket)

	var batch []json.RawMessage
	var mu sync.Mutex
	flushTicker := time.NewTicker(flushInterval)

	flush := func() {
		mu.Lock()
		defer mu.Unlock()
		if len(batch) == 0 {
			return
		}
		body, err := json.Marshal(batch)
		if err != nil {
			log.Printf("json marshal: %v", err)
			return
		}
		key := s3Key(time.Now().UTC())
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(body),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			log.Printf("S3 put error: %v", err)
			s3Errors.Inc()
			return
		}
		batchSize.Observe(float64(len(batch)))
		recordsSinked.Add(float64(len(batch)))
		log.Printf("Flushed %d records → s3://%s/%s", len(batch), bucket, key)
		batch = batch[:0]
	}

	// Flush on schedule regardless of message flow.
	go func() {
		for range flushTicker.C {
			flush()
		}
	}()

	for {
		fetches := kgoClient.PollFetches(ctx)
		if fetches.IsClientClosed() {
			flush()
			return
		}
		fetches.EachError(func(t string, p int32, err error) {
			log.Printf("Fetch error topic=%s partition=%d: %v", t, p, err)
		})
		mu.Lock()
		fetches.EachRecord(func(r *kgo.Record) {
			batch = append(batch, json.RawMessage(r.Value))
			kafkaLag.Set(float64(r.Offset))
		})
		shouldFlush := len(batch) >= 50
		mu.Unlock()
		if shouldFlush {
			flush()
		}
	}
}
