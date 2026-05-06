package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const datadogDefaultURL = "https://api.datadoghq.com/api/v1/events"

// DatadogNotifier sends job events to the Datadog Events API.
type DatadogNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

// NewDatadogNotifier creates a DatadogNotifier using the given API key.
func NewDatadogNotifier(apiKey string) *DatadogNotifier {
	return &DatadogNotifier{
		apiKey: apiKey,
		apiURL: datadogDefaultURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type datadogPayload struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	AlertType string `json:"alert_type"`
	Tags  []string `json:"tags"`
}

// Notify sends a Datadog event for the given job result.
func (d *DatadogNotifier) Notify(jobName string, err error) error {
	return d.notifyWithURL(d.apiURL, jobName, err)
}

func (d *DatadogNotifier) notifyWithURL(url, jobName string, err error) error {
	alertType := "success"
	text := fmt.Sprintf("Job %q completed successfully.", jobName)
	if err != nil {
		alertType = "error"
		text = fmt.Sprintf("Job %q failed: %v", jobName, err)
	}

	payload := datadogPayload{
		Title:     fmt.Sprintf("cronwrap: %s", jobName),
		Text:      text,
		AlertType: alertType,
		Tags:      []string{"source:cronwrap", fmt.Sprintf("job:%s", jobName)},
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("datadog: marshal payload: %w", marshalErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("datadog: create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, doErr := d.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("datadog: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}
