package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/cronwrap/internal/notify"
)

func TestNewRocketChatNotifier_SetsURL(t *testing.T) {
	n := notify.NewRocketChatNotifier("https://chat.example.com/hooks/abc")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestRocketChatNotifier_Notify_Success(t *testing.T) {
	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewRocketChatNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(received, "backup-job") {
		t.Errorf("expected payload to contain job name, got: %s", received)
	}
	if !strings.Contains(received, "successfully") {
		t.Errorf("expected success message in payload, got: %s", received)
	}
}

func TestRocketChatNotifier_Notify_WithError(t *testing.T) {
	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewRocketChatNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("cleanup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(received, "disk full") {
		t.Errorf("expected error message in payload, got: %s", received)
	}
	if !strings.Contains(received, "cleanup-job") {
		t.Errorf("expected job name in payload, got: %s", received)
	}
}

func TestRocketChatNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewRocketChatNotifierWithClient(server.URL, server.Client())
	err := n.Notify("test-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected status code in error, got: %v", err)
	}
}

func TestRocketChatNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewRocketChatNotifier("http://127.0.0.1:0/invalid")
	err := n.Notify("test-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
