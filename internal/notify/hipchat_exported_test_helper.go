package notify

import "net/http"

// NewHipChatNotifierWithClient creates a HipChatNotifier with a custom HTTP
// client and base URL, used in tests to inject a fake server.
func NewHipChatNotifierWithClient(token, roomID, baseURL string, client *http.Client) *HipChatNotifier {
	return newHipChatNotifierWithClient(token, roomID, baseURL, client)
}
