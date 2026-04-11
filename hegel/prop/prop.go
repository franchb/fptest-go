// Package prop provides hegel-native convenience wrappers for property testing.
// Each function creates a HegelRunner and wraps hegel generators, then delegates
// to the corresponding engine-generic function in the core prop package.
package prop

import (
	"testing"

	fpthegel "github.com/franchb/fptest-go/hegel"
	coreprop "github.com/franchb/fptest-go/prop"
	hegellib "hegel.dev/go/hegel"
)

var hegelRunner = fpthegel.HegelRunner{}

// RoundTrip verifies that encode and decode are inverses: decode(encode(a)) == a.
func RoundTrip[A, B any](t *testing.T, name string, genA hegellib.Generator[A], eqA func(A, A) bool, encode func(A) B, decode func(B) A) {
	t.Helper()
	coreprop.RoundTripEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// RoundTripPartial verifies a round-trip where decode may fail, returning (A, bool).
func RoundTripPartial[A, B any](t *testing.T, name string, genA hegellib.Generator[A], eqA func(A, A) bool, encode func(A) B, decode func(B) (A, bool)) {
	t.Helper()
	coreprop.RoundTripPartialEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// RoundTripError verifies a round-trip where decode may return an error.
func RoundTripError[A, B any](t *testing.T, name string, genA hegellib.Generator[A], eqA func(A, A) bool, encode func(A) B, decode func(B) (A, error)) {
	t.Helper()
	coreprop.RoundTripErrorEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// Oracle verifies that an implementation produces the same results as a reference
// implementation for all generated inputs.
func Oracle[A, B any](t *testing.T, name string, genA hegellib.Generator[A], eqB func(B, B) bool, impl func(A) B, reference func(A) B) {
	t.Helper()
	coreprop.OracleEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqB, impl, reference)
}

// Idempotent verifies that applying a function twice yields the same result as applying once.
func Idempotent[A any](t *testing.T, name string, genA hegellib.Generator[A], eqA func(A, A) bool, f func(A) A) {
	t.Helper()
	coreprop.IdempotentEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, f)
}

// Commutative verifies that f(a, b) == f(b, a) for all a, b.
func Commutative[A, B any](t *testing.T, name string, genA hegellib.Generator[A], eqB func(B, B) bool, f func(A, A) B) {
	t.Helper()
	coreprop.CommutativeEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqB, f)
}

// Invariant verifies that a predicate holds for all generated inputs.
func Invariant[A any](t *testing.T, name string, genA hegellib.Generator[A], predicate func(A) bool) {
	t.Helper()
	coreprop.InvariantEngine(t, hegelRunner, name, fpthegel.Wrap(genA), predicate)
}
