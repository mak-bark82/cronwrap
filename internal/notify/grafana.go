package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GrafanaNotifier sends annotations to a Grafana instance when a job
// completes, making it easy to correlate cron job runs with dashboards.
type GrafanaNotifier struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewGrafanaNotifier creates a GrafanaNotifier that posts annotations to
// the given Grafana baseURL (e.g. "https://grafana.example.com") using
// the provided API key.
func NewGrafanaNotifier(baseURL, apiKey string) *GrafanaNotifier {
	return newGrafanaNotifierWithClient(baseURL, apiKey, &http.Client{})
}

func newGrafanaNotifierWithClient(baseURL, apiKey string, client *http.Client) *GrafanaNotifier {
	return &GrafanaNotifier{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
	}
}

// Notify posts a Grafana annotation describing the job result.
func (g *GrafanaNotifier) Notify(jobName string, err error) error {
	text := fmt.Sprintf("cronwrap: job %q succeeded", jobName)
	tags := []string{"cronwrap", jobName, "success"}
	if err != nil {
		text = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
		tags = []string{"cronwrap", jobName, "failure"}
	}

	payload := map[string]interface{}{
		"text": text,
		"tags": tags,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("grafana: marshal payload: %w", marshalErr)
	}

	url := g.baseURL + "/api/annotations"
	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("grafana: build request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, doErr := g.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("grafana: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status %d", resp.StatusCode)
	}
	return nil
}
