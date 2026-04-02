package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// SemigroupLaws verifies the Semigroup law (associativity).
func SemigroupLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
) {
	t.Helper()

	t.Run("Semigroup/Associativity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")
		c := genA.Draw(t, "c")

		// Law: concat(concat(a, b), c) == concat(a, concat(b, c))
		left := concat(concat(a, b), c)
		right := concat(a, concat(b, c))
		if !eqA(left, right) {
			t.Fatalf("Semigroup associativity violated:\n  (a <> b) <> c = %v\n  a <> (b <> c) = %v", left, right)
		}
	}))
}

// SemigroupInterfaceLaws verifies Semigroup laws using fp-go's Semigroup interface.
func SemigroupInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	sg interface{ Concat(A, A) A },
) {
	t.Helper()
	SemigroupLaws(t, genA, eqA, sg.Concat)
}
