package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewStatuspageNotifier_SetsFields(t *testing.T) {
	n := notify.NewStatuspageNotifier("key123", "page1", "comp1")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestStatuspageNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "comp1") {
			t.Errorf("expected component ID in path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewStatuspageNotifierWithClient("key", "page1", "comp1", server.URL, server.Client())
	if err := n.Notify("myjob", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestStatuspageNotifier_Notify_WithError(t *testing.T) {
	var capturedBody strings.Builder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		capturedBody.WriteString(buf.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewStatuspageNotifierWithClient("key", "page1", "comp1", server.URL, server.Client())
	if err := n.Notify("myjob", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(capturedBody.String(), "major_outage") {
		t.Errorf("expected major_outage in payload, got %s", capturedBody.String())
	}
}

func TestStatuspageNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	n := notify.NewStatuspageNotifierWithClient("badkey", "page1", "comp1", server.URL, server.Client())
	if err := n.Notify("myjob", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestStatuspageNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewStatuspageNotifierWithClient("key", "page1", "comp1", "://bad-url", &http.Client{})
	if err := n.Notify("myjob", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
