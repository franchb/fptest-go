package hegel

import (
	"github.com/franchb/fptest/engine"
	hegellib "hegel.dev/go/hegel"
)

// HegelGen wraps a hegel.Generator[A] as an engine.Generator[A].
type HegelGen[A any] struct {
	G hegellib.Generator[A]
}

// Draw produces a value from the wrapped hegel generator.
// The label parameter is ignored because hegel's Draw does not take a label.
func (hg HegelGen[A]) Draw(t engine.T, _ string) A {
	return hegellib.Draw(t.(*hegellib.T), hg.G)
}

// Wrap converts a hegel.Generator[A] into an engine.Generator[A].
func Wrap[A any](g hegellib.Generator[A]) engine.Generator[A] {
	return HegelGen[A]{G: g}
}
