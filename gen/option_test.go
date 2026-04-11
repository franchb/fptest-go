package gen_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest-go/gen"
	"pgregory.net/rapid"
)

func TestGenOption(t *testing.T) {
	// Just verify it produces valid Options without panicking
	rapid.Check(t, func(t *rapid.T) {
		o := gen.GenOption(rapid.Int()).Draw(t, "opt")
		// Should be either Some or None
		_ = option.IsSome(o) || option.IsNone(o)
	})
}

func TestGenSome(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		o := gen.GenSome(rapid.Int()).Draw(t, "opt")
		if !option.IsSome(o) {
			t.Fatal("GenSome produced None")
		}
	})
}

func TestGenNone(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		o := gen.GenNone[int]().Draw(t, "opt")
		if !option.IsNone(o) {
			t.Fatal("GenNone produced Some")
		}
	})
}
