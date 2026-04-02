package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// ApplyAssociativeComposition verifies the Apply associative composition law:
//
//	Ap(Ap(Map(compose)(fbc))(fab))(fa) == Ap(fbc)(Ap(fab)(fa))
//
// This law ensures that applying composed functions equals composing applications.
//
// Type parameters:
//   - FA, FB, FC: container types F[A], F[B], F[C]
//   - FAB, FBC, FAC: container types for lifted functions F[A->B], F[B->C], F[A->C]
//   - FABAC: container type F[(A->B) -> (A->C)]
//   - A, B, C: element types
func ApplyAssociativeComposition[FA, FB, FC, FAB, FBC, FAC, FABAC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	ptdAB Pointed[func(A) B, FAB],
	ptdBC Pointed[func(B) C, FBC],
	fmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC],
	apAB Apply[A, B, FA, FB, FAB],
	apBC Apply[B, C, FB, FC, FBC],
	apAC Apply[A, C, FA, FC, FAC],
	apABAC Apply[func(A) B, func(A) C, FAB, FAC, FABAC],
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
) {
	t.Helper()

	t.Run("Apply/AssociativeComposition", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		ab := genAB.Draw(t, "ab")
		bc := genBC.Draw(t, "bc")

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
			t.Fatalf("Apply associative composition violated:\n  Ap(Ap(Map(compose)(fbc))(fab))(fa) = %v\n  Ap(fbc)(Ap(fab)(fa))               = %v", left, right)
		}
	}))
}
