package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"

	"cronwrap/internal/notify"
)

// fakeSNSPublisher records the last Publish call.
type fakeSNSPublisher struct {
	lastInput *sns.PublishInput
	returnErr error
}

func (f *fakeSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	f.lastInput = params
	return &sns.PublishOutput{}, f.returnErr
}

func TestNewSNSNotifier_SetsTopicARN(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123:my-topic", fake)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSNSNotifier_Notify_Success(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123:my-topic", fake)

	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastInput == nil {
		t.Fatal("expected Publish to be called")
	}
	subject := *fake.lastInput.Subject
	if subject != "[cronwrap] OK: backup-job" {
		t.Errorf("unexpected subject: %q", subject)
	}
}

func TestSNSNotifier_Notify_WithError(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123:my-topic", fake)

	jobErr := errors.New("exit status 1")
	if err := n.Notify("backup-job", jobErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	subject := *fake.lastInput.Subject
	if subject != "[cronwrap] FAILED: backup-job" {
		t.Errorf("unexpected subject: %q", subject)
	}
}

func TestSNSNotifier_Notify_PublishError(t *testing.T) {
	fake := &fakeSNSPublisher{returnErr: errors.New("aws error")}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123:my-topic", fake)

	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, fake.returnErr) {
		t.Errorf("expected wrapped aws error, got: %v", err)
	}
}

func TestSNSNotifier_Notify_MessageContainsJobName(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123:my-topic", fake)

	_ = n.Notify("my-special-job", nil)
	body := *fake.lastInput.Message
	if body == "" {
		t.Fatal("expected non-empty message body")
	}
	if !contains(body, "my-special-job") {
		t.Errorf("expected body to contain job name, got: %q", body)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
