package config

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Config holds all runtime configuration for a cronwrap invocation.
type Config struct {
	Command      string
	Args         []string
	JobName      string
	Retries      int
	Timeout      time.Duration
	RetryDelay   time.Duration
	WebhookURL   string
	SlackURL     string
	PagerDutyKey string
	EmailTo      string
	EmailFrom    string
	SMTPAddr     string
	Quiet        bool
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Retries:    0,
		Timeout:    0,
		RetryDelay: 5 * time.Second,
	}
}

// Validate checks the Config for required fields and logical constraints.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Command) == "" {
		return errors.New("config: command is required")
	}
	if c.Retries < 0 {
		return fmt.Errorf("config: retries must be non-negative, got %d", c.Retries)
	}
	if c.Timeout < 0 {
		return fmt.Errorf("config: timeout must be non-negative, got %v", c.Timeout)
	}
	if c.RetryDelay < 0 {
		return fmt.Errorf("config: retry-delay must be non-negative, got %v", c.RetryDelay)
	}
	if c.JobName == "" {
		c.JobName = c.Command
	}
	if c.EmailTo != "" || c.EmailFrom != "" {
		if c.SMTPAddr == "" {
			return errors.New("config: smtp-addr is required when email-to or email-from is set")
		}
	}
	if c.PagerDutyKey != "" && len(c.PagerDutyKey) < 8 {
		return errors.New("config: pagerduty-key appears invalid (too short)")
	}
	return nil
}
