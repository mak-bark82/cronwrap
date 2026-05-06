// Package notify provides notifier implementations for alerting on cron job
// outcomes. Each notifier implements a common Notify(jobName string, err error)
// interface so they can be composed via MultiNotifier.
//
// # Microsoft Teams
//
// TeamsNotifier delivers alerts to a Microsoft Teams channel using an
// Incoming Webhook connector URL. The card colour is green on success and
// red on failure, and the error message (if any) is included in the card body.
//
// Usage:
//
//	n := notify.NewTeamsNotifier("https://example.webhook.office.com/...")
//	if err := n.Notify("my-job", jobErr); err != nil {
//		log.Printf("teams alert failed: %v", err)
//	}
package notify
