package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewLarkNotifier_SetsURL(t *testing.T) {
	n := notify.NewLarkNotifier("https://open.larksuite.com/webhook/abc")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLarkNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewLarkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLarkNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewLarkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLarkNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewLarkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestLarkNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewLarkNotifier("://invalid-url")
	if err := n.Notify("backup-job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
