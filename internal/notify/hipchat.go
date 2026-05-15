package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultHipChatURL = "https://api.hipchat.com/v2/room/%s/notification"

// HipChatNotifier sends job notifications to a HipChat room.
type HipChatNotifier struct {
	token  string
	roomID string
	baseURL string
	client  *http.Client
}

type hipChatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// NewHipChatNotifier creates a HipChatNotifier for the given room and token.
func NewHipChatNotifier(token, roomID string) *HipChatNotifier {
	return newHipChatNotifierWithClient(token, roomID, defaultHipChatURL, &http.Client{})
}

func newHipChatNotifierWithClient(token, roomID, baseURL string, client *http.Client) *HipChatNotifier {
	return &HipChatNotifier{
		token:   token,
		roomID:  roomID,
		baseURL: baseURL,
		client:  client,
	}
}

// Notify sends a notification to the configured HipChat room.
func (n *HipChatNotifier) Notify(job string, err error) error {
	color := "green"
	var msg string
	if err != nil {
		color = "red"
		msg = fmt.Sprintf("cronwrap: job '%s' failed: %v", job, err)
	} else {
		msg = fmt.Sprintf("cronwrap: job '%s' succeeded", job)
	}

	payload := hipChatPayload{
		Message:       msg,
		MessageFormat: "text",
		Color:         color,
		Notify:        err != nil,
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("hipchat: marshal payload: %w", encErr)
	}

	url := fmt.Sprintf(n.baseURL, n.roomID)
	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if reqErr != nil {
		return fmt.Errorf("hipchat: create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.token)

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("hipchat: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
