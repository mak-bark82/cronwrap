package notify

import "net/http"

// NewBearyChat creates a BearyChat notifier with a custom HTTP client.
// This is exported for use in external test packages.
func NewBearyChat(webhookURL string, client *http.Client) *BearyChat {
	return newBearyChat(webhookURL, client)
}
