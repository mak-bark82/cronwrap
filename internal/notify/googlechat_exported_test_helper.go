package notify

import "net/http"

// NewGoogleChatNotifierWithClient creates a GoogleChatNotifier with a custom
// HTTP client, intended for use in tests.
func NewGoogleChatNotifierWithClient(webhookURL string, client *http.Client) *GoogleChatNotifier {
	return newGoogleChatNotifierWithClient(webhookURL, client)
}
