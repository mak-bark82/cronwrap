package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GotifyNotifier sends notifications to a self-hosted Gotify server.
type GotifyNotifier struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewGotifyNotifier creates a new GotifyNotifier with the given server URL and app token.
func NewGotifyNotifier(baseURL, token string) *GotifyNotifier {
	return &GotifyNotifier{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}
}

type gotifyPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Priority int    `json:"priority"`
}

// Notify sends a message to Gotify. Priority is elevated when err is non-nil.
func (g *GotifyNotifier) Notify(jobName string, err error) error {
	priority := 5
	message := fmt.Sprintf("Job %q completed successfully.", jobName)
	if err != nil {
		priority = 9
		message = fmt.Sprintf("Job %q failed: %v", jobName, err)
	}

	payload := gotifyPayload{
		Title:    fmt.Sprintf("cronwrap: %s", jobName),
		Message:  message,
		Priority: priority,
	}

	body, err2 := json.Marshal(payload)
	if err2 != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err2)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.baseURL, g.token)
	resp, err2 := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err2 != nil {
		return fmt.Errorf("gotify: send request: %w", err2)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
