package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DingTalkNotifier sends job alerts to a DingTalk webhook.
type DingTalkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDingTalkNotifier creates a DingTalkNotifier that posts to the given webhook URL.
func NewDingTalkNotifier(webhookURL string) *DingTalkNotifier {
	return newDingTalkNotifierWithClient(webhookURL, &http.Client{})
}

func newDingTalkNotifierWithClient(webhookURL string, client *http.Client) *DingTalkNotifier {
	return &DingTalkNotifier{webhookURL: webhookURL, client: client}
}

// Notify sends a DingTalk message for the given job result.
func (d *DingTalkNotifier) Notify(job string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf("❌ Job *%s* failed: %v", job, err)
	} else {
		text = fmt.Sprintf("✅ Job *%s* succeeded.", job)
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("dingtalk: marshal payload: %w", encErr)
	}

	resp, httpErr := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if httpErr != nil {
		return fmt.Errorf("dingtalk: post: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
