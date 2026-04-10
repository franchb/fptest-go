package hegel_test

import (
	"testing"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

func TestHegelRunnerExecutesProperty(t *testing.T) {
	var runner engine.Runner = fpthegel.HegelRunner{}
	executed := false
	runner.MakeCheck(t, "test/property", func(et engine.T) {
		executed = true
	})
	if !executed {
		t.Fatal("property was never executed")
	}
}

func TestHegelGenDrawsValues(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/draw", func(et engine.T) {
		gen := fpthegel.Wrap(hegellib.Integers[int](1, 100))
		val := gen.Draw(et, "n")
		if val < 1 || val > 100 {
			t.Fatalf("expected value in [1,100], got %d", val)
		}
	})
}
