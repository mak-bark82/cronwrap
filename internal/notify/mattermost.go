package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MattermostNotifier sends job alerts to a Mattermost incoming webhook.
type MattermostNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewMattermostNotifier creates a new MattermostNotifier with the given webhook URL.
func NewMattermostNotifier(webhookURL string) *MattermostNotifier {
	return &MattermostNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

newMattermostNotifierWithClient := func(webhookURL string, client *http.Client) *MattermostNotifier {
	return &MattermostNotifier{webhookURL: webhookURL, client: client}
}

func newMattermostNotifierWithClient(webhookURL string, client *http.Client) *MattermostNotifier {
	return &MattermostNotifier{webhookURL: webhookURL, client: client}
}

// Notify sends a message to Mattermost describing the job result.
func (m *MattermostNotifier) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf(":x: Cron job *%s* failed: %s", jobName, err.Error())
	} else {
		text = fmt.Sprintf(":white_check_mark: Cron job *%s* succeeded.", jobName)
	}

	payload := map[string]string{"text": text}
	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", encErr)
	}

	resp, httpErr := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("mattermost: http post: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
