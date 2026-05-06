package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
}

// NewPagerDutyNotifier creates a new PagerDutyNotifier with the given integration key.
func NewPagerDutyNotifier(integrationKey string) *PagerDutyNotifier {
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

// Notify sends a trigger event to PagerDuty.
func (p *PagerDutyNotifier) Notify(jobName string, err error) error {
	summary := fmt.Sprintf("cronwrap: job %q succeeded", jobName)
	if err != nil {
		summary = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
	}

	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  summary,
			Source:   "cronwrap",
			Severity: severity(err),
		},
	}

	data, encErr := json.Marshal(body)
	if encErr != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", encErr)
	}

	resp, httpErr := p.client.Post(pagerDutyEventURL, "application/json", bytes.NewReader(data))
	if httpErr != nil {
		return fmt.Errorf("pagerduty: send event: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func severity(err error) string {
	if err != nil {
		return "error"
	}
	return "info"
}
