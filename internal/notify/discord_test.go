package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewDiscordNotifier_SetsURL(t *testing.T) {
	n := notify.NewDiscordNotifier("https://discord.com/api/webhooks/123/abc")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDiscordNotifier_Notify_Success(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n := notify.NewDiscordNotifier(server.URL)
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(gotBody, "backup-job") {
		t.Errorf("expected body to contain job name, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "succeeded") {
		t.Errorf("expected body to indicate success, got: %s", gotBody)
	}
}

func TestDiscordNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n := notify.NewDiscordNotifier(server.URL)
	jobErr := errors.New("exit status 1")
	if err := n.Notify("backup-job", jobErr); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(gotBody, "failed") {
		t.Errorf("expected body to indicate failure, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "exit status 1") {
		t.Errorf("expected body to contain error text, got: %s", gotBody)
	}
}

func TestDiscordNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewDiscordNotifier(server.URL)
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got: %v", err)
	}
}

func TestDiscordNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewDiscordNotifier("http://127.0.0.1:0")
	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
