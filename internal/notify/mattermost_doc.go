// Package notify provides notifier implementations for alerting on cron job
// outcomes via various channels.
//
// # Mattermost
//
// MattermostNotifier sends alerts to a Mattermost channel via an incoming
// webhook URL. Configure a webhook in your Mattermost instance and pass the
// URL to NewMattermostNotifier.
//
// Example:
//
//	n := notify.NewMattermostNotifier("https://mattermost.example.com/hooks/xxx")
//	if err := n.Notify("my-job", jobErr); err != nil {
//		log.Println("alert failed:", err)
//	}
package notify
