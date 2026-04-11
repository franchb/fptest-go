package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplyAssociativeCompositionEngine(t, cfg.runner, eqFC,
		ptdAB, ptdBC, fmapCompose, apAB, apBC, apAC, apABAC,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genAB), enginerapid.Wrap(genBC))
}
