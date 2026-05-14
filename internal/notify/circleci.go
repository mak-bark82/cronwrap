package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/exampleorg/cronwrap/internal/alert"
)

const circleciAPIURL = "https://circleci.com/api/v2/insights/events"

// CircleCINotifier sends job notifications to the CircleCI Insights API.
type CircleCINotifier struct {
	token  string
	apiURL string
	client *http.Client
}

// NewCircleCINotifier creates a CircleCINotifier that authenticates with
// the provided API token.
func NewCircleCINotifier(token string) *CircleCINotifier {
	return &CircleCINotifier{
		token:  token,
		apiURL: circleciAPIURL,
		client: &http.Client{},
	}
}

func newCircleCINotifierWithClient(token, apiURL string, client *http.Client) *CircleCINotifier {
	return &CircleCINotifier{token: token, apiURL: apiURL, client: client}
}

// Notify sends the job event to CircleCI Insights.
func (n *CircleCINotifier) Notify(ev alert.Event) error {
	status := "success"
	if ev.Err != nil {
		status = "failed"
	}

	payload := map[string]string{
		"job_name": ev.JobName,
		"status":   status,
	}
	if ev.Err != nil {
		payload["error"] = ev.Err.Error()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("circleci notifier: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("circleci notifier: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", n.token)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("circleci notifier: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("circleci notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
