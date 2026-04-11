package laws

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
)

// EqLawsEngine verifies the Eq laws (reflexivity, symmetry, transitivity) using the engine abstraction.
func EqLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	equals func(A, A) bool,
) {
	t.Helper()

	runner.MakeCheck(t, "Eq/Reflexivity", func(et engine.T) {
		a := genA.Draw(et, "a")

		// Law: equals(a, a) == true
		if !equals(a, a) {
			et.Fatalf("Eq reflexivity violated: equals(%v, %v) = false", a, a)
		}
	})

	runner.MakeCheck(t, "Eq/Symmetry", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		// Law: equals(a, b) == equals(b, a)
		ab := equals(a, b)
		ba := equals(b, a)
		if ab != ba {
			et.Fatalf("Eq symmetry violated:\n  equals(%v, %v) = %v\n  equals(%v, %v) = %v", a, b, ab, b, a, ba)
		}
	})

	runner.MakeCheck(t, "Eq/Transitivity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")

		// Law: if equals(a, b) && equals(b, c) then equals(a, c)
		if equals(a, b) && equals(b, c) && !equals(a, c) {
			et.Fatalf("Eq transitivity violated:\n  equals(%v, %v) = true\n  equals(%v, %v) = true\n  equals(%v, %v) = false", a, b, b, c, a, c)
		}
	})
}
