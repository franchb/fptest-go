package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// MonoidLaws verifies the Monoid laws (left identity, right identity, associativity).
// Monoid extends Semigroup with an identity element.
func MonoidLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonoidLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), eqA, concat, empty)
}

// MonoidInterfaceLaws verifies Monoid laws using fp-go's Monoid interface.
func MonoidInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	m interface {
		Concat(A, A) A
		Empty() A
	},
	opts ...Option,
) {
	t.Helper()
	MonoidLaws(t, genA, eqA, m.Concat, m.Empty(), opts...)
}
