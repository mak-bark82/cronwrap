package notify_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewHTTPPostNotifier_SetsURL(t *testing.T) {
	n := notify.NewHTTPPostNotifier("http://example.com/hook")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestHTTPPostNotifier_Notify_Success(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewHTTPPostNotifier(ts.URL)
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received["status"] != "success" {
		t.Errorf("expected status=success, got %q", received["status"])
	}
	if received["job_name"] != "backup-job" {
		t.Errorf("expected job_name=backup-job, got %q", received["job_name"])
	}
}

func TestHTTPPostNotifier_Notify_WithError(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewHTTPPostNotifier(ts.URL)
	if err := n.Notify("cleanup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received["status"] != "failure" {
		t.Errorf("expected status=failure, got %q", received["status"])
	}
	if received["message"] != "disk full" {
		t.Errorf("expected message=disk full, got %q", received["message"])
	}
}

func TestHTTPPostNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewHTTPPostNotifier(ts.URL)
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestHTTPPostNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewHTTPPostNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
