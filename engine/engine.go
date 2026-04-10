// Package engine defines the PBT engine abstraction for fptest-go.
//
// This package contains only interfaces — no external dependencies.
// Concrete implementations live in engine/rapid/ (for pgregory.net/rapid)
// and hegel/ (for hegel.dev/go/hegel).
package engine

import (
	"context"
	"testing"
)

// T is the test context provided by a PBT engine during property execution.
// Both rapid.T and hegel.T satisfy this interface.
//
// This is intentionally narrower than testing.TB because not all PBT engines
// implement every testing.TB method (e.g. rapid.T lacks Helper and Name).
type T interface {
	Cleanup(func())
	Context() context.Context
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Log(args ...any)
	Logf(format string, args ...any)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
}

// Generator draws values of type A from the PBT engine's search space.
type Generator[A any] interface {
	Draw(t T, label string) A
}

// Runner executes a property check as a named subtest.
type Runner interface {
	MakeCheck(t *testing.T, name string, prop func(T))
}
