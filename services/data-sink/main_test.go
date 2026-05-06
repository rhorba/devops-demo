package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetenv(t *testing.T) {
	os.Setenv("DS_TEST_VAR", "hello")
	defer os.Unsetenv("DS_TEST_VAR")

	if got := getenv("DS_TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
	if got := getenv("DS_MISSING_VAR", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestS3KeyFormat(t *testing.T) {
	t0 := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	key := s3Key(t0)

	if !strings.HasPrefix(key, "2024-03-15/14/") {
		t.Errorf("s3Key = %q, want prefix 2024-03-15/14/", key)
	}
	if !strings.HasSuffix(key, ".json") {
		t.Errorf("s3Key = %q, want suffix .json", key)
	}
}

func TestS3KeyMidnightHour(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	key := s3Key(t0)

	if !strings.HasPrefix(key, "2024-01-01/00/") {
		t.Errorf("midnight s3Key = %q, want prefix 2024-01-01/00/", key)
	}
}

func TestS3KeyUnique(t *testing.T) {
	t0 := time.Now()
	k1 := s3Key(t0)
	k2 := s3Key(t0)
	if k1 == k2 {
		t.Error("s3Key must return a unique key on each call (UUID-based)")
	}
}

func TestS3KeySegments(t *testing.T) {
	t0 := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	key := s3Key(t0)
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		t.Errorf("expected 3 path segments in s3Key, got %d: %q", len(parts), key)
	}
}
