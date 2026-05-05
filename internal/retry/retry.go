package retry

import (
	"context"
	"time"
)

// Policy defines the retry behaviour for a job.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Jitter      time.Duration
}

// Attempt holds the result of a single execution attempt.
type Attempt struct {
	Number   int
	Err      error
	Duration time.Duration
}

// Func is the function signature that retry executes.
type Func func(ctx context.Context) error

// Do runs fn according to p, returning all attempts made.
// It stops early when fn succeeds or the context is cancelled.
func Do(ctx context.Context, p Policy, fn Func) []Attempt {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}

	attempts := make([]Attempt, 0, p.MaxAttempts)

	for i := 1; i <= p.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			attempts = append(attempts, Attempt{Number: i, Err: err})
			break
		}

		start := time.Now()
		err := fn(ctx)
		dur := time.Since(start)

		attempts = append(attempts, Attempt{Number: i, Err: err, Duration: dur})

		if err == nil {
			break
		}

		if i < p.MaxAttempts {
			sleep := p.Delay
			if p.Jitter > 0 {
				sleep += jitter(p.Jitter)
			}
			select {
			case <-ctx.Done():
			case <-time.After(sleep):
			}
		}
	}

	return attempts
}

// Succeeded returns true when the last attempt in the slice was error-free.
func Succeeded(attempts []Attempt) bool {
	if len(attempts) == 0 {
		return false
	}
	return attempts[len(attempts)-1].Err == nil
}
