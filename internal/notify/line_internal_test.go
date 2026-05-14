package notify

import (
	"strings"
	"testing"
)

func TestLineBody_NoError(t *testing.T) {
	body := lineBody("backup", nil)
	if !strings.Contains(body, "backup") {
		t.Errorf("expected job name in body, got %q", body)
	}
	if !strings.Contains(body, "successfully") {
		t.Errorf("expected success text in body, got %q", body)
	}
}

func TestLineBody_WithError(t *testing.T) {
	body := lineBody("backup", errTest("disk full"))
	if !strings.Contains(body, "backup") {
		t.Errorf("expected job name in body, got %q", body)
	}
	if !strings.Contains(body, "disk full") {
		t.Errorf("expected error text in body, got %q", body)
	}
	if !strings.Contains(body, "failed") {
		t.Errorf("expected 'failed' in body, got %q", body)
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }
