package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// OrdLawsEngine verifies the Ord laws (antisymmetry, transitivity, totality, consistency)
// using the engine abstraction. The compare function must return -1, 0, or 1.
func OrdLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
) {
	t.Helper()

	// Include Eq laws since Ord extends Eq
	EqLawsEngine(t, runner, genA, equals)

	runner.MakeCheck(t, "Ord/Antisymmetry", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		// Law: if compare(a, b) <= 0 && compare(b, a) <= 0 then equals(a, b)
		if compare(a, b) <= 0 && compare(b, a) <= 0 && !equals(a, b) {
			et.Fatalf("Ord antisymmetry violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  but not equal", a, b, compare(a, b), b, a, compare(b, a))
		}
	})

	runner.MakeCheck(t, "Ord/Transitivity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")

		// Law: if compare(a, b) <= 0 && compare(b, c) <= 0 then compare(a, c) <= 0
		if compare(a, b) <= 0 && compare(b, c) <= 0 && compare(a, c) > 0 {
			et.Fatalf("Ord transitivity violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, c, compare(b, c), a, c, compare(a, c))
		}
	})

	runner.MakeCheck(t, "Ord/Totality", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		// Law: compare(a, b) <= 0 || compare(b, a) <= 0
		if compare(a, b) > 0 && compare(b, a) > 0 {
			et.Fatalf("Ord totality violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, a, compare(b, a))
		}
	})

	runner.MakeCheck(t, "Ord/Consistency", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		// Law: compare(a, b) == 0 iff equals(a, b)
		cmp := compare(a, b)
		eq := equals(a, b)
		if (cmp == 0) != eq {
			et.Fatalf("Ord consistency violated:\n  compare(%v, %v) = %d\n  equals(%v, %v) = %v",
				a, b, cmp, a, b, eq)
		}
	})
}
