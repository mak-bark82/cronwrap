package notify

import "net/http"

// NewMattermostNotifierWithClient creates a MattermostNotifier with a custom
// HTTP client, intended for use in tests.
func NewMattermostNotifierWithClient(webhookURL string, client *http.Client) *MattermostNotifier {
	return newMattermostNotifierWithClient(webhookURL, client)
}
