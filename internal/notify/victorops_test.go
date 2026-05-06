package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwrap/cronwrap/internal/notify"
)

func TestNewVictorOpsNotifier_SetsFields(t *testing.T) {
	n := notify.NewVictorOpsNotifier("api-key-123", "routing-key-456")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestVictorOpsNotifier_Notify_Success(t *testing.T) {
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody = make([]byte, r.ContentLength)
		r.Body.Read(gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewVictorOpsNotifier("key", "route")
	n.(*notify.VictorOpsNotifier).BaseURL = server.URL

	if err := n.Notify("my-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotBody) == 0 {
		t.Error("expected body to be sent")
	}
}

func TestVictorOpsNotifier_Notify_WithError(t *testing.T) {
	var msgType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		msgType, _ = payload["message_type"].(string)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewVictorOpsNotifier("key", "route")
	n.(*notify.VictorOpsNotifier).BaseURL = server.URL

	if err := n.Notify("my-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msgType != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %q", msgType)
	}
}

func TestVictorOpsNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := notify.NewVictorOpsNotifier("key", "route")
	n.(*notify.VictorOpsNotifier).BaseURL = server.URL

	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestVictorOpsNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewVictorOpsNotifier("key", "route")
	n.(*notify.VictorOpsNotifier).BaseURL = "http://127.0.0.1:0"

	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
