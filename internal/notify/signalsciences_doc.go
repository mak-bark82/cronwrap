// Package notify provides notification integrations for cronwrap.
//
// # Signal Sciences
//
// SignalSciencesNotifier sends job completion events to the Signal Sciences
// (Fastly Next-Gen WAF) event feed API. It reports both successful and failed
// job runs, including the job name and error message when applicable.
//
// Usage:
//
//	n := notify.NewSignalSciencesNotifier("my-corp", "my-site", "api-token")
//	err := n.Notify("db-backup", jobErr)
//
// The notifier sets the x-api-user and x-api-token headers required by the
// Signal Sciences API and posts a JSON payload containing the event type,
// job name, status, and human-readable message.
package notify
