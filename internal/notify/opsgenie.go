package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey string
	url    string
	client *http.Client
}

// NewOpsGenieNotifier creates a new OpsGenieNotifier with the given API key.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey: apiKey,
		url:    defaultOpsGenieURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type opsGeniePayload struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority"`
	Source      string `json:"source"`
}

// Notify sends an alert to OpsGenie with job name and optional error details.
func (o *OpsGenieNotifier) Notify(jobName string, err error) error {
	priority := "P5"
	description := fmt.Sprintf("Cron job '%s' completed successfully.", jobName)

	if err != nil {
		priority = "P1"
		description = fmt.Sprintf("Cron job '%s' failed: %s", jobName, err.Error())
	}

	payload := opsGeniePayload{
		Message:     fmt.Sprintf("cronwrap: %s", jobName),
		Description: description,
		Priority:    priority,
		Source:      "cronwrap",
	}

	body, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", jsonErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, o.url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("opsgenie: failed to create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, doErr := o.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("opsgenie: request failed: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
