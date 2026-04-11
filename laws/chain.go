package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ChainAssociativityEngine[FA, FB, FC, A, B, C](t, cfg.runner, eqFC,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		chainAB, chainBC, chainAC)
}
