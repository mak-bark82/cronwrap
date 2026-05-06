package notify

// Expose internal fields for testing via exported accessor.
// This file is compiled into the notify package so tests in
// the notify_test (external) package cannot directly set private fields.
// Instead, we expose a thin helper that external tests can use.

// VictorOpsNotifier re-exports the unexported type so external tests
// can type-assert and set BaseURL for test servers.
func (v *VictorOpsNotifier) SetBaseURL(u string) {
	v.baseURL = u
}

// Ensure VictorOpsNotifier satisfies the Notifier interface.
var _ Notifier = (*VictorOpsNotifier)(nil)
