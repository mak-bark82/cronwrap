package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BearyChat sends notifications to a BearyChat incoming webhook.
type BearyChat struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChat creates a new BearyChat notifier with the given webhook URL.
func NewBearyChat(webhookURL string) *BearyChat {
	return newBearyChat(webhookURL, &http.Client{})
}

func newBearyChat(webhookURL string, client *http.Client) *BearyChat {
	return &BearyChat{
		webhookURL: webhookURL,
		client:     client,
	}
}

// Notify sends a job result notification to BearyChat.
func (b *BearyChat) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf("❌ Job *%s* failed: %s", jobName, err.Error())
	} else {
		text = fmt.Sprintf("✅ Job *%s* succeeded.", jobName)
	}

	payload := map[string]interface{}{
		"text":     text,
		"markdown": true,
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("bearychat: marshal payload: %w", encErr)
	}

	resp, httpErr := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("bearychat: send notification: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
