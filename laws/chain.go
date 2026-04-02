package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// ChainAssociativity verifies the Chain associativity law:
//
//	Chain(g)(Chain(f)(fa)) == Chain(x => Chain(g)(f(x)))(fa)
//
// This ensures monadic composition is independent of grouping.
//
// Type parameters:
//   - FA, FB, FC: container types F[A], F[B], F[C]
//   - A, B, C: element types
func ChainAssociativity[FA, FB, FC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	t.Run("Chain/Associativity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		f := genKleisliAB.Draw(t, "f")
		g := genKleisliBC.Draw(t, "g")

		// Left: Chain(g)(Chain(f)(fa))
		left := chainBC(g)(chainAB(f)(fa))

		// Right: Chain(x => Chain(g)(f(x)))(fa)
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)

		if !eqFC(left, right) {
			t.Fatalf("Chain associativity violated:\n  chain(g)(chain(f)(fa))        = %v\n  chain(x=>chain(g)(f(x)))(fa) = %v", left, right)
		}
	}))
}
