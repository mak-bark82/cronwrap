package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// notifyWithURL is an internal helper used by tests to override the PagerDuty URL.
func notifyWithURL(n *PagerDutyNotifier, url, jobName string, err error) error {
	summary := fmt.Sprintf("cronwrap: job %q succeeded", jobName)
	if err != nil {
		summary = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
	}
	body := pdPayload{
		RoutingKey:  n.integrationKey,
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
	resp, httpErr := n.client.Post(url, "application/json", bytes.NewReader(data))
	if httpErr != nil {
		return fmt.Errorf("pagerduty: send event: %w", httpErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func TestSeverity_WithError(t *testing.T) {
	if s := severity(errors.New("oops")); s != "error" {
		t.Errorf("expected %q, got %q", "error", s)
	}
}

func TestSeverity_NoError(t *testing.T) {
	if s := severity(nil); s != "info" {
		t.Errorf("expected %q, got %q", "info", s)
	}
}

func TestNotifyWithURL_PayloadContainsJobName(t *testing.T) {
	var body string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sb strings.Builder
		_, _ = sb.ReadFrom(r.Body)
		body = sb.String()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("routing-key")
	n.client = ts.Client()
	_ = notifyWithURL(n, ts.URL, "my-cron-job", nil)

	if !strings.Contains(body, "my-cron-job") {
		t.Errorf("expected payload to contain job name, got: %s", body)
	}
	if !strings.Contains(body, "routing-key") {
		t.Errorf("expected payload to contain routing key, got: %s", body)
	}
}
