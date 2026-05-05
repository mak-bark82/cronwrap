package retry

import (
	"math/rand"
	"time"
)

// jitter returns a random duration in [0, max).
func jitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	//nolint:gosec // non-cryptographic jitter is intentional
	return time.Duration(rand.Int63n(int64(max)))
}
