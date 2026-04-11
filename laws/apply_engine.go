package laws

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
)

// ApplyAssociativeCompositionEngine verifies the Apply associative composition law
// using the engine abstraction.
func ApplyAssociativeCompositionEngine[FA, FB, FC, FAB, FBC, FAC, FABAC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	eqFC func(FC, FC) bool,
	ptdAB Pointed[func(A) B, FAB],
	ptdBC Pointed[func(B) C, FBC],
	fmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC],
	apAB Apply[A, B, FA, FB, FAB],
	apBC Apply[B, C, FB, FC, FBC],
	apAC Apply[A, C, FA, FC, FAC],
	apABAC Apply[func(A) B, func(A) C, FAB, FAC, FABAC],
	genFA engine.Generator[FA],
	genAB engine.Generator[func(A) B],
	genBC engine.Generator[func(B) C],
) {
	t.Helper()

	runner.MakeCheck(t, "Apply/AssociativeComposition", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		ab := genAB.Draw(et, "ab")
		bc := genBC.Draw(et, "bc")

		fab := ptdAB.Of(ab)
		fbc := ptdBC.Of(bc)

		// compose: (B -> C) -> (A -> B) -> (A -> C)
		compose := func(g func(B) C) func(func(A) B) func(A) C {
			return func(f func(A) B) func(A) C {
				return func(a A) C { return g(f(a)) }
			}
		}

		// Left: Ap(Ap(Map(compose)(fbc))(fab))(fa)
		composed := fmapCompose.Map(compose)(fbc) // F[(A->B) -> (A->C)]
		applied := apABAC.Ap(fab)(composed)        // F[A -> C]
		left := apAC.Ap(fa)(applied)               // F[C]

		// Right: Ap(fbc)(Ap(fab)(fa))
		inner := apAB.Ap(fa)(fab)    // F[B]
		right := apBC.Ap(inner)(fbc) // F[C]

		if !eqFC(left, right) {
			et.Fatalf("Apply associative composition violated:\n  Ap(Ap(Map(compose)(fbc))(fab))(fa) = %v\n  Ap(fbc)(Ap(fab)(fa))               = %v", left, right)
		}
	})
}
