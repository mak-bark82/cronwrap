package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewSplunkNotifier_SetsToken(t *testing.T) {
	n := notify.NewSplunkNotifier("mytoken")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSplunkNotifier_Notify_Success(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewSplunkNotifierWithClient("tok123", server.URL, server.Client())
	if err := n.Notify("my-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotAuth != "Splunk tok123" {
		t.Errorf("unexpected Authorization header: %q", gotAuth)
	}
}

func TestSplunkNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewSplunkNotifierWithClient("tok", server.URL, server.Client())
	if err := n.Notify("failing-job", errors.New("exit status 1")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSplunkNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := notify.NewSplunkNotifierWithClient("tok", server.URL, server.Client())
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestSplunkNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewSplunkNotifierWithClient("tok", "://bad-url", &http.Client{})
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
