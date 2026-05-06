package notify

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailNotifier sends alert notifications via SMTP email.
type EmailNotifier struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
}

// NewEmailNotifier creates an EmailNotifier using plain SMTP auth.
// addr should be in the form "host:port" (e.g. "smtp.example.com:587").
func NewEmailNotifier(addr, username, password, from string, to []string) *EmailNotifier {
	host := strings.Split(addr, ":")[0]
	auth := smtp.PlainAuth("", username, password, host)
	return &EmailNotifier{
		addr: addr,
		auth: auth,
		from: from,
		to:   to,
	}
}

// Notify sends an email with the job name and optional error detail.
func (e *EmailNotifier) Notify(jobName string, err error) error {
	subject := fmt.Sprintf("[cronwrap] Job %q succeeded", jobName)
	body := fmt.Sprintf("Job: %s\nStatus: success", jobName)

	if err != nil {
		subject = fmt.Sprintf("[cronwrap] Job %q failed", jobName)
		body = fmt.Sprintf("Job: %s\nStatus: failed\nError: %v", jobName, err)
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		e.from,
		strings.Join(e.to, ", "),
		subject,
		body,
	))

	return smtp.SendMail(e.addr, e.auth, e.from, e.to, msg)
}
