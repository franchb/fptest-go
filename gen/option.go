package gen

import (
	"github.com/IBM/fp-go/v2/option"
	"pgregory.net/rapid"
)

// GenOption generates an Option[A] that is randomly either Some or None.
func GenOption[A any](genA *rapid.Generator[A]) *rapid.Generator[option.Option[A]] {
	return rapid.Custom(func(t *rapid.T) option.Option[A] {
		if rapid.Bool().Draw(t, "isSome") {
			return option.Some(genA.Draw(t, "value"))
		}
		return option.None[A]()
	})
}

// GenSome generates an Option[A] that is always Some.
func GenSome[A any](genA *rapid.Generator[A]) *rapid.Generator[option.Option[A]] {
	return rapid.Map(genA, option.Some[A])
}

// GenNone generates an Option[A] that is always None.
func GenNone[A any]() *rapid.Generator[option.Option[A]] {
	return rapid.Just(option.None[A]())
}

// MonadicOption generates an Option[A] using the Gen monad.
func MonadicOption[A any](ga Gen[A]) Gen[option.Option[A]] {
	return func(t *rapid.T) option.Option[A] {
		if rapid.Bool().Draw(t, "isSome") {
			return option.Some(ga(t))
		}
		return option.None[A]()
	}
}

// MonadicSome generates a Some Option[A] using the Gen monad.
func MonadicSome[A any](ga Gen[A]) Gen[option.Option[A]] {
	return Map(ga, option.Some[A])
}
