package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
)

// AssertRight extracts the value from a Right either, or fails the test if it is Left.
// Returns the unwrapped value for chained assertions.
func AssertRight[E, A any](t testing.TB, e either.Either[E, A]) A {
	t.Helper()
	return either.Fold(
		func(err E) A {
			t.Fatalf("Expected Right, got Left(%v)", err)
			var zero A
			return zero
		},
		func(a A) A { return a },
	)(e)
}

// AssertLeft extracts the error from a Left either, or fails the test if it is Right.
// Returns the unwrapped error value for chained assertions.
func AssertLeft[E, A any](t testing.TB, e either.Either[E, A]) E {
	t.Helper()
	return either.Fold(
		func(err E) E { return err },
		func(a A) E {
			t.Fatalf("Expected Left, got Right(%v)", a)
			var zero E
			return zero
		},
	)(e)
}

// AssertRightEq extracts the value from a Right and asserts it equals the expected value.
func AssertRightEq[E any, A comparable](t testing.TB, e either.Either[E, A], want A) {
	t.Helper()
	got := AssertRight[E, A](t, e)
	if got != want {
		t.Fatalf("Right value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}

// AssertLeftEq extracts the error from a Left and asserts it equals the expected error.
func AssertLeftEq[E comparable, A any](t testing.TB, e either.Either[E, A], want E) {
	t.Helper()
	got := AssertLeft[E, A](t, e)
	if got != want {
		t.Fatalf("Left value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}
