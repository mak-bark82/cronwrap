package notify

import (
	"bytes"
	"fmt"
	"net/http"
)

// NtfyNotifier sends notifications via ntfy.sh or a self-hosted ntfy instance.
type NtfyNotifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyNotifier creates a new NtfyNotifier targeting the given server and topic.
// serverURL should be the base URL (e.g. "https://ntfy.sh" or your self-hosted instance).
func NewNtfyNotifier(serverURL, topic string) *NtfyNotifier {
	return &NtfyNotifier{
		serverURL: serverURL,
		topic:     topic,
		client:    &http.Client{},
	}
}

func newNtfyNotifierWithClient(serverURL, topic string, client *http.Client) *NtfyNotifier {
	return &NtfyNotifier{
		serverURL: serverURL,
		topic:     topic,
		client:    client,
	}
}

// Notify sends a message to the configured ntfy topic.
func (n *NtfyNotifier) Notify(jobName string, err error) error {
	var body string
	var title string
	var priority string

	if err != nil {
		title = fmt.Sprintf("cronwrap: %s FAILED", jobName)
		body = fmt.Sprintf("Job %q failed: %v", jobName, err)
		priority = "high"
	} else {
		title = fmt.Sprintf("cronwrap: %s succeeded", jobName)
		body = fmt.Sprintf("Job %q completed successfully.", jobName)
		priority = "default"
	}

	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)
	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if reqErr != nil {
		return fmt.Errorf("ntfy: failed to build request: %w", reqErr)
	}

	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority)
	req.Header.Set("Content-Type", "text/plain")

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("ntfy: request failed: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
