package gen

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"pgregory.net/rapid"
)

// GenIO generates an IO[A] that returns a generated value.
// The IO thunk captures the generated value, making it deterministic once drawn.
func GenIO[A any](genA *rapid.Generator[A]) *rapid.Generator[io.IO[A]] {
	return rapid.Map(genA, io.Of[A])
}

// GenIOEither generates an IOEither[E, A] from Either generators.
func GenIOEither[E, A any](genE *rapid.Generator[E], genA *rapid.Generator[A]) *rapid.Generator[ioeither.IOEither[E, A]] {
	return rapid.Map(GenEither[E, A](genE, genA), func(e either.Either[E, A]) ioeither.IOEither[E, A] {
		return func() either.Either[E, A] { return e }
	})
}

// GenIORight generates an IOEither[E, A] that always succeeds.
func GenIORight[E, A any](genA *rapid.Generator[A]) *rapid.Generator[ioeither.IOEither[E, A]] {
	return rapid.Map(genA, ioeither.Right[E, A])
}

// GenIOLeft generates an IOEither[E, A] that always fails.
func GenIOLeft[E, A any](genE *rapid.Generator[E]) *rapid.Generator[ioeither.IOEither[E, A]] {
	return rapid.Map(genE, ioeither.Left[A, E])
}

// MonadicIO generates an IO[A] using the Gen monad.
func MonadicIO[A any](ga Gen[A]) Gen[io.IO[A]] {
	return Map(ga, io.Of[A])
}
