package notify

import "net/http"

// NewNewRelicNotifierWithClient creates a NewRelicNotifier with a custom HTTP
// client, intended for use in tests.
func NewNewRelicNotifierWithClient(apiKey, url string, client *http.Client) *NewRelicNotifier {
	return newNewRelicNotifierWithClient(apiKey, url, client)
}
