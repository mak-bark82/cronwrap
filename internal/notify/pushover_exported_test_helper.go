package notify

import "net/http"

// NewPushoverNotifierWithClient creates a PushoverNotifier with a custom HTTP client and URL,
// intended for use in tests only.
func NewPushoverNotifierWithClient(token, userKey, url string, client *http.Client) *PushoverNotifier {
	return newPushoverNotifierWithClient(token, userKey, url, client)
}
