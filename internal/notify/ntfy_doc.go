// Package notify provides notifier implementations for cronwrap.
//
// # Ntfy Notifier
//
// NtfyNotifier sends job status notifications to a ntfy topic.
// ntfy (https://ntfy.sh) is a simple HTTP-based pub-sub notification service
// that supports self-hosting.
//
// Example usage:
//
//	notifier := notify.NewNtfyNotifier("https://ntfy.sh", "my-cron-alerts")
//	notifier.Notify("backup-job", err)
//
// On failure, the notification is sent with high priority and includes the
// error message. On success, a default-priority message is sent.
package notify
