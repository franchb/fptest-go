// Package prop provides higher-level property testing utilities built on rapid.
package prop

import (
	"testing"

	"pgregory.net/rapid"
)

// RoundTrip verifies that encode and decode are inverses: decode(encode(a)) == a.
// This is the fundamental property for serialization, parsing, and any codec pair.
func RoundTrip[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
) {
	t.Helper()
	t.Run(name+"/RoundTrip", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		got := decode(encode(a))
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	}))
}

// RoundTripPartial verifies a round-trip where decode may fail, returning (A, bool).
// Only checks the property when decoding succeeds.
func RoundTripPartial[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
) {
	t.Helper()
	t.Run(name+"/RoundTripPartial", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		got, ok := decode(encode(a))
		if !ok {
			t.Fatalf("Round-trip decode failed for input: %v", a)
		}
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	}))
}

// RoundTripError verifies a round-trip where decode may return an error.
func RoundTripError[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
) {
	t.Helper()
	t.Run(name+"/RoundTripError", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		got, err := decode(encode(a))
		if err != nil {
			t.Fatalf("Round-trip decode error for input %v: %v", a, err)
		}
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	}))
}
