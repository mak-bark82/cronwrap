// Package main is the entry point for the cronwrap command-line tool.
// It wires together configuration, logging, alerting, retry logic, and metrics
// to provide a robust wrapper around arbitrary cron job commands.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/cronwrap/internal/alert"
	"github.com/yourorg/cronwrap/internal/config"
	"github.com/yourorg/cronwrap/internal/logger"
	"github.com/yourorg/cronwrap/internal/metrics"
	"github.com/yourorg/cronwrap/internal/notify"
	"github.com/yourorg/cronwrap/internal/runner"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cronwrap: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse flags and build configuration.
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Set up structured logger.
	log := logger.New(os.Stdout, cfg.JobName)

	// Build the notifier chain based on configured destinations.
	notifier := buildNotifier(cfg)

	// Set up alerter that wraps the notifier.
	alerter := alert.New(notifier, log)

	// Metrics collector for this run.
	collector := metrics.NewCollector()

	// Honour OS signals for graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Execute the job via the runner (handles retries, timeout, logging).
	result := runner.Run(ctx, cfg, log, alerter, collector)

	// Print a summary to stdout so it appears in cron mail / logs.
	metrics.PrintSummary(os.Stdout, collector.Summary())
	metrics.PrintResults(os.Stdout, collector.All())

	if !result.Succeeded {
		return fmt.Errorf("job %q failed after %d attempt(s): %w",
			cfg.JobName, result.Attempts, result.Err)
	}

	return nil
}

// buildNotifier constructs a MultiNotifier from whichever notification
// back-ends have been enabled in the configuration.
func buildNotifier(cfg *config.Config) notify.Notifier {
	multi := notify.NewMultiNotifier()

	// Always include the stderr notifier so failures are visible in cron logs.
	multi.Add(alert.NewStderrNotifier())

	if cfg.SlackWebhook != "" {
		multi.Add(notify.NewSlackNotifier(cfg.SlackWebhook))
	}

	if cfg.WebhookURL != "" {
		multi.Add(notify.NewWebhookNotifier(cfg.WebhookURL))
	}

	if cfg.PagerDutyKey != "" {
		multi.Add(notify.NewPagerDutyNotifier(cfg.PagerDutyKey))
	}

	if cfg.OpsGenieKey != "" {
		multi.Add(notify.NewOpsGenieNotifier(cfg.OpsGenieKey))
	}

	if cfg.DatadogAPIKey != "" {
		multi.Add(notify.NewDatadogNotifier(cfg.DatadogAPIKey))
	}

	if cfg.EmailTo != "" && cfg.SMTPHost != "" {
		multi.Add(notify.NewEmailNotifier(
			cfg.SMTPHost, cfg.SMTPPort,
			cfg.EmailFrom, cfg.EmailTo,
			cfg.SMTPUsername, cfg.SMTPPassword,
		))
	}

	return multi
}
