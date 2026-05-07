package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwrap/cronwrap/internal/notify"
)

func TestNewPushoverNotifier_SetsFields(t *testing.T) {
	n := notify.NewPushoverNotifier("tok123", "user456")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestPushoverNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewPushoverNotifierWithClient("tok", "user", server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPushoverNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewPushoverNotifierWithClient("tok", "user", server.URL, server.Client())
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPushoverNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	n := notify.NewPushoverNotifierWithClient("tok", "user", server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestPushoverNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewPushoverNotifierWithClient("tok", "user", "http://127.0.0.1:0", &http.Client{})
	if err := n.Notify("backup-job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
