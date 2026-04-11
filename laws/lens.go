package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	"pgregory.net/rapid"
)

// LensLaws verifies the Lens laws (get-set, set-get, set-set).
//
// A lawful lens must satisfy:
//   - GetSet: set(get(s), s) == s (setting what you get changes nothing)
//   - SetGet: get(set(a, s)) == a (you get back what you set)
//   - SetSet: set(b, set(a, s)) == set(b, s) (setting twice is setting once)
func LensLaws[S, A any](
	t *testing.T,
	genS *rapid.Generator[S],
	genA *rapid.Generator[A],
	eqS func(S, S) bool,
	eqA func(A, A) bool,
	get func(S) A,
	set func(A) func(S) S,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	LensLawsEngine(t, cfg.runner, enginerapid.Wrap(genS), enginerapid.Wrap(genA), eqS, eqA, get, set)
}
