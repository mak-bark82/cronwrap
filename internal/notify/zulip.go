package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ibrahimker/cronwrap/internal/alert"
)

// ZulipNotifier sends notifications to a Zulip stream via the Zulip REST API.
type ZulipNotifier struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	client   *http.Client
}

// NewZulipNotifier creates a ZulipNotifier that posts to the given Zulip
// server. baseURL should be e.g. "https://yourorg.zulipchat.com".
func NewZulipNotifier(baseURL, email, apiKey, stream, topic string) *ZulipNotifier {
	return &ZulipNotifier{
		baseURL: baseURL,
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}
}

func newZulipNotifierWithClient(baseURL, email, apiKey, stream, topic string, client *http.Client) *ZulipNotifier {
	n := NewZulipNotifier(baseURL, email, apiKey, stream, topic)
	n.client = client
	return n
}

// Notify implements alert.Notifier. It sends a message to the configured
// Zulip stream and topic.
func (z *ZulipNotifier) Notify(e alert.Event) error {
	body := fmt.Sprintf("**Job:** %s\n**Status:** %s", e.JobName, statusText(e.Err))
	if e.Err != nil {
		body += fmt.Sprintf("\n**Error:** %s", e.Err.Error())
	}

	payload, err := json.Marshal(map[string]string{
		"type":    "stream",
		"to":      z.stream,
		"topic":   z.topic,
		"content": body,
	})
	if err != nil {
		return fmt.Errorf("zulip: marshal payload: %w", err)
	}

	url := z.baseURL + "/api/v1/messages"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("zulip: build request: %w", err)
	}
	req.SetBasicAuth(z.email, z.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// statusText returns a human-readable status string.
func statusText(err error) string {
	if err != nil {
		return "FAILED"
	}
	return "SUCCESS"
}
