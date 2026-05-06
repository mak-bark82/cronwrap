// Package notify provides notification integrations for cronwrap.
//
// # Gotify
//
// GotifyNotifier delivers alerts to a self-hosted Gotify push notification
// server (https://gotify.net). It requires a Gotify base URL and an
// application token.
//
// Usage:
//
//	n := notify.NewGotifyNotifier("https://gotify.example.com", "APP_TOKEN")
//	err := n.Notify("my-cron-job", jobErr)
//
// Message priority is set to 9 (high) on failure and 5 (normal) on success,
// matching Gotify's recommended priority scale.
package notify
