package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwrap/internal/notify"
)

func TestNewDatadogNotifier_SetsKey(t *testing.T) {
	n := notify.NewDatadogNotifier("my-api-key")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDatadogNotifier_Notify_Success(t *testing.T) {
	var gotAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("DD-API-KEY")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewDatadogNotifierWithURL("test-key", server.URL)
	if err := n.Notify("my-job", nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if gotAPIKey != "test-key" {
		t.Errorf("expected API key %q, got %q", "test-key", gotAPIKey)
	}
}

func TestDatadogNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewDatadogNotifierWithURL("key", server.URL)
	if err := n.Notify("fail-job", errors.New("timeout")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody == "" {
		t.Error("expected non-empty body")
	}
}

func TestDatadogNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := notify.NewDatadogNotifierWithURL("key", server.URL)
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDatadogNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewDatadogNotifierWithURL("key", "://bad-url")
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
