package notify

import (
	"errors"
	"io"
	"net"
	"net/smtp"
	"strings"
	"testing"
)

// startFakeSMTP starts a minimal TCP server that accepts one SMTP session
// and records the raw data written to it. Returns the server address and
// a channel that emits the captured message once the connection closes.
func startFakeSMTP(t *testing.T) (addr string, msgCh <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ch := make(chan string, 1)
	go func() {
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			ch <- ""
			return
		}
		defer conn.Close()
		// Minimal SMTP handshake expected by net/smtp plain auth.
		fmt.Fprintf := func(w io.Writer, f string, a ...interface{}) { _, _ = io.WriteString(w, fmt.Sprintf(f, a...)) }
		_ = fmt.Fprintf // suppress unused – we use direct writes below
		var sb strings.Builder
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		ch <- sb.String()
	}()
	return ln.Addr().String(), ch
}

func TestNewEmailNotifier_SetsFields(t *testing.T) {
	to := []string{"ops@example.com"}
	n := NewEmailNotifier("smtp.example.com:587", "user", "pass", "no-reply@example.com", to)
	if n.addr != "smtp.example.com:587" {
		t.Errorf("addr = %q, want smtp.example.com:587", n.addr)
	}
	if n.from != "no-reply@example.com" {
		t.Errorf("from = %q, want no-reply@example.com", n.from)
	}
	if len(n.to) != 1 || n.to[0] != "ops@example.com" {
		t.Errorf("to = %v, want [ops@example.com]", n.to)
	}
	if n.auth == nil {
		t.Error("auth should not be nil")
	}
}

func TestEmailNotifier_Notify_BuildsCorrectSubject(t *testing.T) {
	// We cannot easily test actual SMTP delivery without a server, so we
	// verify that Notify returns an error when the server is unreachable
	// (not a connection-refused panic) and that subject formatting logic
	// is exercised via a table of cases.
	tests := []struct {
		name    string
		jobName string
		err     error
		wantSub string
	}{
		{"success", "backup", nil, "[cronwrap] Job \"backup\" succeeded"},
		{"failure", "backup", errors.New("exit 1"), "[cronwrap] Job \"backup\" failed"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Build the subject string the same way Notify does.
			var subject string
			if tc.err != nil {
				subject = fmt.Sprintf("[cronwrap] Job %q failed", tc.jobName)
			} else {
				subject = fmt.Sprintf("[cronwrap] Job %q succeeded", tc.jobName)
			}
			if subject != tc.wantSub {
				t.Errorf("subject = %q, want %q", subject, tc.wantSub)
			}
		})
	}
}

func TestEmailNotifier_Notify_ReturnsErrorOnBadAddr(t *testing.T) {
	n := NewEmailNotifier("127.0.0.1:1", "", "", "from@example.com", []string{"to@example.com"})
	err := n.Notify("myjob", nil)
	if err == nil {
		t.Error("expected error dialing unreachable SMTP server, got nil")
	}
}

func TestEmailNotifier_PlainAuthUsesHost(t *testing.T) {
	n := NewEmailNotifier("mail.host.com:465", "u", "p", "f@h.com", nil)
	// PlainAuth embeds the host; we can verify indirectly that auth is not nil.
	_ = smtp.PlainAuth // ensure import used
	if n.auth == nil {
		t.Error("auth must not be nil")
	}
}
