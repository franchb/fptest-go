// Package rapid provides the rapid PBT engine adapter for fptest-go.
package rapid

import (
	"testing"

	"github.com/franchb/fptest/engine"
	rapidlib "pgregory.net/rapid"
)

// RapidRunner implements engine.Runner using rapid.MakeCheck.
type RapidRunner struct{}

// MakeCheck runs the property as a named subtest using rapid's property checker.
func (RapidRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
	t.Helper()
	t.Run(name, rapidlib.MakeCheck(func(rt *rapidlib.T) {
		prop(rt)
	}))
}

// RapidGen wraps a *rapid.Generator[A] as an engine.Generator[A].
type RapidGen[A any] struct {
	G *rapidlib.Generator[A]
}

// Draw produces a value from the wrapped rapid generator.
func (rg RapidGen[A]) Draw(t engine.T, label string) A {
	return rg.G.Draw(t.(*rapidlib.T), label)
}

// Wrap converts a *rapid.Generator[A] into an engine.Generator[A].
func Wrap[A any](g *rapidlib.Generator[A]) engine.Generator[A] {
	return RapidGen[A]{G: g}
}
