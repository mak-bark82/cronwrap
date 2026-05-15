package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewNewRelicNotifier_SetsKey(t *testing.T) {
	n := notify.NewNewRelicNotifier("test-key")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewRelicNotifier_Notify_Success(t *testing.T) {
	var gotAPIKey, gotContentType string
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("Api-Key")
		gotContentType = r.Header.Get("Content-Type")
		gotBody = make([]byte, r.ContentLength)
		r.Body.Read(gotBody)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := notify.NewNewRelicNotifierWithClient("my-api-key", ts.URL, ts.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAPIKey != "my-api-key" {
		t.Errorf("expected Api-Key header %q, got %q", "my-api-key", gotAPIKey)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", gotContentType)
	}
	if len(gotBody) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestNewRelicNotifier_Notify_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := notify.NewNewRelicNotifierWithClient("key", ts.URL, ts.Client())
	if err := n.Notify("cleanup-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRelicNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewNewRelicNotifierWithClient("key", ts.URL, ts.Client())
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestNewRelicNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewNewRelicNotifierWithClient("key", "://bad-url", &http.Client{})
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
