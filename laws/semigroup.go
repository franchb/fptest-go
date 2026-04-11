package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	"pgregory.net/rapid"
)

// SemigroupLaws verifies the Semigroup law (associativity).
func SemigroupLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	SemigroupLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), eqA, concat)
}

// SemigroupInterfaceLaws verifies Semigroup laws using fp-go's Semigroup interface.
func SemigroupInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	sg interface{ Concat(A, A) A },
	opts ...Option,
) {
	t.Helper()
	SemigroupLaws(t, genA, eqA, sg.Concat, opts...)
}
