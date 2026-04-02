package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// EqLaws verifies the Eq laws (reflexivity, symmetry, transitivity).
func EqLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
) {
	t.Helper()

	t.Run("Eq/Reflexivity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		// Law: equals(a, a) == true
		if !equals(a, a) {
			t.Fatalf("Eq reflexivity violated: equals(%v, %v) = false", a, a)
		}
	}))

	t.Run("Eq/Symmetry", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		// Law: equals(a, b) == equals(b, a)
		ab := equals(a, b)
		ba := equals(b, a)
		if ab != ba {
			t.Fatalf("Eq symmetry violated:\n  equals(%v, %v) = %v\n  equals(%v, %v) = %v", a, b, ab, b, a, ba)
		}
	}))

	t.Run("Eq/Transitivity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")
		c := genA.Draw(t, "c")

		// Law: if equals(a, b) && equals(b, c) then equals(a, c)
		if equals(a, b) && equals(b, c) && !equals(a, c) {
			t.Fatalf("Eq transitivity violated:\n  equals(%v, %v) = true\n  equals(%v, %v) = true\n  equals(%v, %v) = false", a, b, b, c, a, c)
		}
	}))
}

// EqInterfaceLaws verifies Eq laws using fp-go's Eq interface.
func EqInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA interface{ Equals(A, A) bool },
) {
	t.Helper()
	EqLaws(t, genA, eqA.Equals)
}
