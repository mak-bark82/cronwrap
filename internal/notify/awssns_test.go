package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/example/cronwrap/internal/notify"
)

// fakeSNSPublisher records the last Publish call.
type fakeSNSPublisher struct {
	input     *sns.PublishInput
	returnErr error
}

func (f *fakeSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	f.input = params
	return &sns.PublishOutput{}, f.returnErr
}

func TestNewSNSNotifier_SetsTopicARN(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:my-topic", fake)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSNSNotifier_Notify_Success(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", fake)

	if err := n.Notify("backup-job", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.input == nil {
		t.Fatal("expected Publish to be called")
	}
	if *fake.input.Subject != "cronwrap: backup-job" {
		t.Errorf("unexpected subject: %s", *fake.input.Subject)
	}
}

func TestSNSNotifier_Notify_WithError(t *testing.T) {
	fake := &fakeSNSPublisher{}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", fake)

	jobErr := errors.New("exit status 1")
	if err := n.Notify("backup-job", jobErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.input == nil {
		t.Fatal("expected Publish to be called")
	}
	msg := *fake.input.Message
	if msg == "" {
		t.Error("expected non-empty message body")
	}
}

func TestSNSNotifier_Notify_PublishError(t *testing.T) {
	fake := &fakeSNSPublisher{returnErr: errors.New("aws error")}
	n := notify.NewSNSNotifierWithClient("arn:aws:sns:us-east-1:123456789012:alerts", fake)

	err := n.Notify("backup-job", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
