package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const twilioAPIBase = "https://api.twilio.com/2010-04-01/Accounts"

// TwilioNotifier sends SMS alerts via the Twilio REST API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
	apiBase    string
}

// NewTwilioNotifier creates a TwilioNotifier that sends SMS messages.
func NewTwilioNotifier(accountSID, authToken, from, to string) *TwilioNotifier {
	return newTwilioNotifierWithClient(accountSID, authToken, from, to, &http.Client{}, twilioAPIBase)
}

func newTwilioNotifierWithClient(accountSID, authToken, from, to string, client *http.Client, apiBase string) *TwilioNotifier {
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     client,
		apiBase:    apiBase,
	}
}

// Notify sends an SMS with the job name and optional error via Twilio.
func (t *TwilioNotifier) Notify(jobName string, err error) error {
	body := fmt.Sprintf("cronwrap: job %q succeeded", jobName)
	if err != nil {
		body = fmt.Sprintf("cronwrap: job %q failed: %s", jobName, err.Error())
	}

	endpoint := fmt.Sprintf("%s/%s/Messages.json", t.apiBase, t.accountSID)

	form := url.Values{}
	form.Set("From", t.from)
	form.Set("To", t.to)
	form.Set("Body", body)

	req, err2 := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err2 != nil {
		return fmt.Errorf("twilio: build request: %w", err2)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err2 := t.client.Do(req)
	if err2 != nil {
		return fmt.Errorf("twilio: send request: %w", err2)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		return fmt.Errorf("twilio: unexpected status %d: %s", resp.StatusCode, apiErr.Message)
	}

	return nil
}
