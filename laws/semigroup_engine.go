package laws

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
)

// SemigroupLawsEngine verifies the Semigroup law (associativity) using the engine abstraction.
func SemigroupLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
) {
	t.Helper()

	runner.MakeCheck(t, "Semigroup/Associativity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")

		// Law: concat(concat(a, b), c) == concat(a, concat(b, c))
		left := concat(concat(a, b), c)
		right := concat(a, concat(b, c))
		if !eqA(left, right) {
			et.Fatalf("Semigroup associativity violated:\n  (a <> b) <> c = %v\n  a <> (b <> c) = %v", left, right)
		}
	})
}
