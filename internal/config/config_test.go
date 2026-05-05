package config

import (
	"testing"
	"time"
)

func TestValidate_MissingCommand(t *testing.T) {
	c := &Config{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestValidate_NegativeRetries(t *testing.T) {
	c := &Config{Command: "echo", Retries: -1}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative retries, got nil")
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	c := &Config{Command: "echo", Timeout: -1 * time.Second}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
}

func TestValidate_NegativeRetryDelay(t *testing.T) {
	c := &Config{Command: "echo", RetryDelay: -1 * time.Second}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative retry delay, got nil")
	}
}

func TestValidate_SetsJobNameFromCommand(t *testing.T) {
	c := &Config{Command: "mycommand"}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.JobName != "mycommand" {
		t.Errorf("expected JobName %q, got %q", "mycommand", c.JobName)
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	c := &Config{
		Command:    "echo",
		Args:       []string{"hello"},
		Timeout:    10 * time.Second,
		Retries:    3,
		RetryDelay: 2 * time.Second,
		JobName:    "test-job",
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error for valid config: %v", err)
	}
}

func TestDefault_SetsExpectedValues(t *testing.T) {
	c := Default("backup.sh")

	if c.Command != "backup.sh" {
		t.Errorf("expected Command %q, got %q", "backup.sh", c.Command)
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("expected Timeout 30s, got %v", c.Timeout)
	}
	if c.Retries != 0 {
		t.Errorf("expected Retries 0, got %d", c.Retries)
	}
	if c.RetryDelay != 5*time.Second {
		t.Errorf("expected RetryDelay 5s, got %v", c.RetryDelay)
	}
	if !c.AlertOnFailure {
		t.Error("expected AlertOnFailure to be true")
	}
	if c.JobName != "backup.sh" {
		t.Errorf("expected JobName %q, got %q", "backup.sh", c.JobName)
	}
}
