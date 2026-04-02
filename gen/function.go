package gen

import (
	"pgregory.net/rapid"
)

// GenFunc generates a pure function from A to B by generating a B and ignoring the input.
// This produces "constant functions" useful for law testing where the function's behavior
// must be consistent across invocations with the same *rapid.T draw sequence.
func GenFunc[A, B any](genB *rapid.Generator[B]) *rapid.Generator[func(A) B] {
	return rapid.Map(genB, func(b B) func(A) B {
		return func(_ A) B { return b }
	})
}

// GenEndomorphism generates a function from A to A. For numeric types, this generates
// simple arithmetic transformations. For general types, it generates constant functions.
func GenEndomorphism[A any](genA *rapid.Generator[A]) *rapid.Generator[func(A) A] {
	return GenFunc[A, A](genA)
}

// MonadicFunc generates a function from A to B using the Gen monad.
func MonadicFunc[A, B any](gb Gen[B]) Gen[func(A) B] {
	return Map(gb, func(b B) func(A) B {
		return func(_ A) B { return b }
	})
}
