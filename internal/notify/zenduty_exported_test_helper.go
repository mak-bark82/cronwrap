package notify

import "net/http"

// NewZendutyNotifierWithClient creates a ZendutyNotifier with a custom HTTP
// client and base URL for testing.
func NewZendutyNotifierWithClient(integrationKey, url string, client *http.Client) *ZendutyNotifier {
	return newZendutyNotifierWithClient(integrationKey, url, client)
}
