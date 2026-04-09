package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/result"
)

// AssertOk extracts the value from a successful Result, or fails the test if it is an error.
// Returns the unwrapped value for chained assertions.
func AssertOk[A any](t testing.TB, r result.Result[A]) A {
	t.Helper()
	return AssertRight[error, A](t, r)
}

// AssertErr extracts the error from a failed Result, or fails the test if it is a success.
// Returns the unwrapped error for chained assertions.
func AssertErr[A any](t testing.TB, r result.Result[A]) error {
	t.Helper()
	return AssertLeft[error, A](t, r)
}

// AssertOkEq extracts the value from a successful Result and asserts it equals the expected value.
func AssertOkEq[A comparable](t testing.TB, r result.Result[A], want A) {
	t.Helper()
	got := AssertOk[A](t, r)
	if got != want {
		t.Fatalf("Result value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}
