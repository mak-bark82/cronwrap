package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const victorOpsDefaultURL = "https://alert.victorops.com/integrations/generic/20131114/alert"

// VictorOpsNotifier sends alerts to VictorOps (Splunk On-Call) via REST endpoint.
type VictorOpsNotifier struct {
	apiKey   string
	routingKey string
	baseURL  string
	client   *http.Client
}

// NewVictorOpsNotifier creates a new VictorOpsNotifier with the given API key and routing key.
func NewVictorOpsNotifier(apiKey, routingKey string) *VictorOpsNotifier {
	return &VictorOpsNotifier{
		apiKey:     apiKey,
		routingKey: routingKey,
		baseURL:    victorOpsDefaultURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends an alert event to VictorOps.
func (v *VictorOpsNotifier) Notify(job string, err error) error {
	return v.notifyWithURL(fmt.Sprintf("%s/%s/%s", v.baseURL, v.apiKey, v.routingKey), job, err)
}

func (v *VictorOpsNotifier) notifyWithURL(url, job string, err error) error {
	msgType := "INFO"
	message := fmt.Sprintf("cronwrap: job %q completed successfully", job)
	if err != nil {
		msgType = "CRITICAL"
		message = fmt.Sprintf("cronwrap: job %q failed: %v", job, err)
	}

	payload := map[string]interface{}{
		"message_type":        msgType,
		"entity_id":           fmt.Sprintf("cronwrap/%s", job),
		"entity_display_name": fmt.Sprintf("cronwrap job: %s", job),
		"state_message":       message,
		"timestamp":           time.Now().Unix(),
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("victorops: marshal payload: %w", marshalErr)
	}

	resp, httpErr := v.client.Post(url, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("victorops: send request: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}
