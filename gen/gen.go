// Package gen provides a monadic generator abstraction built on top of pgregory.net/rapid.
//
// The core type Gen[A] is defined as func(*rapid.T) A, which is isomorphic to
// Reader[*rapid.T, A] in fp-go's type system. This gives generators a full
// Functor/Applicative/Monad structure, enabling declarative composition of
// property-based test data generators.
package gen

import (
	"pgregory.net/rapid"
)

// Gen is a generator of values of type A. It is isomorphic to Reader[*rapid.T, A],
// making it a monad that threads rapid's random source through composed generators.
type Gen[A any] func(*rapid.T) A

// Of lifts a pure value into the Gen monad (Pointed/Applicative pure).
func Of[A any](a A) Gen[A] {
	return func(_ *rapid.T) A { return a }
}

// Map applies a function to the generated value (Functor fmap).
func Map[A, B any](ga Gen[A], f func(A) B) Gen[B] {
	return func(t *rapid.T) B { return f(ga(t)) }
}

// Chain sequences two generators where the second depends on the first (Monad bind/flatMap).
// This enables dependent generation — the Go equivalent of Hypothesis's flatmap.
func Chain[A, B any](ga Gen[A], f func(A) Gen[B]) Gen[B] {
	return func(t *rapid.T) B { return f(ga(t))(t) }
}

// Ap applies a generated function to a generated value (Applicative ap).
func Ap[A, B any](gf Gen[func(A) B], ga Gen[A]) Gen[B] {
	return func(t *rapid.T) B { return gf(t)(ga(t)) }
}

// Filter creates a generator that only produces values satisfying the predicate.
// Panics if too many values are rejected (same semantics as rapid's Filter).
func Filter[A any](ga Gen[A], pred func(A) bool) Gen[A] {
	return func(t *rapid.T) A {
		return ToRapid(ga).Filter(pred).Draw(t, "filtered")
	}
}

// ToRapid converts a Gen[A] to rapid's *Generator[A] for use with rapid.Check
// and other rapid combinators.
//
// Note: the underlying Gen must draw from at least one rapid generator during
// execution, as rapid requires Custom generators to consume bitstream data.
// Use rapid.Just for constant values instead.
func ToRapid[A any](g Gen[A]) *rapid.Generator[A] {
	return rapid.Custom(func(t *rapid.T) A { return g(t) })
}

// FromRapid converts a rapid *Generator[A] into a Gen[A].
func FromRapid[A any](g *rapid.Generator[A]) Gen[A] {
	return func(t *rapid.T) A { return g.Draw(t, "") }
}

// FromRapidLabeled converts a rapid *Generator[A] into a Gen[A] with a label
// for debugging shrink output.
func FromRapidLabeled[A any](g *rapid.Generator[A], label string) Gen[A] {
	return func(t *rapid.T) A { return g.Draw(t, label) }
}

// Map2 combines two generators using a binary function (Applicative lift2).
func Map2[A, B, C any](ga Gen[A], gb Gen[B], f func(A, B) C) Gen[C] {
	return func(t *rapid.T) C { return f(ga(t), gb(t)) }
}

// Map3 combines three generators using a ternary function (Applicative lift3).
func Map3[A, B, C, D any](ga Gen[A], gb Gen[B], gc Gen[C], f func(A, B, C) D) Gen[D] {
	return func(t *rapid.T) D { return f(ga(t), gb(t), gc(t)) }
}

// Pair generates a pair of values from two independent generators.
func Pair[A, B any](ga Gen[A], gb Gen[B]) Gen[[2]any] {
	return func(t *rapid.T) [2]any { return [2]any{ga(t), gb(t)} }
}

// Slice generates a slice of values from a single generator with length in [minLen, maxLen].
func Slice[A any](ga Gen[A], minLen, maxLen int) Gen[[]A] {
	return func(t *rapid.T) []A {
		n := rapid.IntRange(minLen, maxLen).Draw(t, "len")
		result := make([]A, n)
		for i := range result {
			result[i] = ga(t)
		}
		return result
	}
}
