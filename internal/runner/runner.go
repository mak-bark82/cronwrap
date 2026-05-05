package runner

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Config holds the configuration for running a cron job command.
type Config struct {
	Command    string
	Args       []string
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// Result captures the outcome of a command execution attempt.
type Result struct {
	Attempt   int
	Output    []byte
	Err       error
	StartTime time.Time
	Duration  time.Duration
}

// Run executes the configured command with retry logic and timeout support.
func Run(cfg Config) Result {
	var result Result

	for attempt := 1; attempt <= cfg.MaxRetries+1; attempt++ {
		result = execute(cfg, attempt)
		if result.Err == nil {
			log.Printf("[cronwrap] command succeeded on attempt %d (duration: %s)", attempt, result.Duration)
			return result
		}

		log.Printf("[cronwrap] attempt %d failed: %v", attempt, result.Err)

		if attempt <= cfg.MaxRetries {
			log.Printf("[cronwrap] retrying in %s...", cfg.RetryDelay)
			time.Sleep(cfg.RetryDelay)
		}
	}

	return result
}

func execute(cfg Config, attempt int) Result {
	ctx := context.Background()
	var cancel context.CancelFunc

	if cfg.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, cfg.Command, cfg.Args...)
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("command timed out after %s", cfg.Timeout)
		}
	}

	return Result{
		Attempt:   attempt,
		Output:    output,
		Err:       err,
		StartTime: start,
		Duration:  duration,
	}
}
