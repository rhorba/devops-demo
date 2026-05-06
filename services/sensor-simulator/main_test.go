package main

import (
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"
)

func TestGetenv(t *testing.T) {
	os.Setenv("SS_TEST_VAR", "hello")
	defer os.Unsetenv("SS_TEST_VAR")

	if got := getenv("SS_TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
	if got := getenv("SS_MISSING_VAR", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestReadingJSONRoundTrip(t *testing.T) {
	r := Reading{
		Zone:      "zone-a",
		Type:      "temperature",
		Value:     23.45,
		Unit:      "celsius",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got Reading
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Zone != r.Zone || got.Type != r.Type || got.Value != r.Value || got.Unit != r.Unit {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, r)
	}
}

func TestSensorsNotEmpty(t *testing.T) {
	if len(sensors) == 0 {
		t.Fatal("sensors slice must not be empty")
	}
}

func TestSensorValueBounds(t *testing.T) {
	tests := []struct{ factor float64 }{
		{0.0}, {0.25}, {0.5}, {0.75}, {1.0},
	}
	for _, s := range sensors {
		for _, tt := range tests {
			v := math.Round((s.min+tt.factor*(s.max-s.min))*100) / 100
			if v < s.min-0.01 || v > s.max+0.01 {
				t.Errorf("sensor %s/%s: value %f out of [%f, %f]",
					s.zone, s.stype, v, s.min, s.max)
			}
		}
	}
}

func TestSensorFields(t *testing.T) {
	for _, s := range sensors {
		if s.zone == "" {
			t.Error("sensor zone must not be empty")
		}
		if s.stype == "" {
			t.Error("sensor type must not be empty")
		}
		if s.unit == "" {
			t.Error("sensor unit must not be empty")
		}
		if s.min >= s.max {
			t.Errorf("sensor %s/%s: min (%f) must be less than max (%f)",
				s.zone, s.stype, s.min, s.max)
		}
	}
}
