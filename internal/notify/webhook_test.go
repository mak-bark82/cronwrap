package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestWebhookNotifier_Notify_Success(t *testing.T) {
	var received notify.WebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewWebhookNotifier(server.URL)
	payload := notify.WebhookPayload{
		JobName:  "backup",
		Status:   "success",
		Duration: 2 * time.Second,
		Attempts: 1,
	}

	if err := n.Notify(payload); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if received.JobName != "backup" {
		t.Errorf("expected job_name 'backup', got '%s'", received.JobName)
	}
	if received.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", received.Status)
	}
}

func TestWebhookNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewWebhookNotifier(server.URL)
	err := n.Notify(notify.WebhookPayload{JobName: "test", Status: "failure"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewWebhookNotifier("http://127.0.0.1:0/no-server")
	err := n.Notify(notify.WebhookPayload{JobName: "test", Status: "failure"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

func TestNewWebhookNotifier_SetsURL(t *testing.T) {
	url := "https://example.com/hook"
	n := notify.NewWebhookNotifier(url)
	if n.URL != url {
		t.Errorf("expected URL %q, got %q", url, n.URL)
	}
	if n.Client == nil {
		t.Error("expected non-nil HTTP client")
	}
}
