package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultSplunkURL = "https://input.splunk.com/services/collector/event"

// SplunkNotifier sends job alerts to a Splunk HTTP Event Collector (HEC).
type SplunkNotifier struct {
	token string
	url   string
	client *http.Client
}

type splunkEvent struct {
	Event splunkPayload `json:"event"`
}

type splunkPayload struct {
	Job     string `json:"job"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewSplunkNotifier creates a SplunkNotifier that posts to Splunk HEC.
func NewSplunkNotifier(token string) *SplunkNotifier {
	return newSplunkNotifierWithClient(token, defaultSplunkURL, &http.Client{})
}

func newSplunkNotifierWithClient(token, url string, client *http.Client) *SplunkNotifier {
	return &SplunkNotifier{token: token, url: url, client: client}
}

// Notify sends a Splunk HEC event for the given job outcome.
func (s *SplunkNotifier) Notify(jobName string, err error) error {
	status := "success"
	var msg string
	if err != nil {
		status = "failure"
		msg = err.Error()
	}

	payload := splunkEvent{
		Event: splunkPayload{
			Job:     jobName,
			Status:  status,
			Message: msg,
		},
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("splunk: marshal payload: %w", encErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("splunk: create request: %w", reqErr)
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, doErr := s.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("splunk: send event: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
