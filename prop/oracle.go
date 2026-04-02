package prop

import (
	"testing"

	"pgregory.net/rapid"
)

// Oracle verifies that an implementation produces the same results as a reference
// implementation for all generated inputs. This is the "test oracle" pattern.
func Oracle[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqB func(B, B) bool,
	impl func(A) B,
	reference func(A) B,
) {
	t.Helper()
	t.Run(name+"/Oracle", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		got := impl(a)
		want := reference(a)
		if !eqB(got, want) {
			t.Fatalf("Oracle mismatch for input %v:\n  impl      = %v\n  reference = %v", a, got, want)
		}
	}))
}

// Idempotent verifies that applying a function twice yields the same result as applying once.
// f(f(x)) == f(x) for all x.
func Idempotent[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
) {
	t.Helper()
	t.Run(name+"/Idempotent", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")

		once := f(a)
		twice := f(once)
		if !eqA(once, twice) {
			t.Fatalf("Idempotency violated for input %v:\n  f(a)    = %v\n  f(f(a)) = %v", a, once, twice)
		}
	}))
}

// Commutative verifies that f(a, b) == f(b, a) for all a, b.
func Commutative[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
) {
	t.Helper()
	t.Run(name+"/Commutative", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		left := f(a, b)
		right := f(b, a)
		if !eqB(left, right) {
			t.Fatalf("Commutativity violated:\n  f(%v, %v) = %v\n  f(%v, %v) = %v", a, b, left, b, a, right)
		}
	}))
}

// Invariant verifies that a predicate holds for all generated inputs.
func Invariant[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	predicate func(A) bool,
) {
	t.Helper()
	t.Run(name+"/Invariant", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		if !predicate(a) {
			t.Fatalf("Invariant violated for input: %v", a)
		}
	}))
}
