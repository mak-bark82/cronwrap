package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const opsGenieDefaultURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpsGenieNotifier creates an OpsGenieNotifier with the given API key.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey:  apiKey,
		baseURL: opsGenieDefaultURL,
		client:  &http.Client{},
	}
}

func newOpsGenieNotifierWithURL(apiKey, url string, client *http.Client) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey:  apiKey,
		baseURL: url,
		client:  client,
	}
}

// Notify sends an OpsGenie alert for the given job result.
func (n *OpsGenieNotifier) Notify(job string, err error) error {
	priority := "P5"
	message := fmt.Sprintf("cronwrap: job '%s' succeeded", job)
	if err != nil {
		priority = "P1"
		message = fmt.Sprintf("cronwrap: job '%s' failed: %v", job, err)
	}

	payload := map[string]interface{}{
		"message":  message,
		"alias":    fmt.Sprintf("cronwrap-%s", job),
		"priority": priority,
		"source":   "cronwrap",
		"tags":     []string{"cronwrap", job},
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", marshalErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, n.baseURL, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("opsgenie: create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+n.apiKey)

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("opsgenie: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
