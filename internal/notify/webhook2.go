// Package notify provides notification integrations for cronwrap.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LineNotifier sends notifications via LINE Notify API.
type LineNotifier struct {
	token  string
	url    string
	client *http.Client
}

// NewLineNotifier creates a LineNotifier using the given LINE Notify token.
func NewLineNotifier(token string) *LineNotifier {
	return newLineNotifierWithClient(token, "https://notify-api.line.me/api/notify", &http.Client{Timeout: 10 * time.Second})
}

func newLineNotifierWithClient(token, url string, client *http.Client) *LineNotifier {
	return &LineNotifier{token: token, url: url, client: client}
}

// Notify sends a LINE Notify message for the given job result.
func (n *LineNotifier) Notify(job string, err error) error {
	body := lineBody(job, err)
	payload, _ := json.Marshal(map[string]string{"message": body})

	req, reqErr := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(payload))
	if reqErr != nil {
		return fmt.Errorf("line: create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.token)

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("line: send request: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("line: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func lineBody(job string, err error) string {
	if err != nil {
		return fmt.Sprintf("[cronwrap] Job %q failed: %v", job, err)
	}
	return fmt.Sprintf("[cronwrap] Job %q completed successfully", job)
}
