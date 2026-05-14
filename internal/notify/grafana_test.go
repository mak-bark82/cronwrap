package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewGrafanaNotifier_SetsFields(t *testing.T) {
	n := notify.NewGrafanaNotifier("https://grafana.example.com", "secret")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGrafanaNotifier_Notify_Success(t *testing.T) {
	var gotAuth, gotContentType string
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGrafanaNotifierWithClient(ts.URL, "mykey", ts.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer mykey" {
		t.Errorf("expected Bearer mykey, got %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %q", gotContentType)
	}
	if !strings.Contains(gotBody, "backup") {
		t.Errorf("expected body to contain job name, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "success") {
		t.Errorf("expected body to contain 'success' tag, got %q", gotBody)
	}
}

func TestGrafanaNotifier_Notify_WithError(t *testing.T) {
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGrafanaNotifierWithClient(ts.URL, "mykey", ts.Client())
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "failure") {
		t.Errorf("expected body to contain 'failure' tag, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "disk full") {
		t.Errorf("expected body to contain error text, got %q", gotBody)
	}
}

func TestGrafanaNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewGrafanaNotifierWithClient(ts.URL, "badkey", ts.Client())
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGrafanaNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewGrafanaNotifierWithClient("http://127.0.0.1:0", "key", &http.Client{})
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
