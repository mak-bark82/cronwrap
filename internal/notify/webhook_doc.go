// Package notify provides notification integrations for cronwrap.
//
// The WebhookNotifier sends an HTTP POST request to a configurable URL
// whenever a job completes. The payload is a JSON object containing the
// job name, success status, and any error message.
//
// Example usage:
//
//	notifier := notify.NewWebhookNotifier("https://example.com/hook")
//	err := notifier.Notify(alert.Event{
//		JobName: "my-job",
//		Err:     nil,
//	})
package notify
