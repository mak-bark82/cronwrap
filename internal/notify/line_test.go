package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwrap/internal/notify"
)

func TestNewLineNotifier_SetsToken(t *testing.T) {
	n := notify.NewLineNotifier("tok123")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLineNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewLineNotifierWithClient("tok", ts.URL, ts.Client())
	if err := n.Notify("myjob", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLineNotifier_Notify_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewLineNotifierWithClient("tok", ts.URL, ts.Client())
	if err := n.Notify("myjob", errors.New("boom")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLineNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewLineNotifierWithClient("tok", ts.URL, ts.Client())
	if err := n.Notify("myjob", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestLineNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewLineNotifierWithClient("tok", "://bad-url", http.DefaultClient)
	if err := n.Notify("myjob", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
