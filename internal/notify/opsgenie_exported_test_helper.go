package notify

import "net/http"

// NewOpsGenieNotifierWithURL creates an OpsGenieNotifier with a custom URL for testing.
func NewOpsGenieNotifierWithURL(apiKey, url string, client *http.Client) *OpsGenieNotifier {
	n := NewOpsGenieNotifier(apiKey)
	n.url = url
	n.client = client
	return n
}
