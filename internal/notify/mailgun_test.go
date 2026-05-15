package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func newMailgunNotifierWithClient(domain, apiKey, from, to string, client *http.Client, baseURL string) *notify.MailgunNotifier {
	return notify.NewMailgunNotifierWithClient(domain, apiKey, from, to, client, baseURL)
}

func TestNewMailgunNotifier_SetsFields(t *testing.T) {
	n := notify.NewMailgunNotifier("mg.example.com", "key-abc", "from@example.com", "to@example.com")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMailgunNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewMailgunNotifierWithClient("mg.example.com", "key-abc", "from@example.com", "to@example.com", ts.Client(), ts.URL)
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMailgunNotifier_Notify_WithError(t *testing.T) {
	var capturedBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		capturedBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewMailgunNotifierWithClient("mg.example.com", "key-abc", "from@example.com", "to@example.com", ts.Client(), ts.URL)
	if err := n.Notify("backup-job", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(capturedBody, "FAILED") && !strings.Contains(capturedBody, "disk+full") {
		// URL-encoded body should contain job failure info
		if !strings.Contains(capturedBody, "failed") {
			t.Errorf("expected body to mention failure, got: %s", capturedBody)
		}
	}
}

func TestMailgunNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := notify.NewMailgunNotifierWithClient("mg.example.com", "key-abc", "from@example.com", "to@example.com", ts.Client(), ts.URL)
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to contain status code, got: %v", err)
	}
}

func TestMailgunNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewMailgunNotifierWithClient("mg.example.com", "key-abc", "from@example.com", "to@example.com", &http.Client{}, "http://127.0.0.1:0")
	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
