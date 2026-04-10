// Package engine defines the PBT engine abstraction for fptest-go.
//
// This package contains only interfaces — no external dependencies.
// Concrete implementations live in engine/rapid/ (for pgregory.net/rapid)
// and hegel/ (for hegel.dev/go/hegel).
package engine

import "testing"

// T is the test context provided by a PBT engine during property execution.
// Both rapid.T and hegel.T satisfy this interface.
type T interface {
	testing.TB
}

// Generator draws values of type A from the PBT engine's search space.
type Generator[A any] interface {
	Draw(t T, label string) A
}

// Runner executes a property check as a named subtest.
type Runner interface {
	MakeCheck(t *testing.T, name string, prop func(T))
}
