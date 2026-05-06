package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Reading struct {
	Zone      string    `json:"zone"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	published = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sensor_messages_published_total",
		Help: "Total MQTT messages published per zone and type",
	}, []string{"zone", "type"})

	publishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sensor_publish_errors_total",
		Help: "Total MQTT publish errors per zone",
	}, []string{"zone"})

	publishLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "sensor_publish_duration_seconds",
		Help:    "MQTT publish round-trip latency",
		Buckets: prometheus.DefBuckets,
	})
)

type sensor struct {
	zone  string
	stype string
	unit  string
	min   float64
	max   float64
}

var sensors = []sensor{
	{"zone-a", "temperature", "celsius", 18.0, 35.0},
	{"zone-a", "pressure", "bar", 0.8, 1.2},
	{"zone-b", "humidity", "percent", 30.0, 90.0},
	{"zone-b", "temperature", "celsius", 20.0, 40.0},
	{"zone-c", "vibration", "mm/s", 0.0, 5.0},
	{"zone-c", "pressure", "bar", 0.9, 1.5},
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newMQTTClient(broker, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetCleanSession(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetOnConnectHandler(func(_ mqtt.Client) {
			log.Printf("Connected to MQTT broker: %s", broker)
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			log.Printf("MQTT connection lost: %v — reconnecting...", err)
		})
	return mqtt.NewClient(opts)
}

func main() {
	broker := getenv("MQTT_BROKER", "tcp://mosquitto:1883")
	metricsAddr := getenv("METRICS_ADDR", ":8080")
	intervalStr := getenv("INTERVAL_MS", "1000")

	interval, _ := time.ParseDuration(intervalStr + "ms")
	if interval == 0 {
		interval = time.Second
	}

	client := newMQTTClient(broker, "sensor-simulator")
	if tok := client.Connect(); tok.Wait() && tok.Error() != nil {
		log.Fatalf("MQTT connect failed: %v", tok.Error())
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"ok","broker":%q}`, broker)
		})
		log.Printf("Metrics server on %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, mux); err != nil {
			log.Fatalf("metrics server: %v", err)
		}
	}()

	log.Printf("Starting sensor simulation — interval: %v, sensors: %d", interval, len(sensors))
	ticker := time.NewTicker(interval)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range ticker.C {
		for _, s := range sensors {
			r := Reading{
				Zone:      s.zone,
				Type:      s.stype,
				Value:     math.Round((s.min+rng.Float64()*(s.max-s.min))*100) / 100,
				Unit:      s.unit,
				Timestamp: time.Now().UTC(),
			}
			payload, _ := json.Marshal(r)
			topic := fmt.Sprintf("sensors/%s/%s", s.zone, s.stype)

			start := time.Now()
			tok := client.Publish(topic, 1, false, payload)
			tok.Wait()
			publishLatency.Observe(time.Since(start).Seconds())

			if err := tok.Error(); err != nil {
				log.Printf("publish error [%s]: %v", topic, err)
				publishErrors.WithLabelValues(s.zone).Inc()
			} else {
				published.WithLabelValues(s.zone, s.stype).Inc()
			}
		}
	}
}
