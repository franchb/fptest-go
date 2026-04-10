package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonadLawsEngine[FA, FB, FC, A, B, C](t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA),
		enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAB, chainBC, chainAC)
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonadLawsFullEngine[FA, FB, FC, A, B, C](t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA),
		enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAA, chainAB, chainBC, chainAC)
}
