package notify

import "net/http"

// NewGrafanaNotifierWithClient creates a GrafanaNotifier with a custom
// HTTP client, intended for use in tests only.
func NewGrafanaNotifierWithClient(baseURL, apiKey string, client *http.Client) *GrafanaNotifier {
	return newGrafanaNotifierWithClient(baseURL, apiKey, client)
}
