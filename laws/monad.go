package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// MonadLaws verifies the Monad laws (left identity, right identity, associativity).
//
// Type parameters:
//   - FA, FB, FC: container types F[A], F[B], F[C]
//   - A, B, C: element types
//
// Parameters:
//   - of: pure/return (lifts a value into the monad)
//   - chainAB, chainBC, chainAC: monadic bind (>>=) specialized per type pair
//   - genA: generator for values of type A
//   - genKleisliAB: generator for Kleisli arrows A -> F[B]
//   - genKleisliBC: generator for Kleisli arrows B -> F[C]
func MonadLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	t.Run("Monad/LeftIdentity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genKleisliAB.Draw(t, "f")

		// Law: chain(f)(of(a)) == f(a)
		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			t.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	}))

	t.Run("Monad/Associativity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		f := genKleisliAB.Draw(t, "f")
		g := genKleisliBC.Draw(t, "g")

		// Law: chain(g)(chain(f)(fa)) == chain(x => chain(g)(f(x)))(fa)
		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			t.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	}))
}

// MonadLawsFull verifies all Monad laws including right identity, requiring an additional
// chain operation for the same-type case (A -> FA -> FA).
func MonadLawsFull[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
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

	t.Run("Monad/LeftIdentity", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genKleisliAB.Draw(t, "f")

		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			t.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	}))

	t.Run("Monad/RightIdentity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")

		// Law: chain(of)(fa) == fa
		got := chainAA(of)(fa)
		if !eqFA(got, fa) {
			t.Fatalf("Monad right identity violated:\n  chain(of)(fa) = %v\n  fa            = %v", got, fa)
		}
	}))

	t.Run("Monad/Associativity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		f := genKleisliAB.Draw(t, "f")
		g := genKleisliBC.Draw(t, "g")

		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			t.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	}))
}
