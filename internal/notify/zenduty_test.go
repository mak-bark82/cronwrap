package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewZendutyNotifier_SetsKey(t *testing.T) {
	n := notify.NewZendutyNotifier("key-123")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestZendutyNotifier_Notify_Success(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	n := notify.NewZendutyNotifierWithClient("key-abc", svr.URL, svr.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestZendutyNotifier_Notify_WithError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	n := notify.NewZendutyNotifierWithClient("key-abc", svr.URL, svr.Client())
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestZendutyNotifier_Notify_NonOKStatus(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	n := notify.NewZendutyNotifierWithClient("key-abc", svr.URL, svr.Client())
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestZendutyNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewZendutyNotifierWithClient("key-abc", "http://127.0.0.1:0", &http.Client{})
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
