package notify

import "net/http"

// NewLineNotifierWithClient creates a LineNotifier with a custom HTTP client for testing.
func NewLineNotifierWithClient(token, url string, client *http.Client) *LineNotifier {
	return newLineNotifierWithClient(token, url, client)
}
