package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultSignalSciencesURL = "https://dashboard.signalsciences.net/api/v0/corps/%s/sites/%s/feed/requests"

// SignalSciencesNotifier sends job alerts to Signal Sciences (Fastly Next-Gen WAF)
// via their event feed API.
type SignalSciencesNotifier struct {
	corpName string
	siteName string
	apiToken string
	baseURL  string
	client   *http.Client
}

type signalSciencesPayload struct {
	Event   string `json:"event"`
	JobName string `json:"job_name"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewSignalSciencesNotifier creates a new SignalSciencesNotifier.
func NewSignalSciencesNotifier(corpName, siteName, apiToken string) *SignalSciencesNotifier {
	return newSignalSciencesNotifierWithClient(corpName, siteName, apiToken, &http.Client{})
}

func newSignalSciencesNotifierWithClient(corpName, siteName, apiToken string, client *http.Client) *SignalSciencesNotifier {
	return &SignalSciencesNotifier{
		corpName: corpName,
		siteName: siteName,
		apiToken: apiToken,
		baseURL:  fmt.Sprintf(defaultSignalSciencesURL, corpName, siteName),
		client:   client,
	}
}

// Notify sends a job event notification to Signal Sciences.
func (n *SignalSciencesNotifier) Notify(jobName string, err error) error {
	status := "success"
	msg := fmt.Sprintf("cronwrap: job %q completed successfully", jobName)
	if err != nil {
		status = "failure"
		msg = fmt.Sprintf("cronwrap: job %q failed: %v", jobName, err)
	}

	payload := signalSciencesPayload{
		Event:   "cronwrap",
		JobName: jobName,
		Message: msg,
		Status:  status,
	}

	body, err2 := json.Marshal(payload)
	if err2 != nil {
		return fmt.Errorf("signalsciences: marshal payload: %w", err2)
	}

	req, err2 := http.NewRequest(http.MethodPost, n.baseURL, bytes.NewReader(body))
	if err2 != nil {
		return fmt.Errorf("signalsciences: create request: %w", err2)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-user", n.corpName)
	req.Header.Set("x-api-token", n.apiToken)

	resp, err2 := n.client.Do(req)
	if err2 != nil {
		return fmt.Errorf("signalsciences: send request: %w", err2)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalsciences: unexpected status %d", resp.StatusCode)
	}
	return nil
}
