package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const statuspageDefaultURL = "https://api.statuspage.io/v1"

// StatuspageNotifier sends incident notifications to Atlassian Statuspage.
type StatuspageNotifier struct {
	apiKey     string
	pageID     string
	componentID string
	url        string
	client     *http.Client
}

// NewStatuspageNotifier creates a new StatuspageNotifier.
func NewStatuspageNotifier(apiKey, pageID, componentID string) *StatuspageNotifier {
	return &StatuspageNotifier{
		apiKey:      apiKey,
		pageID:      pageID,
		componentID: componentID,
		url:         statuspageDefaultURL,
		client:      &http.Client{},
	}
}

func newStatuspageNotifierWithClient(apiKey, pageID, componentID, url string, client *http.Client) *StatuspageNotifier {
	n := NewStatuspageNotifier(apiKey, pageID, componentID)
	n.url = url
	n.client = client
	return n
}

// Notify sends an incident to Statuspage if err is non-nil, or resolves it on success.
func (n *StatuspageNotifier) Notify(jobName string, err error) error {
	status := "operational"
	if err != nil {
		status = "major_outage"
	}

	body := map[string]interface{}{
		"component": map[string]string{
			"status": status,
		},
	}

	payload, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		return fmt.Errorf("statuspage: marshal payload: %w", jsonErr)
	}

	endpoint := fmt.Sprintf("%s/pages/%s/components/%s", n.url, n.pageID, n.componentID)
	req, reqErr := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(payload))
	if reqErr != nil {
		return fmt.Errorf("statuspage: build request: %w", reqErr)
	}
	req.Header.Set("Authorization", "OAuth "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("statuspage: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
