package notify

import "net/http"

// NewMatrixNotifierWithClient creates a MatrixNotifier with a custom HTTP
// client, intended for use in tests only.
func NewMatrixNotifierWithClient(homeserver, token, roomID string, client *http.Client) *MatrixNotifier {
	return newMatrixNotifierWithClient(homeserver, token, roomID, client)
}
