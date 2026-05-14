package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewTwilioNotifier_SetsFields(t *testing.T) {
	n := notify.NewTwilioNotifier("ACabc", "token123", "+10000000000", "+19999999999")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTwilioNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"sid":"SM123"}`))
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifierWithClient("ACabc", "token", "+1000", "+1999", ts.Client(), ts.URL)
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTwilioNotifier_Notify_WithError(t *testing.T) {
	var captured string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.FormValue("Body")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifierWithClient("ACabc", "token", "+1000", "+1999", ts.Client(), ts.URL)
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(captured, "disk full") {
		t.Errorf("expected body to contain error message, got %q", captured)
	}
	if !strings.Contains(captured, "backup") {
		t.Errorf("expected body to contain job name, got %q", captured)
	}
}

func TestTwilioNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid credentials"}`))
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifierWithClient("ACabc", "bad", "+1000", "+1999", ts.Client(), ts.URL)
	err := n.Notify("backup", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected status code in error, got %v", err)
	}
}

func TestTwilioNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewTwilioNotifierWithClient("ACabc", "token", "+1000", "+1999", &http.Client{}, "http://127.0.0.1:0")
	err := n.Notify("backup", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
