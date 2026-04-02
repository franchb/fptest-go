package laws

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
)

// EitherApplicativeInstances constructs an ApplicativeInstances bundle wired to
// fp-go's Either type. E, A, B, C must be comparable so that Either equality works.
func EitherApplicativeInstances[E, A, B, C comparable]() *ApplicativeInstances[
	either.Either[E, A],
	either.Either[E, B],
	either.Either[E, C],
	either.Either[E, func(A) B],
	either.Either[E, func(B) C],
	either.Either[E, func(A) C],
	either.Either[E, func(func(A) B) func(A) C],
	either.Either[E, func(func(A) B) B],
	A, B, C,
] {
	eqE := eq.FromStrictEquals[E]()
	return &ApplicativeInstances[
		either.Either[E, A],
		either.Either[E, B],
		either.Either[E, C],
		either.Either[E, func(A) B],
		either.Either[E, func(B) C],
		either.Either[E, func(A) C],
		either.Either[E, func(func(A) B) func(A) C],
		either.Either[E, func(func(A) B) B],
		A, B, C,
	]{
		EqFA: either.Eq(eqE, eq.FromStrictEquals[A]()).Equals,
		EqFB: either.Eq(eqE, eq.FromStrictEquals[B]()).Equals,
		EqFC: either.Eq(eqE, eq.FromStrictEquals[C]()).Equals,

		ApAB:   MakeApplicative[A, B](either.Of[E, A], either.Map[E, A, B], either.Ap[B, E, A]),
		PtdB:   MakePointed[B](either.Of[E, B]),
		FmapAA: either.Map[E, A, A],

		PtdAB:  MakePointed[func(A) B](either.Of[E, func(A) B]),
		ApABB:  MakeApply[func(A) B, B](either.Map[E, func(A) B, B], either.Ap[B, E, func(A) B]),
		PtdABB: MakePointed[func(func(A) B) B](either.Of[E, func(func(A) B) B]),

		PtdBC:       MakePointed[func(B) C](either.Of[E, func(B) C]),
		FmapCompose: MakeFunctor[func(B) C, func(func(A) B) func(A) C](either.Map[E, func(B) C, func(func(A) B) func(A) C]),
		ApBC:        MakeApply[B, C](either.Map[E, B, C], either.Ap[C, E, B]),
		ApAC:        MakeApply[A, C](either.Map[E, A, C], either.Ap[C, E, A]),
		ApABAC:      MakeApply[func(A) B, func(A) C](either.Map[E, func(A) B, func(A) C], either.Ap[func(A) C, E, func(A) B]),
	}
}
