package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronwrap/internal/notify"
)

func TestNewDingTalkNotifier_SetsURL(t *testing.T) {
	n := notify.NewDingTalkNotifier("https://oapi.dingtalk.com/robot/send?access_token=abc")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDingTalkNotifier_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewDingTalkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDingTalkNotifier_Notify_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewDingTalkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup", errors.New("disk full")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDingTalkNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewDingTalkNotifierWithClient(server.URL, server.Client())
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestDingTalkNotifier_Notify_InvalidURL(t *testing.T) {
	n := notify.NewDingTalkNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify("backup", nil); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
