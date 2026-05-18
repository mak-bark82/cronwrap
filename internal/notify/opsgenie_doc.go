// Package notify provides integrations with external notification services.
//
// OpsGenie
//
// The OpsGenieNotifier sends alerts to OpsGenie using the Alert API v2.
// Jobs that fail are sent with priority P1; successful jobs use P5.
//
// Usage:
//
//	n := notify.NewOpsGenieNotifier("<your-api-key>")
//	err := n.Notify("my-cron-job", jobErr)
//
// Authentication is performed via the GenieKey scheme in the Authorization
// header. The alias field is set to "cronwrap-<jobname>" to allow OpsGenie
// to deduplicate repeated alerts for the same job.
package notify
