// Package laws provides typeclass law verification using property-based testing.
//
// Each law test function runs as subtests via rapid.MakeCheck, producing output
// like TestOptionLaws/Functor/Identity that integrates with go test -run.
// Functions accept typeclass operations as parameters (not interfaces), making them
// work with any type that provides the right operations — including fp-go's Option,
// Either, IO, and user-defined types.
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// FunctorLaws verifies the Functor laws (identity and composition) for a type constructor.
//
// Type parameters:
//   - FA: the container type F[A]
//   - FB: the container type F[B]
//   - FC: the container type F[C]
//   - A, B, C: element types
//
// The fmap parameters correspond to the Map operation specialized for each type combination.
// In fp-go terms: fmapAB = option.Map[A, B], etc.
func FunctorLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	fmapAA func(func(A) A) func(FA) FA,
	fmapAB func(func(A) B) func(FA) FB,
	fmapBC func(func(B) C) func(FB) FC,
	fmapAC func(func(A) C) func(FA) FC,
	identity func(A) A,
	compose func(func(A) B, func(B) C) func(A) C,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	FunctorLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genAB), enginerapid.Wrap(genBC),
		eqFA, eqFC, fmapAA, fmapAB, fmapBC, fmapAC, identity, compose)
}
