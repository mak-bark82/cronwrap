package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// snsPublisher abstracts the SNS Publish call for testing.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends job alerts to an AWS SNS topic.
type SNSNotifier struct {
	topicARN string
	client   snsPublisher
}

// NewSNSNotifier creates an SNSNotifier using the default AWS credential chain.
// topicARN must be a valid SNS topic ARN.
func NewSNSNotifier(topicARN string) (*SNSNotifier, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: load aws config: %w", err)
	}
	return &SNSNotifier{
		topicARN: topicARN,
		client:   sns.NewFromConfig(cfg),
	}, nil
}

// newSNSNotifierWithClient creates an SNSNotifier with a custom publisher (used in tests).
func newSNSNotifierWithClient(topicARN string, client snsPublisher) *SNSNotifier {
	return &SNSNotifier{topicARN: topicARN, client: client}
}

// Notify publishes a message to the configured SNS topic.
func (n *SNSNotifier) Notify(jobName string, err error) error {
	subject, body := snsMessage(jobName, err)
	_, pubErr := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Subject:  aws.String(subject),
		Message:  aws.String(body),
	})
	if pubErr != nil {
		return fmt.Errorf("sns: publish: %w", pubErr)
	}
	return nil
}

func snsMessage(jobName string, err error) (subject, body string) {
	if err != nil {
		return fmt.Sprintf("[cronwrap] FAILED: %s", jobName),
			fmt.Sprintf("Job %q failed with error: %v", jobName, err)
	}
	return fmt.Sprintf("[cronwrap] OK: %s", jobName),
		fmt.Sprintf("Job %q completed successfully.", jobName)
}
