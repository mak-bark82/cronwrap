// Package notify provides integrations for sending job alerts to various
// external services.
//
// # VictorOps (Splunk On-Call)
//
// VictorOpsNotifier sends alert events to the VictorOps REST endpoint.
// It maps a successful job run to an INFO message type and a failed run
// to a CRITICAL message type so that on-call engineers are paged only
// when action is required.
//
// Usage:
//
//	n := notify.NewVictorOpsNotifier(apiKey, routingKey)
//	if err := n.Notify(jobName, jobErr); err != nil {
//		log.Printf("alert delivery failed: %v", err)
//	}
package notify
