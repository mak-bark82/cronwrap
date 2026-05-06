package notify

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewPagerDutyNotifier_SetsKey(t *testing.T) {
	n := NewPagerDutyNotifier("test-key")
	if n.integrationKey != "test-key" {
		t.Errorf("expected integration key %q, got %q", "test-key", n.integrationKey)
	}
}

func TestPagerDutyNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key")
	n.client = ts.Client()
	// Override URL via a local helper for testability.
	err := notifyPD(n, ts.URL, "myjob", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPagerDutyNotifier_Notify_WithError(t *testing.T) {
	var capturedBody strings.Builder
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		capturedBody.WriteString(buf.String())
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key")
	n.client = ts.Client()
	notifyErr := errors.New("disk full")
	err := notifyPD(n, ts.URL, "backup", notifyErr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(capturedBody.String(), "disk full") {
		t.Errorf("expected body to contain error message, got: %s", capturedBody.String())
	}
}

func TestPagerDutyNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key")
	n.client = ts.Client()
	err := notifyPD(n, ts.URL, "myjob", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to mention status 400, got: %v", err)
	}
}

func TestPagerDutyNotifier_Notify_InvalidURL(t *testing.T) {
	n := NewPagerDutyNotifier("key")
	err := notifyPD(n, "http://127.0.0.1:0", "myjob", nil)
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

// notifyPD is a test helper that overrides the PagerDuty endpoint URL.
func notifyPD(n *PagerDutyNotifier, url, jobName string, err error) error {
	origURL := pagerDutyEventURL
	_ = origURL // kept for reference; we shadow via closure below
	// Temporarily swap client base URL by building the request manually.
	// We re-implement the core logic here to inject the test server URL.
	import_bytes := func() {}
	_ = import_bytes
	return notifyWithURL(n, url, jobName, err)
}
