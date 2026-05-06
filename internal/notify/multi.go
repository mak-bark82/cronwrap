package notify

import "fmt"

// Notifier is the interface implemented by every alert backend.
type Notifier interface {
	Notify(jobName string, err error) error
}

// MultiNotifier fans a single notification out to multiple Notifier
// implementations, collecting all errors.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that will call every provided
// Notifier in order when Notify is invoked.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify calls every registered notifier and returns a combined error if
// one or more of them fail. Remaining notifiers are still called even if
// an earlier one returns an error.
func (m *MultiNotifier) Notify(jobName string, err error) error {
	var errs []error
	for _, n := range m.notifiers {
		if nerr := n.Notify(jobName, err); nerr != nil {
			errs = append(errs, nerr)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("multi-notifier encountered %d error(s): %v", len(errs), errs)
}

// Add appends a Notifier to the set at runtime.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Len returns the number of registered notifiers.
func (m *MultiNotifier) Len() int {
	return len(m.notifiers)
}
