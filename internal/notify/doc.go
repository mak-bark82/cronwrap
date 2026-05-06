// Package notify provides pluggable notifiers for cronwrap.
//
// Each notifier implements a simple Notify(jobName string, err error) error
// interface so that multiple backends can be composed together via
// MultiNotifier.
//
// Available notifiers:
//
//   - StderrNotifier  – writes alerts to standard error (zero config)
//   - WebhookNotifier – HTTP POST to an arbitrary webhook endpoint
//   - EmailNotifier   – sends SMTP email alerts
//   - SlackNotifier   – posts to a Slack incoming-webhook URL
//   - MultiNotifier   – fans out to any number of notifiers
//
// Example usage:
//
//	slack := notify.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK"))
//	email := notify.NewEmailNotifier("smtp.example.com:587", "from@example.com", "to@example.com")
//	multi := notify.NewMultiNotifier(slack, email)
//	_ = multi.Notify("my-job", err)
package notify
