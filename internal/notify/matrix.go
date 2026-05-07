package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MatrixNotifier sends notifications to a Matrix room via the Client-Server API.
type MatrixNotifier struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

// NewMatrixNotifier creates a MatrixNotifier that posts to the given Matrix
// homeserver and room using the provided access token.
func NewMatrixNotifier(homeserver, token, roomID string) *MatrixNotifier {
	return &MatrixNotifier{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}
}

func newMatrixNotifierWithClient(homeserver, token, roomID string, client *http.Client) *MatrixNotifier {
	n := NewMatrixNotifier(homeserver, token, roomID)
	n.client = client
	return n
}

// Notify implements the Notifier interface.
func (n *MatrixNotifier) Notify(jobName string, err error) error {
	var body string
	if err != nil {
		body = fmt.Sprintf("❌ cronwrap: job *%s* failed: %s", jobName, err.Error())
	} else {
		body = fmt.Sprintf("✅ cronwrap: job *%s* succeeded.", jobName)
	}

	payload := map[string]string{
		"msgtype": "m.text",
		"body":    body,
	}

	data, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("matrix: marshal payload: %w", encErr)
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		n.homeserver, n.roomID)

	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if reqErr != nil {
		return fmt.Errorf("matrix: build request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.token)

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("matrix: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
