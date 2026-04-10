package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// OrdLaws verifies the Ord laws (antisymmetry, transitivity, totality).
// The compare function must return -1, 0, or 1.
func OrdLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	OrdLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), equals, compare)
}

// OrdInterfaceLaws verifies Ord laws using fp-go's Ord interface.
func OrdInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	ordA interface {
		Equals(A, A) bool
		Compare(A, A) int
	},
	opts ...Option,
) {
	t.Helper()
	OrdLaws(t, genA, ordA.Equals, ordA.Compare, opts...)
}
