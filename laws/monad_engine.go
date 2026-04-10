package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// MonadLawsEngine verifies the Monad laws (left identity, associativity) using the engine abstraction.
func MonadLawsEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Monad/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genKleisliAB.Draw(et, "f")

		// Law: chain(f)(of(a)) == f(a)
		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			et.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	})

	runner.MakeCheck(t, "Monad/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")

		// Law: chain(g)(chain(f)(fa)) == chain(x => chain(g)(f(x)))(fa)
		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			et.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	})
}

// MonadLawsFullEngine verifies all Monad laws including right identity using the engine abstraction.
func MonadLawsFullEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAA func(func(A) FA) func(FA) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Monad/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genKleisliAB.Draw(et, "f")

		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			et.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	})

	runner.MakeCheck(t, "Monad/RightIdentity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")

		// Law: chain(of)(fa) == fa
		got := chainAA(of)(fa)
		if !eqFA(got, fa) {
			et.Fatalf("Monad right identity violated:\n  chain(of)(fa) = %v\n  fa            = %v", got, fa)
		}
	})

	runner.MakeCheck(t, "Monad/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")

		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			et.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	})
}
