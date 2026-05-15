// Package notify provides notifiers for alerting on cron job outcomes.
//
// # New Relic
//
// NewRelicNotifier sends structured log events to the New Relic Logs API
// (https://log-api.newrelic.com/log/v1). Each notification includes the job
// name, a human-readable message, and a status field set to either "success"
// or "failure".
//
// Usage:
//
//	n := notify.NewNewRelicNotifier("YOUR_API_KEY")
//	n.Notify("my-job", nil)          // success
//	n.Notify("my-job", someErr)      // failure
package notify
