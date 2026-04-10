package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// ChainAssociativityEngine verifies the Chain associativity law using the engine abstraction.
func ChainAssociativityEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	eqFC func(FC, FC) bool,
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Chain/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")

		// Left: Chain(g)(Chain(f)(fa))
		left := chainBC(g)(chainAB(f)(fa))

		// Right: Chain(x => Chain(g)(f(x)))(fa)
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)

		if !eqFC(left, right) {
			et.Fatalf("Chain associativity violated:\n  chain(g)(chain(f)(fa))        = %v\n  chain(x=>chain(g)(f(x)))(fa) = %v", left, right)
		}
	})
}
