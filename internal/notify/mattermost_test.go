package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewMattermostNotifier_SetsURL(t *testing.T) {
	n := notify.NewMattermostNotifier("https://example.com/hook")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMattermostNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewMattermostNotifierWithClient(ts.URL, ts.Client())
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMattermostNotifier_Notify_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewMattermostNotifierWithClient(ts.URL, ts.Client())
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMattermostNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewMattermostNotifierWithClient(ts.URL, ts.Client())
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMattermostNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewMattermostNotifierWithClient("http://127.0.0.1:0", &http.Client{})
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
