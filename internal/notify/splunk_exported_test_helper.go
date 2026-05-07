package notify

import "net/http"

// NewSplunkNotifierWithClient creates a SplunkNotifier with a custom HTTP
// client and endpoint, intended for use in tests.
func NewSplunkNotifierWithClient(token, url string, client *http.Client) *SplunkNotifier {
	return newSplunkNotifierWithClient(token, url, client)
}
