package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func newTestLogger(jobName string) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return New(jobName, &buf), &buf
}

func TestLogger_Info(t *testing.T) {
	l, buf := newTestLogger("test-job")
	l.Info("hello world")
	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected INFO level, got: %s", out)
	}
	if !strings.Contains(out, "job=test-job") {
		t.Errorf("expected job name in output, got: %s", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestLogger_Warn(t *testing.T) {
	l, buf := newTestLogger("warn-job")
	l.Warn("something odd")
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Errorf("expected WARN level, got: %s", buf.String())
	}
}

func TestLogger_Error(t *testing.T) {
	l, buf := newTestLogger("err-job")
	l.Error("critical failure")
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected ERROR level, got: %s", buf.String())
	}
}

func TestLogger_JobStarted(t *testing.T) {
	l, buf := newTestLogger("myjob")
	l.JobStarted(2)
	out := buf.String()
	if !strings.Contains(out, "attempt 2") {
		t.Errorf("expected attempt number in output, got: %s", out)
	}
}

func TestLogger_JobSucceeded(t *testing.T) {
	l, buf := newTestLogger("myjob")
	l.JobSucceeded(1, 350*time.Millisecond)
	out := buf.String()
	if !strings.Contains(out, "succeeded") {
		t.Errorf("expected 'succeeded' in output, got: %s", out)
	}
	if !strings.Contains(out, "350ms") {
		t.Errorf("expected duration in output, got: %s", out)
	}
}

func TestLogger_JobFailed(t *testing.T) {
	l, buf := newTestLogger("myjob")
	l.JobFailed(1, 100*time.Millisecond, errors.New("exit status 1"))
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output, got: %s", out)
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestLogger_JobExhausted(t *testing.T) {
	l, buf := newTestLogger("myjob")
	l.JobExhausted(3, errors.New("timeout"))
	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected ERROR level for exhausted, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected attempt count in output, got: %s", out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	l := New("job", nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}
