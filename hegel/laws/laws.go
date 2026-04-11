// Package laws provides hegel-native convenience wrappers for algebraic law
// verification. Each function creates a HegelRunner and wraps hegel generators,
// then delegates to the corresponding engine-generic function in the core
// laws package.
package laws

import (
	"testing"

	fpthegel "github.com/franchb/fptest-go/hegel"
	corelaws "github.com/franchb/fptest-go/laws"
	hegellib "hegel.dev/go/hegel"
)

var hegelRunner = fpthegel.HegelRunner{}

// SemigroupLaws verifies the Semigroup law (associativity) using hegel generators.
func SemigroupLaws[A any](t *testing.T, genA hegellib.Generator[A], eqA func(A, A) bool, concat func(A, A) A) {
	t.Helper()
	corelaws.SemigroupLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), eqA, concat)
}

// MonoidLaws verifies the Monoid laws (left identity, right identity, associativity)
// using hegel generators.
func MonoidLaws[A any](t *testing.T, genA hegellib.Generator[A], eqA func(A, A) bool, concat func(A, A) A, empty A) {
	t.Helper()
	corelaws.MonoidLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), eqA, concat, empty)
}

// EqLaws verifies the Eq laws (reflexivity, symmetry, transitivity) using hegel generators.
func EqLaws[A any](t *testing.T, genA hegellib.Generator[A], equals func(A, A) bool) {
	t.Helper()
	corelaws.EqLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), equals)
}

// OrdLaws verifies the Ord laws (antisymmetry, transitivity, totality, consistency)
// using hegel generators.
func OrdLaws[A any](t *testing.T, genA hegellib.Generator[A], equals func(A, A) bool, compare func(A, A) int) {
	t.Helper()
	corelaws.OrdLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), equals, compare)
}

// LensLaws verifies the Lens laws (get-set, set-get, set-set) using hegel generators.
func LensLaws[S, A any](t *testing.T, genS hegellib.Generator[S], genA hegellib.Generator[A], eqS func(S, S) bool, eqA func(A, A) bool, get func(S) A, set func(A) func(S) S) {
	t.Helper()
	corelaws.LensLawsEngine(t, hegelRunner, fpthegel.Wrap(genS), fpthegel.Wrap(genA), eqS, eqA, get, set)
}

// FunctorLaws verifies the Functor laws (identity and composition) using hegel generators.
func FunctorLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genFA hegellib.Generator[FA],
	genAB hegellib.Generator[func(A) B],
	genBC hegellib.Generator[func(B) C],
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
	corelaws.FunctorLawsEngine(t, hegelRunner,
		fpthegel.Wrap(genFA), fpthegel.Wrap(genAB), fpthegel.Wrap(genBC),
		eqFA, eqFC, fmapAA, fmapAB, fmapBC, fmapAC, identity, compose)
}

// MonadLaws verifies the Monad laws (left identity, associativity) using hegel generators.
func MonadLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA hegellib.Generator[A],
	genFA hegellib.Generator[FA],
	genKleisliAB hegellib.Generator[func(A) FB],
	genKleisliBC hegellib.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()
	corelaws.MonadLawsEngine[FA, FB, FC, A, B, C](t, hegelRunner,
		fpthegel.Wrap(genA), fpthegel.Wrap(genFA),
		fpthegel.Wrap(genKleisliAB), fpthegel.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAB, chainBC, chainAC)
}

// MonadLawsFull verifies all Monad laws including right identity using hegel generators.
func MonadLawsFull[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA hegellib.Generator[A],
	genFA hegellib.Generator[FA],
	genKleisliAB hegellib.Generator[func(A) FB],
	genKleisliBC hegellib.Generator[func(B) FC],
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
	corelaws.MonadLawsFullEngine[FA, FB, FC, A, B, C](t, hegelRunner,
		fpthegel.Wrap(genA), fpthegel.Wrap(genFA),
		fpthegel.Wrap(genKleisliAB), fpthegel.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAA, chainAB, chainBC, chainAC)
}

// ChainAssociativity verifies the Chain associativity law using hegel generators.
func ChainAssociativity[FA, FB, FC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	genFA hegellib.Generator[FA],
	genKleisliAB hegellib.Generator[func(A) FB],
	genKleisliBC hegellib.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()
	corelaws.ChainAssociativityEngine[FA, FB, FC, A, B, C](t, hegelRunner, eqFC,
		fpthegel.Wrap(genFA), fpthegel.Wrap(genKleisliAB), fpthegel.Wrap(genKleisliBC),
		chainAB, chainBC, chainAC)
}
