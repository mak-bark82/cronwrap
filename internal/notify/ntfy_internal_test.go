package notify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNtfyNotify_TopicAppearsInURL(t *testing.T) {
	var gotPath string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := newNtfyNotifierWithClient(ts.URL, "my-topic", ts.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotPath, "my-topic") {
		t.Errorf("expected path to contain topic, got %q", gotPath)
	}
}

func TestNtfyNotify_BodyContainsJobName(t *testing.T) {
	var gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := newNtfyNotifierWithClient(ts.URL, "alerts", ts.Client())
	if err := n.Notify("nightly-backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotBody, "nightly-backup") {
		t.Errorf("expected body to contain job name, got %q", gotBody)
	}
}
