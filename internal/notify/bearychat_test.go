package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewBearyChat_SetsURL(t *testing.T) {
	n := notify.NewBearyChat("https://hook.bearychat.com/abc", &http.Client{})
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBearyChat_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewBearyChat(server.URL, server.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearyChat_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewBearyChat(server.URL, server.Client())
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearyChat_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewBearyChat(server.URL, server.Client())
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestBearyChat_Notify_InvalidURL(t *testing.T) {
	n := notify.NewBearyChat("http://127.0.0.1:0", &http.Client{})
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
