package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleChatNotifier sends job notifications to a Google Chat webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a GoogleChatNotifier that posts to the given
// Google Chat incoming webhook URL.
func NewGoogleChatNotifier(webhookURL string) *GoogleChatNotifier {
	return newGoogleChatNotifierWithClient(webhookURL, &http.Client{})
}

func newGoogleChatNotifierWithClient(webhookURL string, client *http.Client) *GoogleChatNotifier {
	return &GoogleChatNotifier{webhookURL: webhookURL, client: client}
}

// Notify sends a Google Chat message describing the job outcome.
func (n *GoogleChatNotifier) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf("❌ *%s* failed: %s", jobName, err.Error())
	} else {
		text = fmt.Sprintf("✅ *%s* succeeded.", jobName)
	}

	payload := map[string]interface{}{
		"text": text,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", marshalErr)
	}

	resp, httpErr := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("googlechat: post: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}

	return nil
}
