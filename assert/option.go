// Package assert provides FP-aware assertion helpers for testing.
//
// Assertions on sum types unwrap and return the inner value, enabling chained
// extraction: user := assert.AssertSome(t, assert.AssertRight(t, result)).
// All assertion functions call t.Helper() so errors point to the caller's line.
package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
)

// AssertSome extracts the value from a Some option, or fails the test if it is None.
// Returns the unwrapped value for chained assertions.
func AssertSome[A any](t testing.TB, o option.Option[A]) A {
	t.Helper()
	return option.Fold(
		func() A {
			t.Fatalf("Expected Some, got None")
			var zero A
			return zero
		},
		func(a A) A { return a },
	)(o)
}

// AssertNone asserts that the option is None.
func AssertNone[A any](t testing.TB, o option.Option[A]) {
	t.Helper()
	option.Fold(
		func() any { return nil },
		func(a A) any {
			t.Fatalf("Expected None, got Some(%v)", a)
			return nil
		},
	)(o)
}

// AssertSomeEq extracts the value from a Some option and asserts it equals the expected value.
func AssertSomeEq[A comparable](t testing.TB, o option.Option[A], want A) {
	t.Helper()
	got := AssertSome(t, o)
	if got != want {
		t.Fatalf("Option value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}
