package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverNotifier sends notifications via the Pushover API.
type PushoverNotifier struct {
	token   string
	userKey string
	url     string
	client  *http.Client
}

// NewPushoverNotifier creates a new PushoverNotifier with the given API token and user key.
func NewPushoverNotifier(token, userKey string) *PushoverNotifier {
	return &PushoverNotifier{
		token:   token,
		userKey: userKey,
		url:     pushoverAPIURL,
		client:  &http.Client{},
	}
}

func newPushoverNotifierWithClient(token, userKey, url string, client *http.Client) *PushoverNotifier {
	return &PushoverNotifier{
		token:   token,
		userKey: userKey,
		url:     url,
		client:  client,
	}
}

// Notify sends a Pushover notification for the given job result.
func (p *PushoverNotifier) Notify(job string, err error) error {
	title := fmt.Sprintf("cronwrap: %s", job)
	var message string
	if err != nil {
		message = fmt.Sprintf("Job failed: %s", err.Error())
	} else {
		message = "Job completed successfully."
	}

	payload := map[string]string{
		"token":   p.token,
		"user":    p.userKey,
		"title":   title,
		"message": message,
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("pushover: marshal payload: %w", encErr)
	}

	resp, httpErr := p.client.Post(p.url, "application/json", strings.NewReader(string(body)))
	if httpErr != nil {
		return fmt.Errorf("pushover: send request: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
