package notify

// NewDatadogNotifierWithURL creates a DatadogNotifier with a custom API URL.
// This is exported for testing purposes only.
func NewDatadogNotifierWithURL(apiKey, apiURL string) *DatadogNotifier {
	n := NewDatadogNotifier(apiKey)
	n.apiURL = apiURL
	return n
}
