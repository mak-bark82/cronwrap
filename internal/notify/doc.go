// Package notify provides outbound notification integrations for cronwrap.
//
// Currently supported notifiers:
//
//   - WebhookNotifier: sends a JSON POST request to a configurable HTTP endpoint
//     whenever a job completes, succeeds, or fails.
//
// Example usage:
//
//	n := notify.NewWebhookNotifier("https://hooks.example.com/cronwrap")
//	err := n.Notify(notify.WebhookPayload{
//		JobName:  "db-backup",
//		Status:   "failure",
//		Error:    "exit status 1",
//		Duration: 3 * time.Second,
//		Attempts: 3,
//	})
//
// Additional notifiers (Slack, PagerDuty, email) can be added by implementing
// a Notify(WebhookPayload) error method and registering them in the runner.
package notify
