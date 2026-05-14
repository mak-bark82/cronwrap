package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exampleorg/cronwrap/internal/alert"
	"github.com/exampleorg/cronwrap/internal/notify"
)

func TestNewCircleCINotifier_SetsToken(t *testing.T) {
	n := notify.NewCircleCINotifier("my-token")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCircleCINotifier_Notify_Success(t *testing.T) {
	var gotToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("Circle-Token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewCircleCINotifierWithClient("tok123", server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "backup", Err: nil})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotToken != "tok123" {
		t.Errorf("expected token tok123, got %q", gotToken)
	}
}

func TestCircleCINotifier_Notify_WithError(t *testing.T) {
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := notify.NewCircleCINotifierWithClient("tok", server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "sync", Err: errors.New("timeout")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(gotBody, []byte("timeout")) {
		t.Errorf("expected body to contain error text, got %s", gotBody)
	}
}

func TestCircleCINotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	n := notify.NewCircleCINotifierWithClient("bad-tok", server.URL, server.Client())
	err := n.Notify(alert.Event{JobName: "job", Err: nil})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestCircleCINotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewCircleCINotifierWithClient("tok", "://bad-url", &http.Client{})
	err := n.Notify(alert.Event{JobName: "job", Err: nil})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
