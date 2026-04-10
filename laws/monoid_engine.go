package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// MonoidLawsEngine verifies the Monoid laws (left identity, right identity, associativity)
// using the engine abstraction. Monoid extends Semigroup with an identity element.
func MonoidLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
) {
	t.Helper()

	// Semigroup associativity
	SemigroupLawsEngine(t, runner, genA, eqA, concat)

	runner.MakeCheck(t, "Monoid/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")

		// Law: concat(empty, a) == a
		got := concat(empty, a)
		if !eqA(got, a) {
			et.Fatalf("Monoid left identity violated:\n  empty <> a = %v\n  a          = %v", got, a)
		}
	})

	runner.MakeCheck(t, "Monoid/RightIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")

		// Law: concat(a, empty) == a
		got := concat(a, empty)
		if !eqA(got, a) {
			et.Fatalf("Monoid right identity violated:\n  a <> empty = %v\n  a          = %v", got, a)
		}
	})
}
