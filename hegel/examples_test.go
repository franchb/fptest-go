package hegel_test

import (
	"math"
	"strconv"
	"testing"

	hegellaws "github.com/franchb/fptest/hegel/laws"
	hegelprop "github.com/franchb/fptest/hegel/prop"
	hegellib "hegel.dev/go/hegel"
)

// Example: Monoid law verification with hegel's Hypothesis-backed engine.
func TestExample_HegelMonoidLaws(t *testing.T) {
	hegellaws.MonoidLaws(t,
		hegellib.Text(0, 100),
		func(a, b string) bool { return a == b },
		func(a, b string) string { return a + b },
		"",
	)
}

// Example: Ord law verification with hegel.
func TestExample_HegelOrdLaws(t *testing.T) {
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

// Example: Round-trip property with hegel's integer generator.
func TestExample_HegelRoundTrip(t *testing.T) {
	hegelprop.RoundTrip(t, "itoa",
		hegellib.Integers[int](0, 9999),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		func(s string) int { n, _ := strconv.Atoi(s); return n },
	)
}

// Example: Using hegel's domain generators for invariant testing.
func TestExample_HegelDomainGenerators(t *testing.T) {
	hegelprop.Invariant(t, "emails-have-at",
		hegellib.Emails(),
		func(email string) bool {
			for _, c := range email {
				if c == '@' {
					return true
				}
			}
			return false
		},
	)
}
