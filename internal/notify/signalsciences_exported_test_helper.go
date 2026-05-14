package notify

import "net/http"

// NewSignalSciencesNotifierWithClient creates a SignalSciencesNotifier with a
// custom HTTP client, intended for use in tests.
func NewSignalSciencesNotifierWithClient(corpName, siteName, apiToken string, client *http.Client) *SignalSciencesNotifier {
	return newSignalSciencesNotifierWithClient(corpName, siteName, apiToken, client)
}
