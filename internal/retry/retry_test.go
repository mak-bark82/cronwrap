package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/cronwrap/internal/retry"
)

var errFail = errors.New("fail")

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	attempts := retry.Do(context.Background(), p, func(_ context.Context) error {
		return nil
	})

	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if attempts[0].Err != nil {
		t.Fatalf("expected no error, got %v", attempts[0].Err)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	attempts := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errFail
		}
		return nil
	})

	if len(attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(attempts))
	}
	if !retry.Succeeded(attempts) {
		t.Fatal("expected success on final attempt")
	}
}

func TestDo_AllAttemptsFail(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	attempts := retry.Do(context.Background(), p, func(_ context.Context) error {
		return errFail
	})

	if len(attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(attempts))
	}
	if retry.Succeeded(attempts) {
		t.Fatal("expected failure")
	}
}

func TestDo_ContextCancelledBetweenRetries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	p := retry.Policy{MaxAttempts: 5, Delay: 10 * time.Millisecond}
	attempts := retry.Do(ctx, p, func(_ context.Context) error {
		calls++
		if calls == 2 {
			cancel()
		}
		return errFail
	})

	if retry.Succeeded(attempts) {
		t.Fatal("expected failure after cancel")
	}
	if len(attempts) > 4 {
		t.Fatalf("too many attempts after cancel: %d", len(attempts))
	}
}

func TestDo_AttemptNumbersAreSequential(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	attempts := retry.Do(context.Background(), p, func(_ context.Context) error {
		return errFail
	})

	for i, a := range attempts {
		if a.Number != i+1 {
			t.Errorf("attempt[%d].Number = %d, want %d", i, a.Number, i+1)
		}
	}
}

func TestSucceeded_EmptySlice(t *testing.T) {
	if retry.Succeeded(nil) {
		t.Fatal("expected false for empty slice")
	}
}
