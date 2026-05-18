package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewOpsGenieNotifier_SetsKey(t *testing.T) {
	n := notify.NewOpsGenieNotifier("test-key")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestOpsGenieNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewOpsGenieNotifierWithURL("key", server.URL, server.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpsGenieNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewOpsGenieNotifierWithURL("key", server.URL, server.Client())
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpsGenieNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := notify.NewOpsGenieNotifierWithURL("key", server.URL, server.Client())
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestOpsGenieNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewOpsGenieNotifierWithURL("key", "http://127.0.0.1:0", &http.Client{})
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestOpsGenieNotifier_Notify_SetsAuthHeader(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewOpsGenieNotifierWithURL("my-secret-key", server.URL, server.Client())
	_ = n.Notify("myjob", nil)

	if gotAuth != "GenieKey my-secret-key" {
		t.Errorf("expected 'GenieKey my-secret-key', got %q", gotAuth)
	}
}
