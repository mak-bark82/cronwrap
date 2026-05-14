package notify

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignalSciencesNotify_PayloadContainsJobName(t *testing.T) {
	var captured signalSciencesPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := newSignalSciencesNotifierWithClient("corp", "site", "tok", server.Client())
	n.baseURL = server.URL

	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.JobName != "backup-job" {
		t.Errorf("expected job_name %q, got %q", "backup-job", captured.JobName)
	}
	if captured.Status != "success" {
		t.Errorf("expected status success, got %q", captured.Status)
	}
}

func TestSignalSciencesNotify_FailureStatusOnError(t *testing.T) {
	var captured signalSciencesPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := newSignalSciencesNotifierWithClient("corp", "site", "tok", server.Client())
	n.baseURL = server.URL

	if err := n.Notify("cleanup-job", errors.New("disk full")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.Status != "failure" {
		t.Errorf("expected status failure, got %q", captured.Status)
	}
}

func TestSignalSciencesNotify_NonOKStatusReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n := newSignalSciencesNotifierWithClient("corp", "site", "tok", server.Client())
	n.baseURL = server.URL

	err := n.Notify("my-job", nil)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSignalSciencesNotify_SetsAuthHeaders(t *testing.T) {
	var gotToken, gotCorp string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("x-api-token")
		gotCorp = r.Header.Get("x-api-user")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := newSignalSciencesNotifierWithClient("mycorp", "mysite", "secret-token", server.Client())
	n.baseURL = server.URL

	if err := n.Notify("auth-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotToken != "secret-token" {
		t.Errorf("expected token %q, got %q", "secret-token", gotToken)
	}
	if gotCorp != "mycorp" {
		t.Errorf("expected corp %q, got %q", "mycorp", gotCorp)
	}
}
