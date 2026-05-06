package notify

import (
	"errors"
	"strings"
	"testing"
)

// fakeNotifier records calls made to it.
type fakeNotifier struct {
	called  int
	jobName string
	err     error
	retErr  error // error to return from Notify
}

func (f *fakeNotifier) Notify(jobName string, err error) error {
	f.called++
	f.jobName = jobName
	f.err = err
	return f.retErr
}

func TestMultiNotifier_CallsAllNotifiers(t *testing.T) {
	a := &fakeNotifier{}
	b := &fakeNotifier{}
	m := NewMultiNotifier(a, b)

	if err := m.Notify("myjob", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.called != 1 || b.called != 1 {
		t.Errorf("expected each notifier called once, got a=%d b=%d", a.called, b.called)
	}
	if a.jobName != "myjob" || b.jobName != "myjob" {
		t.Error("jobName not propagated correctly")
	}
}

func TestMultiNotifier_ContinuesOnError(t *testing.T) {
	a := &fakeNotifier{retErr: errors.New("smtp down")}
	b := &fakeNotifier{}
	m := NewMultiNotifier(a, b)

	err := m.Notify("job", errors.New("exit 1"))
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
	// b must still have been called despite a failing.
	if b.called != 1 {
		t.Errorf("second notifier should be called even after first fails, got called=%d", b.called)
	}
}

func TestMultiNotifier_CombinesErrors(t *testing.T) {
	a := &fakeNotifier{retErr: errors.New("err-a")}
	b := &fakeNotifier{retErr: errors.New("err-b")}
	m := NewMultiNotifier(a, b)

	err := m.Notify("job", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "2 error(s)") {
		t.Errorf("error message should mention count, got: %v", err)
	}
}

func TestMultiNotifier_Add_IncreasesLen(t *testing.T) {
	m := NewMultiNotifier()
	if m.Len() != 0 {
		t.Fatalf("expected 0 notifiers, got %d", m.Len())
	}
	m.Add(&fakeNotifier{})
	m.Add(&fakeNotifier{})
	if m.Len() != 2 {
		t.Errorf("expected 2 notifiers, got %d", m.Len())
	}
}

func TestMultiNotifier_Empty_NoError(t *testing.T) {
	m := NewMultiNotifier()
	if err := m.Notify("job", nil); err != nil {
		t.Errorf("empty MultiNotifier should return nil, got %v", err)
	}
}
