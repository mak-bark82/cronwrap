package notify

// NewSNSNotifierWithClient exposes newSNSNotifierWithClient for use in
// external test packages that need to inject a fake SNS publisher.
func NewSNSNotifierWithClient(topicARN string, client snsPublisher) *SNSNotifier {
	return newSNSNotifierWithClient(topicARN, client)
}

// SNSPublisher re-exports the snsPublisher interface so test helpers outside
// this package can implement it without importing internal symbols directly.
type SNSPublisher = snsPublisher
