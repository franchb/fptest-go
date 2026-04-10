package laws_test

import (
	"math"
	"testing"

	hegellaws "github.com/franchb/fptest/hegel/laws"
	hegellib "hegel.dev/go/hegel"
)

func TestStringMonoidLaws_Hegel(t *testing.T) {
	hegellaws.MonoidLaws(t,
		hegellib.Text(0, 50),
		func(a, b string) bool { return a == b },
		func(a, b string) string { return a + b },
		"",
	)
}

func TestIntSemigroupLaws_Hegel(t *testing.T) {
	hegellaws.SemigroupLaws(t,
		hegellib.Integers[int](math.MinInt, math.MaxInt),
		func(a, b int) bool { return a == b },
		func(a, b int) int { return a + b },
	)
}

func TestIntEqLaws_Hegel(t *testing.T) {
	hegellaws.EqLaws(t,
		hegellib.Integers[int](math.MinInt, math.MaxInt),
		func(a, b int) bool { return a == b },
	)
}

func TestIntOrdLaws_Hegel(t *testing.T) {
	hegellaws.OrdLaws(t,
		hegellib.Integers[int](math.MinInt, math.MaxInt),
		func(a, b int) bool { return a == b },
		func(a, b int) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		},
	)
}
