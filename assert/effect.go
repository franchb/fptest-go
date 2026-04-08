package assert

import (
	"context"
	"testing"

	"github.com/IBM/fp-go/v2/effect"
)

// AssertEffect runs an Effect with the given dependencies and context,
// then extracts the success value or fails the test.
func AssertEffect[C, A any](t testing.TB, deps C, ctx context.Context, eff effect.Effect[C, A]) A {
	t.Helper()
	val, err := effect.RunSync(effect.Provide[A](deps)(eff))(ctx)
	if err != nil {
		t.Fatalf("Expected Effect to succeed, got error: %v", err)
	}
	return val
}

// AssertEffectErr runs an Effect with the given dependencies and context,
// then extracts the error or fails the test if the Effect succeeded.
func AssertEffectErr[C, A any](t testing.TB, deps C, ctx context.Context, eff effect.Effect[C, A]) error {
	t.Helper()
	_, err := effect.RunSync(effect.Provide[A](deps)(eff))(ctx)
	if err == nil {
		t.Fatal("Expected Effect to fail, but it succeeded")
	}
	return err
}

// AssertEffectEq runs an Effect with the given dependencies and context,
// then asserts the success value equals the expected value.
func AssertEffectEq[C, A comparable](t testing.TB, deps C, ctx context.Context, eff effect.Effect[C, A], want A) {
	t.Helper()
	got := AssertEffect[C, A](t, deps, ctx, eff)
	if got != want {
		t.Fatalf("Effect value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}
