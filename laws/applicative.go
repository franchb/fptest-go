package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// ApplicativeLaws verifies the Applicative functor laws.
//
// Due to Go's type system limitations (no HKTs), this tests the most important
// Applicative properties: identity via fmap and homomorphism via ap.
//
// Type parameters:
//   - FA, FB, FAB: container types F[A], F[B], F[func(A) B]
//   - A, B: element types
func ApplicativeLaws[FA, FB, FAB, A, B any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	eqFA func(FA, FA) bool,
	eqFB func(FB, FB) bool,
	ofA func(A) FA,
	ofB func(B) FB,
	ofAB func(func(A) B) FAB,
	fmapAA func(func(A) A) func(FA) FA,
	apAB func(FA) func(FAB) FB,
	identity func(A) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA), enginerapid.Wrap(genAB),
		eqFA, eqFB, ofA, ofB, ofAB, fmapAA, apAB, identity)
}

// ApplicativeInterchange verifies the Applicative interchange law:
//
//	Ap(Of(a))(u) == Ap(u)(Of(f => f(a)))
//
// Where u : F[func(A) B]. This says: applying a wrapped function to a pure value
// is the same as applying "call-with-that-value" to the wrapped function.
//
// Type parameters:
//   - FA, FB, FAB: container types F[A], F[B], F[func(A) B]
//   - FABB: container type F[func(func(A) B) B]
//   - A, B: element types
func ApplicativeInterchange[FA, FB, FAB, FABB, A, B any](
	t *testing.T,
	eqFB func(FB, FB) bool,
	apAB Applicative[A, B, FA, FB, FAB],
	ptdAB Pointed[func(A) B, FAB],
	apABB Apply[func(A) B, B, FAB, FB, FABB],
	ptdABB Pointed[func(func(A) B) B, FABB],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeInterchangeEngine(t, cfg.runner, eqFB,
		apAB, ptdAB, apABB, ptdABB,
		enginerapid.Wrap(genA), enginerapid.Wrap(genAB))
}

// ApplicativeFullLaws verifies all four Applicative functor laws using a
// pre-built [ApplicativeInstances] bundle:
//
//  1. Identity:     fmap(id)(v) == v
//  2. Homomorphism: ap(of(f))(of(a)) == of(f(a))
//  3. Interchange:  ap(of(a))(u) == ap(u)(of(f => f(a)))
//  4. Composition:  ap(ap(map(compose)(fbc))(fab))(fa) == ap(fbc)(ap(fab)(fa))
//
// The instances bundle avoids forcing callers to construct 12+ interface values
// manually.
func ApplicativeFullLaws[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	inst *ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C],
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeFullLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genA),
		enginerapid.Wrap(genAB), enginerapid.Wrap(genBC),
		inst)
}
