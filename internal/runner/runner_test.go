package runner

import (
	"strings"
	"testing"
	"time"
)

func TestRun_SuccessOnFirstAttempt(t *testing.T) {
	cfg := Config{
		Command:    "echo",
		Args:       []string{"hello"},
		MaxRetries: 2,
		RetryDelay: 10 * time.Millisecond,
		Timeout:    5 * time.Second,
	}

	result := Run(cfg)

	if result.Err != nil {
		t.Fatalf("expected no error, got: %v", result.Err)
	}
	if result.Attempt != 1 {
		t.Errorf("expected attempt 1, got %d", result.Attempt)
	}
	if !strings.Contains(string(result.Output), "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", result.Output)
	}
}

func TestRun_FailsAllAttempts(t *testing.T) {
	cfg := Config{
		Command:    "false",
		MaxRetries: 2,
		RetryDelay: 10 * time.Millisecond,
		Timeout:    5 * time.Second,
	}

	result := Run(cfg)

	if result.Err == nil {
		t.Fatal("expected an error, got nil")
	}
	if result.Attempt != 3 {
		t.Errorf("expected final attempt to be 3, got %d", result.Attempt)
	}
}

func TestRun_Timeout(t *testing.T) {
	cfg := Config{
		Command:    "sleep",
		Args:       []string{"10"},
		MaxRetries: 0,
		RetryDelay: 0,
		Timeout:    50 * time.Millisecond,
	}

	result := Run(cfg)

	if result.Err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(result.Err.Error(), "timed out") {
		t.Errorf("expected timeout error message, got: %v", result.Err)
	}
}

func TestRun_DurationTracked(t *testing.T) {
	cfg := Config{
		Command:    "echo",
		Args:       []string{"timing test"},
		MaxRetries: 0,
		Timeout:    5 * time.Second,
	}

	result := Run(cfg)

	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %s", result.Duration)
	}
	if result.StartTime.IsZero() {
		t.Error("expected non-zero start time")
	}
}
