package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewGoogleChatNotifier_SetsURL(t *testing.T) {
	n := notify.NewGoogleChatNotifier("https://chat.googleapis.com/webhook")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGoogleChatNotifier_Notify_Success(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifierWithClient(ts.URL, ts.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "backup-job") {
		t.Errorf("expected body to contain job name, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "✅") {
		t.Errorf("expected success emoji in body, got: %s", gotBody)
	}
}

func TestGoogleChatNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifierWithClient(ts.URL, ts.Client())
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "disk full") {
		t.Errorf("expected error message in body, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "❌") {
		t.Errorf("expected failure emoji in body, got: %s", gotBody)
	}
}

func TestGoogleChatNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifierWithClient(ts.URL, ts.Client())
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected 403 in error, got: %v", err)
	}
}

func TestGoogleChatNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewGoogleChatNotifier("http://127.0.0.1:0/no-server")
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
