package alert_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/alert"
)

func TestNew_SetsFields(t *testing.T) {
	before := time.Now()
	err := errors.New("something broke")
	a := alert.New("backup-job", alert.LevelError, "job failed", err)
	after := time.Now()

	if a.JobName != "backup-job" {
		t.Errorf("expected job name %q, got %q", "backup-job", a.JobName)
	}
	if a.Level != alert.LevelError {
		t.Errorf("expected level %q, got %q", alert.LevelError, a.Level)
	}
	if a.Message != "job failed" {
		t.Errorf("expected message %q, got %q", "job failed", a.Message)
	}
	if a.Err != err {
		t.Errorf("expected err %v, got %v", err, a.Err)
	}
	if a.OccuredAt.Before(before) || a.OccuredAt.After(after) {
		t.Errorf("OccuredAt %v out of expected range", a.OccuredAt)
	}
}

func TestStderrNotifier_Notify_NoError(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.StderrNotifier{Writer: &buf}
	a := alert.New("cleanup", alert.LevelWarn, "disk usage high", nil)

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error from Notify: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "level=WARN") {
		t.Errorf("expected level=WARN in output, got: %s", out)
	}
	if !strings.Contains(out, `job="cleanup"`) {
		t.Errorf("expected job=\"cleanup\" in output, got: %s", out)
	}
	if !strings.Contains(out, `message="disk usage high"`) {
		t.Errorf("expected message in output, got: %s", out)
	}
	if strings.Contains(out, "error=") {
		t.Errorf("did not expect error= field when err is nil, got: %s", out)
	}
}

func TestStderrNotifier_Notify_WithError(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.StderrNotifier{Writer: &buf}
	a := alert.New("sync", alert.LevelError, "sync failed", errors.New("timeout"))

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error from Notify: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `error="timeout"`) {
		t.Errorf("expected error=\"timeout\" in output, got: %s", out)
	}
}

func TestNewStderrNotifier_WritesToStderr(t *testing.T) {
	n := alert.NewStderrNotifier()
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.Writer == nil {
		t.Fatal("expected non-nil writer")
	}
}
