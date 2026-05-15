package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewHipChatNotifier_SetsFields(t *testing.T) {
	n := notify.NewHipChatNotifier("tok", "room42")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestHipChatNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mytoken" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := notify.NewHipChatNotifierWithClient("mytoken", "room1", ts.URL+"/%s", ts.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHipChatNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 512)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := notify.NewHipChatNotifierWithClient("tok", "room1", ts.URL+"/%s", ts.Client())
	if err := n.Notify("cleanup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody == "" {
		t.Fatal("expected non-empty request body")
	}
}

func TestHipChatNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewHipChatNotifierWithClient("badtoken", "room1", ts.URL+"/%s", ts.Client())
	err := n.Notify("job", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestHipChatNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewHipChatNotifierWithClient("tok", "room", "://bad-url/%s", &http.Client{})
	err := n.Notify("job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
