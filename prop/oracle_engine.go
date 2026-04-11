package prop

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
)

// OracleEngine verifies that an implementation produces the same results as a reference
// implementation for all generated inputs. This is the engine-generic variant of Oracle.
func OracleEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqB func(B, B) bool,
	impl func(A) B,
	reference func(A) B,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Oracle", func(et engine.T) {
		a := genA.Draw(et, "a")

		got := impl(a)
		want := reference(a)
		if !eqB(got, want) {
			et.Fatalf("Oracle mismatch for input %v:\n  impl      = %v\n  reference = %v", a, got, want)
		}
	})
}

// IdempotentEngine verifies that applying a function twice yields the same result as applying once.
// This is the engine-generic variant of Idempotent.
func IdempotentEngine[A any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Idempotent", func(et engine.T) {
		a := genA.Draw(et, "a")

		once := f(a)
		twice := f(once)
		if !eqA(once, twice) {
			et.Fatalf("Idempotency violated for input %v:\n  f(a)    = %v\n  f(f(a)) = %v", a, once, twice)
		}
	})
}

// CommutativeEngine verifies that f(a, b) == f(b, a) for all a, b.
// This is the engine-generic variant of Commutative.
func CommutativeEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Commutative", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		left := f(a, b)
		right := f(b, a)
		if !eqB(left, right) {
			et.Fatalf("Commutativity violated:\n  f(%v, %v) = %v\n  f(%v, %v) = %v", a, b, left, b, a, right)
		}
	})
}

// InvariantEngine verifies that a predicate holds for all generated inputs.
// This is the engine-generic variant of Invariant.
func InvariantEngine[A any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	predicate func(A) bool,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Invariant", func(et engine.T) {
		a := genA.Draw(et, "a")
		if !predicate(a) {
			et.Fatalf("Invariant violated for input: %v", a)
		}
	})
}
