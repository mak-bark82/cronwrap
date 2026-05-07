package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// LarkNotifier sends job notifications to a Lark (Feishu) webhook.
type LarkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewLarkNotifier creates a LarkNotifier that posts to the given webhook URL.
func NewLarkNotifier(webhookURL string) *LarkNotifier {
	return newLarkNotifierWithClient(webhookURL, &http.Client{})
}

func newLarkNotifierWithClient(webhookURL string, client *http.Client) *LarkNotifier {
	return &LarkNotifier{webhookURL: webhookURL, client: client}
}

// Notify sends a Lark message describing the job result.
func (n *LarkNotifier) Notify(jobName string, err error) error {
	var text string
	if err != nil {
		text = fmt.Sprintf(":red_circle: *%s* failed: %v", jobName, err)
	} else {
		text = fmt.Sprintf(":large_green_circle: *%s* succeeded", jobName)
	}

	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("lark: marshal payload: %w", marshalErr)
	}

	resp, postErr := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if postErr != nil {
		return fmt.Errorf("lark: post: %w", postErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark: unexpected status %d", resp.StatusCode)
	}
	return nil
}
