package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ibrahimker/cronwrap/internal/alert"
	"github.com/ibrahimker/cronwrap/internal/notify"
)

func TestNewZulipNotifier_SetsFields(t *testing.T) {
	n := notify.NewZulipNotifier("https://org.zulipchat.com", "bot@org", "secret", "general", "cron")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestZulipNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := newZulipNotifier(server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "backup", Err: nil})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestZulipNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := newZulipNotifier(server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "backup", Err: errors.New("exit status 1")})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestZulipNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	n := newZulipNotifier(server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "backup"})
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestZulipNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewZulipNotifier("://bad-url", "bot@org", "key", "general", "cron")
	err := n.Notify(alert.Event{JobName: "backup"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

// newZulipNotifier is a test helper that injects a custom HTTP client.
func newZulipNotifier(baseURL string, client *http.Client) *notify.ZulipNotifier {
	return notify.NewZulipNotifierWithClient(baseURL, "bot@org", "secret", "general", "cron", client)
}
