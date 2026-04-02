package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// MonoidLaws verifies the Monoid laws (left identity, right identity, associativity).
// Monoid extends Semigroup with an identity element.
func MonoidLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
) {
	t.Helper()

	// Semigroup associativity
	SemigroupLaws(t, genA, eqA, concat)

	t.Run("Monoid/LeftIdentity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		// Law: concat(empty, a) == a
		got := concat(empty, a)
		if !eqA(got, a) {
			t.Fatalf("Monoid left identity violated:\n  empty <> a = %v\n  a          = %v", got, a)
		}
	}))

	t.Run("Monoid/RightIdentity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		// Law: concat(a, empty) == a
		got := concat(a, empty)
		if !eqA(got, a) {
			t.Fatalf("Monoid right identity violated:\n  a <> empty = %v\n  a          = %v", got, a)
		}
	}))
}

// MonoidInterfaceLaws verifies Monoid laws using fp-go's Monoid interface.
func MonoidInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	m interface {
		Concat(A, A) A
		Empty() A
	},
) {
	t.Helper()
	MonoidLaws(t, genA, eqA, m.Concat, m.Empty())
}
