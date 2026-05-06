package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/cronwrap/internal/notify"
)

func TestNewTelegramNotifier_SetsFields(t *testing.T) {
	n := notify.NewTelegramNotifier("mytoken", "123456")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTelegramNotifier_Notify_Success(t *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		capturedBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewTelegramNotifierWithBase("token", "42", server.URL)
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(capturedBody, "backup-job") {
		t.Errorf("expected body to contain job name, got: %s", capturedBody)
	}
	if !strings.Contains(capturedBody, "42") {
		t.Errorf("expected body to contain chat_id, got: %s", capturedBody)
	}
}

func TestTelegramNotifier_Notify_WithError(t *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		capturedBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewTelegramNotifierWithBase("token", "42", server.URL)
	if err := n.Notify("cleanup-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(capturedBody, "disk full") {
		t.Errorf("expected body to mention the error, got: %s", capturedBody)
	}
}

func TestTelegramNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	n := notify.NewTelegramNotifierWithBase("badtoken", "42", server.URL)
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestTelegramNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewTelegramNotifierWithBase("token", "42", "http://127.0.0.1:0")
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
