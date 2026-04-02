package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// OrdLaws verifies the Ord laws (antisymmetry, transitivity, totality).
// The compare function must return -1, 0, or 1.
func OrdLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
) {
	t.Helper()

	// Include Eq laws since Ord extends Eq
	EqLaws(t, genA, equals)

	t.Run("Ord/Antisymmetry", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		// Law: if compare(a, b) <= 0 && compare(b, a) <= 0 then equals(a, b)
		if compare(a, b) <= 0 && compare(b, a) <= 0 && !equals(a, b) {
			t.Fatalf("Ord antisymmetry violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  but not equal", a, b, compare(a, b), b, a, compare(b, a))
		}
	}))

	t.Run("Ord/Transitivity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")
		c := genA.Draw(t, "c")

		// Law: if compare(a, b) <= 0 && compare(b, c) <= 0 then compare(a, c) <= 0
		if compare(a, b) <= 0 && compare(b, c) <= 0 && compare(a, c) > 0 {
			t.Fatalf("Ord transitivity violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, c, compare(b, c), a, c, compare(a, c))
		}
	}))

	t.Run("Ord/Totality", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		// Law: compare(a, b) <= 0 || compare(b, a) <= 0
		if compare(a, b) > 0 && compare(b, a) > 0 {
			t.Fatalf("Ord totality violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, a, compare(b, a))
		}
	}))

	t.Run("Ord/Consistency", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		// Law: compare(a, b) == 0 iff equals(a, b)
		cmp := compare(a, b)
		eq := equals(a, b)
		if (cmp == 0) != eq {
			t.Fatalf("Ord consistency violated:\n  compare(%v, %v) = %d\n  equals(%v, %v) = %v",
				a, b, cmp, a, b, eq)
		}
	}))
}

// OrdInterfaceLaws verifies Ord laws using fp-go's Ord interface.
func OrdInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	ordA interface {
		Equals(A, A) bool
		Compare(A, A) int
	},
) {
	t.Helper()
	OrdLaws(t, genA, ordA.Equals, ordA.Compare)
}
