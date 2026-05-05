package config

import (
	"flag"
	"time"
)

// ParseFlags parses command-line flags and returns a populated Config.
// The remaining non-flag arguments are treated as the command and its args.
func ParseFlags(args []string) (*Config, []string, error) {
	fs := flag.NewFlagSet("cronwrap", flag.ContinueOnError)

	jobName := fs.String("job", "", "human-readable name for the job")
	timeout := fs.Duration("timeout", 30*time.Second, "maximum duration for the command")
	retries := fs.Int("retries", 0, "number of retry attempts on failure")
	retryDelay := fs.Duration("retry-delay", 5*time.Second, "delay between retries")
	alertOnFailure := fs.Bool("alert", true, "send alert on job failure")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	remaining := fs.Args()

	cmd := ""
	var cmdArgs []string
	if len(remaining) > 0 {
		cmd = remaining[0]
		cmdArgs = remaining[1:]
	}

	name := *jobName
	if name == "" {
		name = cmd
	}

	cfg := &Config{
		Command:        cmd,
		Args:           cmdArgs,
		Timeout:        *timeout,
		Retries:        *retries,
		RetryDelay:     *retryDelay,
		JobName:        name,
		AlertOnFailure: *alertOnFailure,
	}

	if err := cfg.Validate(); err != nil {
		return nil, remaining, err
	}

	return cfg, remaining, nil
}
