package prop

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// RoundTripEngine verifies that encode and decode are inverses: decode(encode(a)) == a.
// This is the engine-generic variant of RoundTrip.
func RoundTripEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTrip", func(et engine.T) {
		a := genA.Draw(et, "a")

		got := decode(encode(a))
		if !eqA(got, a) {
			et.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}

// RoundTripPartialEngine verifies a round-trip where decode may fail, returning (A, bool).
// This is the engine-generic variant of RoundTripPartial.
func RoundTripPartialEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTripPartial", func(et engine.T) {
		a := genA.Draw(et, "a")

		got, ok := decode(encode(a))
		if !ok {
			et.Fatalf("Round-trip decode failed for input: %v", a)
		}
		if !eqA(got, a) {
			et.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}

// RoundTripErrorEngine verifies a round-trip where decode may return an error.
// This is the engine-generic variant of RoundTripError.
func RoundTripErrorEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTripError", func(et engine.T) {
		a := genA.Draw(et, "a")

		got, err := decode(encode(a))
		if err != nil {
			et.Fatalf("Round-trip decode error for input %v: %v", a, err)
		}
		if !eqA(got, a) {
			et.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}
