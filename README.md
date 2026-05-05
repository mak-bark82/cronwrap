# cronwrap

Lightweight wrapper for cron jobs that adds logging, alerting, and retry logic.

## Installation

```bash
go install github.com/yourusername/cronwrap@latest
```

Or add it as a dependency:

```bash
go get github.com/yourusername/cronwrap
```

## Usage

Wrap any command to get automatic logging, failure alerts, and retries:

```bash
cronwrap --retries 3 --alert-on-failure --job-name "db-backup" -- /usr/local/bin/backup.sh
```

Or use it programmatically in your Go code:

```go
package main

import "github.com/yourusername/cronwrap"

func main() {
    job := cronwrap.NewJob("db-backup", func() error {
        return runBackup()
    })

    job.WithRetries(3).
        WithLogger(cronwrap.DefaultLogger).
        WithAlert(cronwrap.SlackAlert("https://hooks.slack.com/...")).
        Run()
}
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--retries` | Number of retry attempts on failure | `0` |
| `--job-name` | Identifier used in logs and alerts | `""` |
| `--alert-on-failure` | Send alert if job fails after all retries | `false` |
| `--timeout` | Maximum job duration before killing | `0` (none) |

## Features

- **Structured logging** — JSON or text output for every job run
- **Retry logic** — configurable backoff and retry count
- **Alerting** — pluggable alert backends (Slack, PagerDuty, email)
- **Timeout support** — kill long-running jobs automatically
- **Exit code passthrough** — preserves original command exit codes

## License

MIT © [yourusername](https://github.com/yourusername)