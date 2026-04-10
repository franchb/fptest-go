package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// FunctorLawsEngine verifies the Functor laws (identity and composition) using the engine abstraction.
func FunctorLawsEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genFA engine.Generator[FA],
	genAB engine.Generator[func(A) B],
	genBC engine.Generator[func(B) C],
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	fmapAA func(func(A) A) func(FA) FA,
	fmapAB func(func(A) B) func(FA) FB,
	fmapBC func(func(B) C) func(FB) FC,
	fmapAC func(func(A) C) func(FA) FC,
	identity func(A) A,
	compose func(func(A) B, func(B) C) func(A) C,
) {
	t.Helper()

	runner.MakeCheck(t, "Functor/Identity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")

		// Law: fmap(id)(fa) == fa
		got := fmapAA(identity)(fa)
		if !eqFA(got, fa) {
			et.Fatalf("Functor identity law violated:\n  fa       = %v\n  fmap(id) = %v", fa, got)
		}
	})

	runner.MakeCheck(t, "Functor/Composition", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genAB.Draw(et, "f")
		g := genBC.Draw(et, "g")

		// Law: fmap(g . f)(fa) == fmap(g)(fmap(f)(fa))
		composed := fmapAC(compose(f, g))(fa)
		chained := fmapBC(g)(fmapAB(f)(fa))
		if !eqFC(composed, chained) {
			et.Fatalf("Functor composition law violated:\n  fmap(g.f)(fa) = %v\n  fmap(g)(fmap(f)(fa)) = %v", composed, chained)
		}
	})
}
