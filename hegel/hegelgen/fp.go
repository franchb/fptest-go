package hegelgen

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

// Option returns a generator that produces Option[A] values (Some or None).
func Option[A any](genA hegellib.Generator[A]) engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Map(
		hegellib.Optional(genA),
		func(ptr *A) option.Option[A] {
			if ptr == nil {
				return option.None[A]()
			}
			return option.Some(*ptr)
		},
	))
}

// Some returns a generator that always produces Some[A].
func Some[A any](genA hegellib.Generator[A]) engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Map(genA, option.Some[A]))
}

// None returns a generator that always produces None[A].
func None[A any]() engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Just(option.None[A]()))
}

// Either returns a generator that produces Either[E, A] (Left or Right).
func Either[E, A any](genE hegellib.Generator[E], genA hegellib.Generator[A]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.FlatMap(hegellib.Booleans(), func(isRight bool) hegellib.Generator[either.Either[E, A]] {
		if isRight {
			return hegellib.Map(genA, either.Right[E, A])
		}
		return hegellib.Map(genE, either.Left[A, E])
	}))
}

// Right returns a generator that always produces Right[E, A].
func Right[E, A any](genA hegellib.Generator[A]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.Map(genA, either.Right[E, A]))
}

// Left returns a generator that always produces Left[E, A].
func Left[E, A any](genE hegellib.Generator[E]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.Map(genE, either.Left[A, E]))
}
