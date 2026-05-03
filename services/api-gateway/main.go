package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ringSize = 100

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "apigateway_http_requests_total",
		Help: "Total HTTP requests by method and path",
	}, []string{"method", "path", "status"})

	natsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigateway_nats_messages_received_total",
		Help: "Total messages received from NATS JetStream",
	})

	sseClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "apigateway_sse_clients",
		Help: "Current number of SSE clients connected",
	})
)

// Gateway holds the NATS connection, ring buffer, and SSE client registry.
type Gateway struct {
	mu      sync.RWMutex
	ring    []json.RawMessage
	head    int
	count   int
	clients map[chan json.RawMessage]struct{}
}

func newGateway() *Gateway {
	return &Gateway{
		ring:    make([]json.RawMessage, ringSize),
		clients: make(map[chan json.RawMessage]struct{}),
	}
}

func (g *Gateway) push(msg json.RawMessage) {
	g.mu.Lock()
	g.ring[g.head] = msg
	g.head = (g.head + 1) % ringSize
	if g.count < ringSize {
		g.count++
	}
	clients := make([]chan json.RawMessage, 0, len(g.clients))
	for ch := range g.clients {
		clients = append(clients, ch)
	}
	g.mu.Unlock()

	for _, ch := range clients {
		select {
		case ch <- msg:
		default: // slow client — drop
		}
	}
}

func (g *Gateway) latest(n int) []json.RawMessage {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if n > g.count {
		n = g.count
	}
	result := make([]json.RawMessage, n)
	for i := 0; i < n; i++ {
		idx := (g.head - n + i + ringSize) % ringSize
		result[i] = g.ring[idx]
	}
	return result
}

func (g *Gateway) addSSEClient(ch chan json.RawMessage) {
	g.mu.Lock()
	g.clients[ch] = struct{}{}
	g.mu.Unlock()
	sseClients.Inc()
}

func (g *Gateway) removeSSEClient(ch chan json.RawMessage) {
	g.mu.Lock()
	delete(g.clients, ch)
	g.mu.Unlock()
	sseClients.Dec()
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	natsURL := getenv("NATS_URL", "nats://nats:4222")
	listenAddr := getenv("LISTEN_ADDR", ":8080")
	streamName := getenv("NATS_STREAM", "TELEMETRY")
	subject := getenv("NATS_SUBJECT", "telemetry.>")

	gw := newGateway()

	// Connect to NATS with JetStream
	nc, err := nats.Connect(natsURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Println("NATS reconnected")
		}),
	)
	if err != nil {
		log.Fatalf("NATS connect: %v", err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("JetStream context: %v", err)
	}

	ctx := context.Background()

	// Ensure the stream exists (idempotent)
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		MaxAge:   24 * time.Hour,
		Storage:  jetstream.MemoryStorage,
	})
	if err != nil {
		log.Printf("Stream create/update: %v (may already exist)", err)
	}

	cons, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:          "api-gateway",
		FilterSubject: subject,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("Consumer create: %v", err)
	}

	// Start consuming in background
	go func() {
		_, err := cons.Consume(func(msg jetstream.Msg) {
			msg.Ack()
			natsReceived.Inc()
			gw.push(json.RawMessage(msg.Data()))
		})
		if err != nil {
			log.Fatalf("Consume: %v", err)
		}
		log.Printf("Consuming from NATS stream %s, subject %s", streamName, subject)
	}()

	// HTTP handlers
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "ok"
		statusCode := http.StatusOK
		if !nc.IsConnected() {
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
			w.WriteHeader(statusCode)
		}
		fmt.Fprintf(w, `{"status":%q,"nats":%q}`, status, nc.ConnectedUrl())
		httpRequests.WithLabelValues(r.Method, "/health", fmt.Sprint(statusCode)).Inc()
	})

	mux.HandleFunc("/metrics/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		msgs := gw.latest(10)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":    len(msgs),
			"readings": msgs,
		})
		httpRequests.WithLabelValues(r.Method, "/metrics/latest", "200").Inc()
	})

	// Server-Sent Events — live telemetry stream
	mux.HandleFunc("/metrics/stream", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ch := make(chan json.RawMessage, 16)
		gw.addSSEClient(ch)
		defer gw.removeSSEClient(ch)

		httpRequests.WithLabelValues(r.Method, "/metrics/stream", "200").Inc()

		for {
			select {
			case msg := <-ch:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	log.Printf("API Gateway listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("HTTP server: %v", err)
	}
}
