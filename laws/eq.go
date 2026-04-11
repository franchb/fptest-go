package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	"pgregory.net/rapid"
)

// EqLaws verifies the Eq laws (reflexivity, symmetry, transitivity).
func EqLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	EqLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), equals)
}

// EqInterfaceLaws verifies Eq laws using fp-go's Eq interface.
func EqInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA interface{ Equals(A, A) bool },
	opts ...Option,
) {
	t.Helper()
	EqLaws(t, genA, eqA.Equals, opts...)
}
