package notify

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpsGenieNotify_PayloadContainsJobName(t *testing.T) {
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := newOpsGenieNotifierWithURL("key", server.URL, server.Client())
	_ = n.Notify("cleanup", nil)

	msg, _ := got["message"].(string)
	if msg == "" {
		t.Fatal("expected message field in payload")
	}
	alias, _ := got["alias"].(string)
	if alias != "cronwrap-cleanup" {
		t.Errorf("expected alias 'cronwrap-cleanup', got %q", alias)
	}
}

func TestOpsGenieNotify_P1PriorityOnError(t *testing.T) {
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := newOpsGenieNotifierWithURL("key", server.URL, server.Client())
	_ = n.Notify("cleanup", errors.New("timeout"))

	priority, _ := got["priority"].(string)
	if priority != "P1" {
		t.Errorf("expected priority P1 on error, got %q", priority)
	}
}

func TestOpsGenieNotify_P5PriorityOnSuccess(t *testing.T) {
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n := newOpsGenieNotifierWithURL("key", server.URL, server.Client())
	_ = n.Notify("cleanup", nil)

	priority, _ := got["priority"].(string)
	if priority != "P5" {
		t.Errorf("expected priority P5 on success, got %q", priority)
	}
}
