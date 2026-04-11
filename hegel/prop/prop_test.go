package prop_test

import (
	"math"
	"strconv"
	"testing"

	hegelprop "github.com/franchb/fptest-go/hegel/prop"
	hegellib "hegel.dev/go/hegel"
)

func TestRoundTrip_Hegel(t *testing.T) {
	hegelprop.RoundTrip(t, "itoa",
		hegellib.Integers[int](0, 999),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		func(s string) int { n, _ := strconv.Atoi(s); return n },
	)
}

func TestCommutative_Hegel(t *testing.T) {
	hegelprop.Commutative(t, "add",
		hegellib.Integers[int](math.MinInt/2, math.MaxInt/2),
		func(a, b int) bool { return a == b },
		func(a, b int) int { return a + b },
	)
}

func TestIdempotent_Hegel(t *testing.T) {
	hegelprop.Idempotent(t, "abs",
		hegellib.Integers[int](0, math.MaxInt),
		func(a, b int) bool { return a == b },
		func(a int) int {
			if a < 0 {
				return -a
			}
			return a
		},
	)
}
