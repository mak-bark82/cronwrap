package notify

import "net/http"

// NewNtfyNotifierWithClient creates an NtfyNotifier with a custom HTTP client.
// Exported for use in external test packages.
func NewNtfyNotifierWithClient(serverURL, topic string, client *http.Client) *NtfyNotifier {
	return newNtfyNotifierWithClient(serverURL, topic, client)
}
