// Package retry provides a simple, context-aware retry mechanism for cronwrap.
//
// Use [Do] to execute a [Func] up to Policy.MaxAttempts times, sleeping
// Policy.Delay (plus optional jitter) between failures.  The returned
// []Attempt slice records every execution so callers can log or surface
// per-attempt diagnostics.
//
// Example:
//
//	p := retry.Policy{MaxAttempts: 3, Delay: 5 * time.Second}
//	attempts := retry.Do(ctx, p, func(ctx context.Context) error {
//		return runMyJob(ctx)
//	})
//	if !retry.Succeeded(attempts) {
//		log.Printf("job failed after %d attempts", len(attempts))
//	}
package retry
