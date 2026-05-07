package notify

import "net/http"

// NewLarkNotifierWithClient creates a LarkNotifier with a custom HTTP client.
// Exported for use in external test packages.
func NewLarkNotifierWithClient(webhookURL string, client *http.Client) *LarkNotifier {
	return newLarkNotifierWithClient(webhookURL, client)
}
