package config

import (
	"errors"
	"time"
)

// Config holds the configuration for a cronwrap job execution.
type Config struct {
	// Command is the shell command to execute.
	Command string

	// Args are optional arguments passed to the command.
	Args []string

	// Timeout is the maximum duration allowed for the command to run.
	// A zero value means no timeout.
	Timeout time.Duration

	// Retries is the number of times to retry the command on failure.
	Retries int

	// RetryDelay is the duration to wait between retry attempts.
	RetryDelay time.Duration

	// JobName is a human-readable name for the job used in logs and alerts.
	JobName string

	// AlertOnFailure controls whether an alert is sent when the job fails.
	AlertOnFailure bool
}

// Validate checks that the Config has all required fields and valid values.
func (c *Config) Validate() error {
	if c.Command == "" {
		return errors.New("config: command must not be empty")
	}
	if c.Retries < 0 {
		return errors.New("config: retries must be non-negative")
	}
	if c.Timeout < 0 {
		return errors.New("config: timeout must be non-negative")
	}
	if c.RetryDelay < 0 {
		return errors.New("config: retry_delay must be non-negative")
	}
	if c.JobName == "" {
		c.JobName = c.Command
	}
	return nil
}

// Default returns a Config populated with sensible defaults.
func Default(command string) *Config {
	return &Config{
		Command:        command,
		Timeout:        30 * time.Second,
		Retries:        0,
		RetryDelay:     5 * time.Second,
		JobName:        command,
		AlertOnFailure: true,
	}
}
