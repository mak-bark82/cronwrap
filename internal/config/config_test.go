package config

import (
	"strings"
	"testing"
	"time"
)

func TestValidate_MissingCommand(t *testing.T) {
	c := Default()
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "command") {
		t.Errorf("expected command error, got %v", err)
	}
}

func TestValidate_NegativeRetries(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.Retries = -1
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "retries") {
		t.Errorf("expected retries error, got %v", err)
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.Timeout = -1 * time.Second
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestValidate_NegativeRetryDelay(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.RetryDelay = -1 * time.Second
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "retry-delay") {
		t.Errorf("expected retry-delay error, got %v", err)
	}
}

func TestValidate_SetsJobNameFromCommand(t *testing.T) {
	c := Default()
	c.Command = "mycommand"
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.JobName != "mycommand" {
		t.Errorf("expected JobName %q, got %q", "mycommand", c.JobName)
	}
}

func TestValidate_EmailRequiresSMTP(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.EmailTo = "ops@example.com"
	c.EmailFrom = "cronwrap@example.com"
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "smtp-addr") {
		t.Errorf("expected smtp-addr error, got %v", err)
	}
}

func TestValidate_EmailWithSMTP_OK(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.EmailTo = "ops@example.com"
	c.EmailFrom = "cronwrap@example.com"
	c.SMTPAddr = "localhost:25"
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_PagerDutyKeyTooShort(t *testing.T) {
	c := Default()
	c.Command = "echo"
	c.PagerDutyKey = "short"
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "pagerduty-key") {
		t.Errorf("expected pagerduty-key error, got %v", err)
	}
}

func TestDefault_ReturnsDefaults(t *testing.T) {
	c := Default()
	if c.Retries != 0 {
		t.Errorf("expected Retries=0, got %d", c.Retries)
	}
	if c.RetryDelay != 5*time.Second {
		t.Errorf("expected RetryDelay=5s, got %v", c.RetryDelay)
	}
}
