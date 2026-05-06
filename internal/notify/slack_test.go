package notify

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewSlackNotifier_SetsURL(t *testing.T) {
	n := NewSlackNotifier("https://hooks.slack.com/test")
	if n.webhookURL != "https://hooks.slack.com/test" {
		t.Errorf("expected webhookURL to be set, got %q", n.webhookURL)
	}
	if n.client == nil {
		t.Error("expected http client to be initialised")
	}
}

func TestSlackNotifier_Notify_Success(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "backup") {
		t.Errorf("expected body to contain job name, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "succeeded") {
		t.Errorf("expected body to contain 'succeeded', got %q", gotBody)
	}
}

func TestSlackNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	if err := n.Notify("cleanup", errors.New("exit status 1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "failed") {
		t.Errorf("expected body to contain 'failed', got %q", gotBody)
	}
	if !strings.Contains(gotBody, "exit status 1") {
		t.Errorf("expected body to contain error message, got %q", gotBody)
	}
}

func TestSlackNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	err := n.Notify("job", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got %v", err)
	}
}

func TestSlackNotifier_Notify_InvalidURL(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:0/no-server")
	err := n.Notify("job", nil)
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
