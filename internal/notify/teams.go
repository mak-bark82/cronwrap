package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsNotifier creates a TeamsNotifier that posts to the given webhook URL.
func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	ActivityText  string `json:"activityText"`
}

// Notify sends a notification to Microsoft Teams.
func (t *TeamsNotifier) Notify(jobName string, err error) error {
	color, summary, text := teamsFields(jobName, err)

	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: color,
		Summary:    summary,
		Sections: []teamsSection{
			{ActivityTitle: summary, ActivityText: text},
		},
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("teams: marshal payload: %w", encErr)
	}

	resp, httpErr := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("teams: post: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func teamsFields(jobName string, err error) (color, summary, text string) {
	if err != nil {
		return "FF0000",
			fmt.Sprintf("cronwrap: job '%s' failed", jobName),
			fmt.Sprintf("Job **%s** failed with error: %s", jobName, err.Error())
	}
	return "00FF00",
		fmt.Sprintf("cronwrap: job '%s' succeeded", jobName),
		fmt.Sprintf("Job **%s** completed successfully.", jobName)
}
