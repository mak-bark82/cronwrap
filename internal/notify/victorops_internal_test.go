package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVictorOpsNotifyWithURL_PayloadContainsJobName(t *testing.T) {
	var payload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	v := NewVictorOpsNotifier("key", "route")
	if err := v.notifyWithURL(server.URL, "backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entityID, _ := payload["entity_id"].(string)
	if !strings.Contains(entityID, "backup-job") {
		t.Errorf("entity_id %q does not contain job name", entityID)
	}
}

func TestVictorOpsNotifyWithURL_CriticalOnError(t *testing.T) {
	var payload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	v := NewVictorOpsNotifier("key", "route")
	if err := v.notifyWithURL(server.URL, "db-backup", errors.New("timeout")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msgType, _ := payload["message_type"].(string)
	if msgType != "CRITICAL" {
		t.Errorf("expected CRITICAL message_type, got %q", msgType)
	}
}

func TestVictorOpsNotifyWithURL_InfoOnSuccess(t *testing.T) {
	var payload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	v := NewVictorOpsNotifier("key", "route")
	if err := v.notifyWithURL(server.URL, "db-backup", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msgType, _ := payload["message_type"].(string)
	if msgType != "INFO" {
		t.Errorf("expected INFO message_type, got %q", msgType)
	}
}
