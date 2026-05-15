package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultNewRelicURL = "https://log-api.newrelic.com/log/v1"

// NewRelicNotifier sends job notifications to New Relic Logs API.
type NewRelicNotifier struct {
	apiKey string
	url    string
	client *http.Client
}

// NewNewRelicNotifier creates a NewRelicNotifier using the given API key.
func NewNewRelicNotifier(apiKey string) *NewRelicNotifier {
	return newNewRelicNotifierWithClient(apiKey, defaultNewRelicURL, &http.Client{})
}

func newNewRelicNotifierWithClient(apiKey, url string, client *http.Client) *NewRelicNotifier {
	return &NewRelicNotifier{
		apiKey: apiKey,
		url:    url,
		client: client,
	}
}

// Notify sends a log event to New Relic for the given job.
func (n *NewRelicNotifier) Notify(jobName string, err error) error {
	status := "success"
	message := fmt.Sprintf("cronwrap: job %q completed successfully", jobName)
	if err != nil {
		status = "failure"
		message = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
	}

	payload := map[string]interface{}{
		"message":  message,
		"job_name": jobName,
		"status":   status,
		"logtype":  "cronwrap",
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("newrelic: marshal payload: %w", marshalErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("newrelic: create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", n.apiKey)

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("newrelic: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
