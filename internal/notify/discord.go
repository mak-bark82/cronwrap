package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordNotifier sends job alerts to a Discord channel via an incoming webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a DiscordNotifier that posts to the given webhook URL.
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type discordPayload struct {
	Content string `json:"content"`
}

// Notify sends a Discord message describing the job result.
func (d *DiscordNotifier) Notify(jobName string, err error) error {
	var msg string
	if err != nil {
		msg = fmt.Sprintf(":x: cronwrap: job *%s* failed: %v", jobName, err)
	} else {
		msg = fmt.Sprintf(":white_check_mark: cronwrap: job *%s* succeeded", jobName)
	}

	payload := discordPayload{Content: msg}
	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("discord: marshal payload: %w", encErr)
	}

	resp, httpErr := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("discord: send notification: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
