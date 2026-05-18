package notify

import "net/http"

// NewOpsGenieNotifierWithURL creates an OpsGenieNotifier with a custom URL and
// HTTP client, intended for use in tests.
func NewOpsGenieNotifierWithURL(apiKey, url string, client *http.Client) *OpsGenieNotifier {
	return newOpsGenieNotifierWithURL(apiKey, url, client)
}
