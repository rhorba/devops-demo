package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestGetenv(t *testing.T) {
	os.Setenv("EC_TEST_VAR", "hello")
	defer os.Unsetenv("EC_TEST_VAR")

	if got := getenv("EC_TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
	if got := getenv("EC_MISSING_VAR", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestEnvelopeJSONStructure(t *testing.T) {
	payload := []byte(`{"zone":"zone-a","type":"temperature","value":23.5}`)
	now := time.Now().UTC()

	envelope := map[string]interface{}{
		"source":      "ibm-mq",
		"qmgr":        "QM1",
		"queue":       "DEV.QUEUE.1",
		"payload":     string(payload),
		"received_at": now,
		"protocol":    "amqp-1.0",
	}

	b, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	checks := map[string]string{
		"source":   "ibm-mq",
		"qmgr":     "QM1",
		"queue":    "DEV.QUEUE.1",
		"protocol": "amqp-1.0",
	}
	for field, want := range checks {
		if got[field] != want {
			t.Errorf("envelope[%q] = %v, want %q", field, got[field], want)
		}
	}

	if _, ok := got["received_at"]; !ok {
		t.Error("envelope must contain received_at field")
	}
	if _, ok := got["payload"]; !ok {
		t.Error("envelope must contain payload field")
	}
}

func TestEnvelopePayloadTypes(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
	}{
		{"json_object", []byte(`{"key":"value"}`)},
		{"plain_string", []byte(`hello world`)},
		{"empty", []byte(``)},
		{"binary_like", []byte(`\x00\x01\x02`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envelope := map[string]interface{}{
				"source":  "ibm-mq",
				"payload": string(tt.payload),
			}
			if _, err := json.Marshal(envelope); err != nil {
				t.Errorf("marshal failed for payload type %q: %v", tt.name, err)
			}
		})
	}
}
