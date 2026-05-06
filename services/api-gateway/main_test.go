package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGetenv(t *testing.T) {
	os.Setenv("GW_TEST_VAR", "hello")
	defer os.Unsetenv("GW_TEST_VAR")

	if got := getenv("GW_TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
	if got := getenv("GW_MISSING_VAR", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestGatewayEmptyLatest(t *testing.T) {
	gw := newGateway()
	msgs := gw.latest(10)
	if len(msgs) != 0 {
		t.Errorf("expected 0 messages on empty gateway, got %d", len(msgs))
	}
}

func TestGatewayPushAndLatest(t *testing.T) {
	gw := newGateway()
	for i := 0; i < 5; i++ {
		gw.push(json.RawMessage(`{}`))
	}
	got := gw.latest(3)
	if len(got) != 3 {
		t.Errorf("expected 3 messages, got %d", len(got))
	}
}

func TestGatewayLatestCappedAtCount(t *testing.T) {
	gw := newGateway()
	gw.push(json.RawMessage(`{}`))
	gw.push(json.RawMessage(`{}`))

	got := gw.latest(100)
	if len(got) != 2 {
		t.Errorf("latest(100) with 2 messages: expected 2, got %d", len(got))
	}
}

func TestGatewayRingOverflow(t *testing.T) {
	gw := newGateway()
	for i := 0; i < ringSize+10; i++ {
		gw.push(json.RawMessage(`{}`))
	}
	got := gw.latest(ringSize)
	if len(got) != ringSize {
		t.Errorf("after %d pushes, latest(%d) returned %d (want %d)",
			ringSize+10, ringSize, len(got), ringSize)
	}
}

func TestGatewaySSEClientTracking(t *testing.T) {
	gw := newGateway()
	ch := make(chan json.RawMessage, 1)

	gw.addSSEClient(ch)
	if len(gw.clients) != 1 {
		t.Errorf("expected 1 SSE client after add, got %d", len(gw.clients))
	}

	gw.removeSSEClient(ch)
	if len(gw.clients) != 0 {
		t.Errorf("expected 0 SSE clients after remove, got %d", len(gw.clients))
	}
}

func TestGatewayPushNotifiesSSEClient(t *testing.T) {
	gw := newGateway()
	ch := make(chan json.RawMessage, 1)
	gw.addSSEClient(ch)

	msg := json.RawMessage(`{"test":true}`)
	gw.push(msg)

	select {
	case got := <-ch:
		if string(got) != string(msg) {
			t.Errorf("SSE received %s, want %s", got, msg)
		}
	default:
		t.Error("SSE client was not notified by push")
	}
}
