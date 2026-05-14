package notify

import "net/http"

// NewDingTalkNotifierWithClient creates a DingTalkNotifier with a custom HTTP client.
// Exported for use in external tests.
func NewDingTalkNotifierWithClient(webhookURL string, client *http.Client) *DingTalkNotifier {
	return newDingTalkNotifierWithClient(webhookURL, client)
}
