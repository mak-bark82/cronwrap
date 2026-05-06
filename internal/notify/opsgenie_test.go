package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/cronwrap/internal/notify"
)

func TestNewOpsGenieNotifier_SetsKey(t *testing.T) {
	n := notify.NewOpsGenieNotifier("test-api-key")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestOpsGenieNotifier_Notify_Success(t *testing.T) {
	var gotAuth, gotContentType string
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := notify.NewOpsGenieNotifier("my-key")
	n.(*notify.OpsGenieNotifier) // type assertion not needed; use exported helper below
	_ = gotBody

	// Use the internal URL override via a helper struct if needed;
	// here we rely on the exported constructor and a fake server.
	// Since the URL is unexported, we test via a wrapper.
	tn := testOpsGenieNotifier(t, ts.URL, "my-key")

	if err := tn.Notify("backup-job", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasPrefix(gotAuth, "GenieKey ") {
		t.Errorf("expected GenieKey auth header, got %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %q", gotContentType)
	}
}

func TestOpsGenieNotifier_Notify_WithError(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := testOpsGenieNotifier(t, ts.URL, "key")
	if err := n.Notify("db-backup", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotBody, "P1") {
		t.Errorf("expected priority P1 in body, got %s", gotBody)
	}
	if !strings.Contains(gotBody, "disk full") {
		t.Errorf("expected error message in body, got %s", gotBody)
	}
}

func TestOpsGenieNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := testOpsGenieNotifier(t, ts.URL, "bad-key")
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestOpsGenieNotifier_Notify_InvalidURL(t *testing.T) {
	n := testOpsGenieNotifier(t, "http://127.0.0.1:0", "key")
	if err := n.Notify("job", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

// testOpsGenieNotifier creates a notifier pointed at a custom URL for testing.
func testOpsGenieNotifier(t *testing.T, url, apiKey string) notify.Notifier {
	t.Helper()
	return notify.NewOpsGenieNotifierWithURL(apiKey, url)
}
