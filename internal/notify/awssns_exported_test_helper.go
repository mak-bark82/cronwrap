package notify

// NewSNSNotifierWithClient exposes the internal constructor for use in external test packages.
func NewSNSNotifierWithClient(topicARN string, client snsPublisher) *SNSNotifier {
	return newSNSNotifierWithClient(topicARN, client)
}
