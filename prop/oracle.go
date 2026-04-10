package prop

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	OracleEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqB, impl, reference)
}

// Idempotent verifies that applying a function twice yields the same result as applying once.
// f(f(x)) == f(x) for all x.
func Idempotent[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	IdempotentEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, f)
}

// Commutative verifies that f(a, b) == f(b, a) for all a, b.
func Commutative[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	CommutativeEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqB, f)
}

// Invariant verifies that a predicate holds for all generated inputs.
func Invariant[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	predicate func(A) bool,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	InvariantEngine(t, cfg.runner, name, enginerapid.Wrap(genA), predicate)
}
