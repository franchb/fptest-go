package laws

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/option"
)

// ApplicativeInstances bundles all typed interface instances needed to test the
// full Applicative law suite (identity, homomorphism, interchange, composition).
type ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any] struct {
	// Equality
	EqFA func(FA, FA) bool
	EqFB func(FB, FB) bool
	EqFC func(FC, FC) bool

	// For identity + homomorphism
	ApAB   Applicative[A, B, FA, FB, FAB]
	PtdB   Pointed[B, FB]
	FmapAA func(func(A) A) func(FA) FA // for identity law: fmap(id)

	// For interchange
	PtdAB  Pointed[func(A) B, FAB]
	ApABB  Apply[func(A) B, B, FAB, FB, FABB]
	PtdABB Pointed[func(func(A) B) B, FABB]

	// For composition (via Apply)
	PtdBC       Pointed[func(B) C, FBC]
	FmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC]
	ApBC        Apply[B, C, FB, FC, FBC]
	ApAC        Apply[A, C, FA, FC, FAC]
	ApABAC      Apply[func(A) B, func(A) C, FAB, FAC, FABAC]
}

// OptionApplicativeInstances constructs an ApplicativeInstances bundle wired to
// fp-go's Option type. A, B, C must be comparable so that Option equality works.
func OptionApplicativeInstances[A, B, C comparable]() *ApplicativeInstances[
	option.Option[A], option.Option[B], option.Option[C],
	option.Option[func(A) B], option.Option[func(B) C],
	option.Option[func(A) C],
	option.Option[func(func(A) B) func(A) C],
	option.Option[func(func(A) B) B],
	A, B, C,
] {
	return &ApplicativeInstances[
		option.Option[A], option.Option[B], option.Option[C],
		option.Option[func(A) B], option.Option[func(B) C],
		option.Option[func(A) C],
		option.Option[func(func(A) B) func(A) C],
		option.Option[func(func(A) B) B],
		A, B, C,
	]{
		EqFA: option.Eq(eq.FromStrictEquals[A]()).Equals,
		EqFB: option.Eq(eq.FromStrictEquals[B]()).Equals,
		EqFC: option.Eq(eq.FromStrictEquals[C]()).Equals,

		ApAB:   MakeApplicative[A, B](option.Of[A], option.Map[A, B], option.Ap[B, A]),
		PtdB:   MakePointed[B](option.Of[B]),
		FmapAA: option.Map[A, A],

		PtdAB:  MakePointed[func(A) B](option.Of[func(A) B]),
		ApABB:  MakeApply[func(A) B, B](option.Map[func(A) B, B], option.Ap[B, func(A) B]),
		PtdABB: MakePointed[func(func(A) B) B](option.Of[func(func(A) B) B]),

		PtdBC:       MakePointed[func(B) C](option.Of[func(B) C]),
		FmapCompose: MakeFunctor[func(B) C, func(func(A) B) func(A) C](option.Map[func(B) C, func(func(A) B) func(A) C]),
		ApBC:        MakeApply[B, C](option.Map[B, C], option.Ap[C, B]),
		ApAC:        MakeApply[A, C](option.Map[A, C], option.Ap[C, A]),
		ApABAC:      MakeApply[func(A) B, func(A) C](option.Map[func(A) B, func(A) C], option.Ap[func(A) C, func(A) B]),
	}
}
