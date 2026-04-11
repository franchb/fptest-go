package hegel

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
	hegellib "hegel.dev/go/hegel"
)

// HegelRunner implements engine.Runner using hegel.Case.
type HegelRunner struct{}

// MakeCheck runs the property as a named subtest using hegel's property checker.
func (HegelRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
	t.Helper()
	t.Run(name, hegellib.Case(func(ht *hegellib.T) {
		prop(ht)
	}))
}
