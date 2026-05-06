package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewGotifyNotifier_SetsFields(t *testing.T) {
	n := notify.NewGotifyNotifier("http://gotify.example.com", "mytoken")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGotifyNotifier_Notify_Success(t *testing.T) {
	var received string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.URL.Query().Get("token")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGotifyNotifier(ts.URL, "testtoken")
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received != "testtoken" {
		t.Errorf("expected token %q, got %q", "testtoken", received)
	}
}

func TestGotifyNotifier_Notify_WithError(t *testing.T) {
	var body strings.Builder
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 512)
		n, _ := r.Body.Read(buf)
		body.Write(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGotifyNotifier(ts.URL, "tok")
	jobErr := errors.New("exit status 1")
	if err := n.Notify("cleanup", jobErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body.String(), "exit status 1") {
		t.Errorf("expected error in body, got: %s", body.String())
	}
}

func TestGotifyNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewGotifyNotifier(ts.URL, "badtoken")
	err := n.Notify("myjob", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 in error, got: %v", err)
	}
}

func TestGotifyNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewGotifyNotifier("http://127.0.0.1:0", "tok")
	err := n.Notify("myjob", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
