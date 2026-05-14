package notify

import "net/http"

// NewCircleCINotifierWithClient creates a CircleCINotifier with a custom
// HTTP client and API URL, intended for use in tests only.
func NewCircleCINotifierWithClient(token, apiURL string, client *http.Client) *CircleCINotifier {
	return newCircleCINotifierWithClient(token, apiURL, client)
}
