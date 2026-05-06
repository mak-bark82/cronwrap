package notify

import "net/http"

// NewGotifyNotifierWithClient creates a GotifyNotifier with a custom HTTP
// client, intended for use in tests that need to inspect or intercept requests.
func NewGotifyNotifierWithClient(baseURL, token string, client *http.Client) *GotifyNotifier {
	return &GotifyNotifier{
		baseURL: baseURL,
		token:   token,
		client:  client,
	}
}
