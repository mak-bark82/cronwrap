package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/celest-dev/cronwrap/internal/notify"
)

func TestNewNtfyNotifier_SetsFields(t *testing.T) {
	n := notify.NewNtfyNotifier("https://ntfy.sh", "alerts")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNtfyNotifier_Notify_Success(t *testing.T) {
	var gotTitle, gotPriority string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTitle = r.Header.Get("Title")
		gotPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewNtfyNotifierWithClient(ts.URL, "alerts", ts.Client())
	if err := n.Notify("my-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(gotTitle, "succeeded") {
		t.Errorf("expected title to contain 'succeeded', got %q", gotTitle)
	}
	if gotPriority != "default" {
		t.Errorf("expected priority 'default', got %q", gotPriority)
	}
}

func TestNtfyNotifier_Notify_WithError(t *testing.T) {
	var gotTitle, gotPriority string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTitle = r.Header.Get("Title")
		gotPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewNtfyNotifierWithClient(ts.URL, "alerts", ts.Client())
	if err := n.Notify("my-job", errors.New("exit status 1")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(gotTitle, "FAILED") {
		t.Errorf("expected title to contain 'FAILED', got %q", gotTitle)
	}
	if gotPriority != "high" {
		t.Errorf("expected priority 'high', got %q", gotPriority)
	}
}

func TestNtfyNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewNtfyNotifierWithClient(ts.URL, "alerts", ts.Client())
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got %v", err)
	}
}

func TestNtfyNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewNtfyNotifier("http://127.0.0.1:0", "alerts")
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
