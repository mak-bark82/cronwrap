package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPPostNotifier sends a generic JSON POST notification to a configurable URL.
// It includes the job name, status, and optional error details in the payload.
type HTTPPostNotifier struct {
	url    string
	client *http.Client
}

type httpPostPayload struct {
	JobName string `json:"job_name"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewHTTPPostNotifier creates an HTTPPostNotifier that posts to the given URL.
func NewHTTPPostNotifier(url string) *HTTPPostNotifier {
	return newHTTPPostNotifierWithClient(url, &http.Client{})
}

func newHTTPPostNotifierWithClient(url string, client *http.Client) *HTTPPostNotifier {
	return &HTTPPostNotifier{url: url, client: client}
}

// Notify sends a POST request with a JSON payload describing the job result.
func (n *HTTPPostNotifier) Notify(jobName string, err error) error {
	payload := httpPostPayload{
		JobName: jobName,
		Status:  "success",
	}
	if err != nil {
		payload.Status = "failure"
		payload.Message = err.Error()
	}

	body, encErr := json.Marshal(payload)
	if encErr != nil {
		return fmt.Errorf("httppost: failed to encode payload: %w", encErr)
	}

	resp, doErr := n.client.Post(n.url, "application/json", bytes.NewReader(body))
	if doErr != nil {
		return fmt.Errorf("httppost: request failed: %w", doErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("httppost: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
