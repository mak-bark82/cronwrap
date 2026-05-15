package notify

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// MailgunNotifier sends alert notifications via the Mailgun API.
type MailgunNotifier struct {
	domain string
	apiKey string
	from   string
	to     string
	client *http.Client
	baseURL string
}

// NewMailgunNotifier creates a MailgunNotifier that sends email via Mailgun.
// domain is your Mailgun sending domain, apiKey is your Mailgun private API key,
// from and to are the sender and recipient email addresses.
func NewMailgunNotifier(domain, apiKey, from, to string) *MailgunNotifier {
	return newMailgunNotifierWithClient(domain, apiKey, from, to, &http.Client{}, "https://api.mailgun.net")
}

func newMailgunNotifierWithClient(domain, apiKey, from, to string, client *http.Client, baseURL string) *MailgunNotifier {
	return &MailgunNotifier{
		domain:  domain,
		apiKey:  apiKey,
		from:    from,
		to:      to,
		client:  client,
		baseURL: baseURL,
	}
}

// Notify sends an email via Mailgun reporting the job result.
func (n *MailgunNotifier) Notify(job string, err error) error {
	subject, body := mailgunMessage(job, err)

	endpoint := fmt.Sprintf("%s/v3/%s/messages", n.baseURL, n.domain)

	form := url.Values{}
	form.Set("from", n.from)
	form.Set("to", n.to)
	form.Set("subject", subject)
	form.Set("text", body)

	req, reqErr := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if reqErr != nil {
		return fmt.Errorf("mailgun: build request: %w", reqErr)
	}
	req.SetBasicAuth("api", n.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, doErr := n.client.Do(req)
	if doErr != nil {
		return fmt.Errorf("mailgun: send request: %w", doErr)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mailgun: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func mailgunMessage(job string, err error) (subject, body string) {
	if err != nil {
		subject = fmt.Sprintf("[cronwrap] FAILED: %s", job)
		body = fmt.Sprintf("Job %q failed with error: %v", job, err)
	} else {
		subject = fmt.Sprintf("[cronwrap] OK: %s", job)
		body = fmt.Sprintf("Job %q completed successfully.", job)
	}
	return
}
