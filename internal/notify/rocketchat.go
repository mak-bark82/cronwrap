package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RocketChatNotifier sends job notifications to a Rocket.Chat incoming webhook.
type RocketChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewRocketChatNotifier creates a new RocketChatNotifier with the given webhook URL.
func NewRocketChatNotifier(webhookURL string) *RocketChatNotifier {
	return &RocketChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func newRocketChatNotifierWithClient(webhookURL string, client *http.Client) *RocketChatNotifier {
	return &RocketChatNotifier{
		webhookURL: webhookURL,
		client:     client,
	}
}

// Notify sends a notification to Rocket.Chat with the job name and optional error.
func (r *RocketChatNotifier) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf(":x: Job *%s* failed: %s", jobName, err.Error())
	} else {
		text = fmt.Sprintf(":white_check_mark: Job *%s* completed successfully.", jobName)
	}

	payload := map[string]string{"text": text}
	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", marshalErr)
	}

	resp, httpErr := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("rocketchat: send notification: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
