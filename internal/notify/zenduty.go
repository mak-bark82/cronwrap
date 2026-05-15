package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const zendutyDefaultURL = "https://www.zenduty.com/api/incidents/"

// ZendutyNotifier sends alert notifications to Zenduty.
type ZendutyNotifier struct {
	integrationKey string
	url            string
	client         *http.Client
}

// NewZendutyNotifier creates a ZendutyNotifier with the given integration key.
func NewZendutyNotifier(integrationKey string) *ZendutyNotifier {
	return newZendutyNotifierWithClient(integrationKey, zendutyDefaultURL, &http.Client{})
}

func newZendutyNotifierWithClient(integrationKey, url string, client *http.Client) *ZendutyNotifier {
	return &ZendutyNotifier{
		integrationKey: integrationKey,
		url:            url,
		client:         client,
	}
}

// Notify sends a Zenduty alert for the given job. If err is non-nil the alert
// is treated as a failure.
func (z *ZendutyNotifier) Notify(jobName string, err error) error {
	message := fmt.Sprintf("cronwrap: job %q succeeded", jobName)
	alertType := "info"
	if err != nil {
		message = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
		alertType = "critical"
	}

	payload := map[string]string{
		"integration_key": z.integrationKey,
		"message":         message,
		"alert_type":      alertType,
		"summary":         message,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("zenduty: marshal payload: %w", marshalErr)
	}

	resp, httpErr := z.client.Post(z.url, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("zenduty: send alert: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
