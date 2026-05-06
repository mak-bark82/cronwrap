package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewTeamsNotifier_SetsURL(t *testing.T) {
	n := notify.NewTeamsNotifier("https://example.webhook.office.com/abc")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTeamsNotifier_Notify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewTeamsNotifier(ts.URL)
	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTeamsNotifier_Notify_WithError(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody = make([]byte, r.ContentLength)
		r.Body.Read(gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewTeamsNotifier(ts.URL)
	jobErr := errors.New("exit status 1")
	if err := n.Notify("cleanup-job", jobErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(gotBody)
	if len(body) == 0 {
		t.Fatal("expected non-empty request body")
	}
	for _, want := range []string{"cleanup-job", "exit status 1", "FF0000"} {
		if !containsString(body, want) {
			t.Errorf("body missing %q; body=%s", want, body)
		}
	}
}

func TestTeamsNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := notify.NewTeamsNotifier(ts.URL)
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestTeamsNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewTeamsNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify("my-job", nil); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestTeamsNotifier_Notify_SuccessPayloadGreen(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody = make([]byte, r.ContentLength)
		r.Body.Read(gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewTeamsNotifier(ts.URL)
	if err := n.Notify("daily-report", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(gotBody)
	if !containsString(body, "00FF00") {
		t.Errorf("expected green color in success payload; body=%s", body)
	}
	if !containsString(body, "daily-report") {
		t.Errorf("expected job name in payload; body=%s", body)
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
