package prop_test

import (
	"strconv"
	"testing"

	"github.com/franchb/fptest/prop"
	"pgregory.net/rapid"
)

func TestRoundTrip(t *testing.T) {
	prop.RoundTrip(t, "IntToString",
		rapid.IntRange(-1000, 1000),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		},
	)
}

func TestRoundTripError(t *testing.T) {
	prop.RoundTripError(t, "IntToString",
		rapid.IntRange(-1000, 1000),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		strconv.Atoi,
	)
}

func TestOracle(t *testing.T) {
	// Test that our multiply-by-2 is the same as adding to itself
	prop.Oracle(t, "Double",
		rapid.Int(),
		func(a, b int) bool { return a == b },
		func(x int) int { return x * 2 },
		func(x int) int { return x + x },
	)
}

func TestIdempotent(t *testing.T) {
	// abs is idempotent
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	prop.Idempotent(t, "Abs",
		rapid.Int(),
		func(a, b int) bool { return a == b },
		abs,
	)
}

func TestCommutative(t *testing.T) {
	prop.Commutative(t, "IntAdd",
		rapid.Int(),
		func(a, b int) bool { return a == b },
		func(a, b int) int { return a + b },
	)
}

func TestInvariant(t *testing.T) {
	prop.Invariant(t, "SliceLenNonNegative",
		rapid.SliceOf(rapid.Int()),
		func(s []int) bool { return len(s) >= 0 },
	)
}
