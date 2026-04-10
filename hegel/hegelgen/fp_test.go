package hegelgen_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	"github.com/franchb/fptest/hegel/hegelgen"
	hegellib "hegel.dev/go/hegel"
)

func TestOptionGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	hasSome := false
	hasNone := false
	runner.MakeCheck(t, "test/option", func(et engine.T) {
		opt := hegelgen.Option(hegellib.Integers[int](1, 100)).Draw(et, "opt")
		option.Fold(func() string {
			hasNone = true
			return ""
		}, func(v int) string {
			hasSome = true
			if v < 1 || v > 100 {
				t.Fatalf("expected value in [1,100], got %d", v)
			}
			return ""
		})(opt)
	})
	// Over 100 iterations, we should see both Some and None.
	if !hasSome {
		t.Error("never generated Some")
	}
	if !hasNone {
		t.Error("never generated None")
	}
}

func TestSomeGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/some", func(et engine.T) {
		opt := hegelgen.Some(hegellib.Integers[int](1, 100)).Draw(et, "opt")
		if option.IsNone(opt) {
			t.Fatal("expected Some, got None")
		}
	})
}

func TestNoneGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/none", func(et engine.T) {
		opt := hegelgen.None[int]().Draw(et, "opt")
		if option.IsSome(opt) {
			t.Fatal("expected None, got Some")
		}
	})
}
