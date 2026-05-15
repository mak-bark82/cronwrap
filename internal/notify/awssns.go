package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// snsPublisher is the subset of the SNS client we use.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends job notifications via AWS SNS.
type SNSNotifier struct {
	topicARN string
	client   snsPublisher
}

// NewSNSNotifier creates an SNSNotifier using the default AWS SDK configuration.
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

// newSNSNotifierWithClient creates an SNSNotifier with a custom publisher (for testing).
func newSNSNotifierWithClient(topicARN string, client snsPublisher) *SNSNotifier {
	return &SNSNotifier{topicARN: topicARN, client: client}
}

// Notify sends a message to the configured SNS topic.
func (n *SNSNotifier) Notify(jobName string, err error) error {
	body := snsMessage(jobName, err)
	_, pubErr := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Message:  aws.String(body),
		Subject:  aws.String(fmt.Sprintf("cronwrap: %s", jobName)),
	})
	if pubErr != nil {
		return fmt.Errorf("sns: publish: %w", pubErr)
	}
	return nil
}

func snsMessage(jobName string, err error) string {
	if err != nil {
		return fmt.Sprintf("Job %q failed: %v", jobName, err)
	}
	return fmt.Sprintf("Job %q completed successfully.", jobName)
}
