package gen_test

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
	enginerapid "github.com/franchb/fptest-go/engine/rapid"
	"github.com/franchb/fptest-go/gen"
	"pgregory.net/rapid"
)

func TestToEngineDrawsValues(t *testing.T) {
	// Use FromRapid(rapid.Just(...)) for constant values inside rapid,
	// because gen.Of does not consume bitstream data (see ToRapid doc).
	g := gen.FromRapid(rapid.Just(42))
	var eg engine.Generator[int] = gen.ToEngine(g)

	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/ToEngine", func(et engine.T) {
		val := eg.Draw(et, "n")
		if val != 42 {
			t.Fatalf("expected 42, got %d", val)
		}
	})
}

func TestToEngineWithRapidGen(t *testing.T) {
	g := gen.FromRapid(rapid.IntRange(10, 20))
	eg := gen.ToEngine(g)

	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/ToEngineRapid", func(et engine.T) {
		val := eg.Draw(et, "n")
		if val < 10 || val > 20 {
			t.Fatalf("expected value in [10,20], got %d", val)
		}
	})
}
