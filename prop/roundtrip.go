// Package prop provides higher-level property testing utilities built on rapid.
package prop

import (
	"testing"

	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	"pgregory.net/rapid"
)

// RoundTrip verifies that encode and decode are inverses: decode(encode(a)) == a.
// This is the fundamental property for serialization, parsing, and any codec pair.
func RoundTrip[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}

// RoundTripPartial verifies a round-trip where decode may fail, returning (A, bool).
// Only checks the property when decoding succeeds.
func RoundTripPartial[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripPartialEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}

// RoundTripError verifies a round-trip where decode may return an error.
func RoundTripError[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripErrorEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}
