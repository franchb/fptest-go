package rapid_test

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	rapidlib "pgregory.net/rapid"
)

func TestRapidRunnerExecutesProperty(t *testing.T) {
	var runner engine.Runner = enginerapid.RapidRunner{}
	executed := false
	runner.MakeCheck(t, "test/property", func(et engine.T) {
		executed = true
	})
	if !executed {
		t.Fatal("property was never executed")
	}
}

func TestRapidGenDrawsValues(t *testing.T) {
	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/draw", func(et engine.T) {
		gen := enginerapid.Wrap(rapidlib.IntRange(1, 100))
		val := gen.Draw(et, "n")
		if val < 1 || val > 100 {
			t.Fatalf("expected value in [1,100], got %d", val)
		}
	})
}
