package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/cronwrap/internal/notify"
)

func TestNewMatrixNotifier_SetsFields(t *testing.T) {
	n := notify.NewMatrixNotifier("https://matrix.example.com", "tok", "!room:example.com")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMatrixNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewMatrixNotifierWithClient(server.URL, "tok", "!room:example.com", server.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMatrixNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewMatrixNotifierWithClient(server.URL, "tok", "!room:example.com", server.Client())
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "disk full") {
		t.Errorf("expected body to contain error text, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "backup") {
		t.Errorf("expected body to contain job name, got: %s", gotBody)
	}
}

func TestMatrixNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := notify.NewMatrixNotifierWithClient(server.URL, "bad-token", "!room:example.com", server.Client())
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMatrixNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewMatrixNotifierWithClient("://bad-url", "tok", "!room:example.com", &http.Client{})
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestMatrixNotifier_AuthHeaderSet(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewMatrixNotifierWithClient(server.URL, "mytoken", "!room:example.com", server.Client())
	_ = n.Notify("job", nil)

	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Authorization header 'Bearer mytoken', got %q", gotAuth)
	}
}
