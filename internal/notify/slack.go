package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends alert notifications to a Slack webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier that posts to the given Slack
// incoming-webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends a Slack message describing the job outcome.
func (s *SlackNotifier) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf(":red_circle: cronwrap: job *%s* failed: %v", jobName, err)
	} else {
		text = fmt.Sprintf(":white_check_mark: cronwrap: job *%s* succeeded", jobName)
	}

	body, encErr := json.Marshal(slackPayload{Text: text})
	if encErr != nil {
		return fmt.Errorf("slack: marshal payload: %w", encErr)
	}

	resp, reqErr := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("slack: post: %w", reqErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
