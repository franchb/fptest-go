package gen

import (
	"github.com/IBM/fp-go/v2/either"
	"pgregory.net/rapid"
)

// GenEither generates an Either[E, A] that is randomly either Left or Right.
func GenEither[E, A any](genE *rapid.Generator[E], genA *rapid.Generator[A]) *rapid.Generator[either.Either[E, A]] {
	return rapid.Custom(func(t *rapid.T) either.Either[E, A] {
		if rapid.Bool().Draw(t, "isRight") {
			return either.Right[E](genA.Draw(t, "value"))
		}
		return either.Left[A](genE.Draw(t, "error"))
	})
}

// GenRight generates an Either[E, A] that is always Right.
func GenRight[E, A any](genA *rapid.Generator[A]) *rapid.Generator[either.Either[E, A]] {
	return rapid.Map(genA, either.Right[E, A])
}

// GenLeft generates an Either[E, A] that is always Left.
func GenLeft[E, A any](genE *rapid.Generator[E]) *rapid.Generator[either.Either[E, A]] {
	return rapid.Map(genE, either.Left[A, E])
}

// MonadicEither generates an Either[E, A] using the Gen monad.
func MonadicEither[E, A any](ge Gen[E], ga Gen[A]) Gen[either.Either[E, A]] {
	return func(t *rapid.T) either.Either[E, A] {
		if rapid.Bool().Draw(t, "isRight") {
			return either.Right[E](ga(t))
		}
		return either.Left[A](ge(t))
	}
}

// MonadicRight generates a Right Either[E, A] using the Gen monad.
func MonadicRight[E, A any](ga Gen[A]) Gen[either.Either[E, A]] {
	return Map(ga, either.Right[E, A])
}

// MonadicLeft generates a Left Either[E, A] using the Gen monad.
func MonadicLeft[E, A any](ge Gen[E]) Gen[either.Either[E, A]] {
	return Map(ge, either.Left[A, E])
}
