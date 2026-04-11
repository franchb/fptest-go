package gen_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/franchb/fptest-go/gen"
	"pgregory.net/rapid"
)

func TestGenEither(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		e := gen.GenEither(rapid.String(), rapid.Int()).Draw(t, "either")
		_ = either.IsLeft(e) || either.IsRight(e)
	})
}

func TestGenRight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		e := gen.GenRight[string](rapid.Int()).Draw(t, "right")
		if !either.IsRight(e) {
			t.Fatal("GenRight produced Left")
		}
	})
}

func TestGenLeft(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		e := gen.GenLeft[string, int](rapid.String()).Draw(t, "left")
		if !either.IsLeft(e) {
			t.Fatal("GenLeft produced Right")
		}
	})
}
