package notify

import "net/http"

// NewLarkNotifierWithClient creates a LarkNotifier with a custom HTTP client.
// This is exported for use in external test packages to allow injecting a mock
// or test HTTP client without exposing internal constructors in production code.
func NewLarkNotifierWithClient(webhookURL string, client *http.Client) *LarkNotifier {
	return newLarkNotifierWithClient(webhookURL, client)
}

// NewLarkNotifierWithDefaults creates a LarkNotifier using the default HTTP client.
// Exported for use in external test packages.
func NewLarkNotifierWithDefaults(webhookURL string) *LarkNotifier {
	return newLarkNotifierWithClient(webhookURL, http.DefaultClient)
}
