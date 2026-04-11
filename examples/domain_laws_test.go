package examples_test

import (
	"strings"
	"testing"

	"github.com/franchb/fptest-go/laws"
	"pgregory.net/rapid"
)

// --- Domain types ---

// Money represents a monetary amount in the smallest unit (e.g., cents).
// Semigroup/Monoid only valid for same-currency values.
type Money struct {
	Amount   int64
	Currency string
}

// Email is a string wrapper with case-insensitive equality.
type Email string

// Priority represents task urgency levels.
type Priority int

const (
	PriorityLow      Priority = 1
	PriorityMedium   Priority = 2
	PriorityHigh     Priority = 3
	PriorityCritical Priority = 4
)

// --- Generators ---

func genUSDMoney() *rapid.Generator[Money] {
	return rapid.Custom(func(t *rapid.T) Money {
		return Money{
			Amount:   rapid.Int64Range(0, 100_000_00).Draw(t, "cents"),
			Currency: "USD",
		}
	})
}

func genEmail() *rapid.Generator[Email] {
	return rapid.Custom(func(t *rapid.T) Email {
		// Generate emails with random casing to exercise case-insensitive equality
		user := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9]{1,8}`).Draw(t, "user")
		domain := rapid.StringMatching(`[a-zA-Z]{2,6}`).Draw(t, "domain")
		tld := rapid.SampledFrom([]string{"com", "org", "net", "IO", "Dev"}).Draw(t, "tld")
		return Email(user + "@" + domain + "." + tld)
	})
}

func genPriority() *rapid.Generator[Priority] {
	return rapid.SampledFrom([]Priority{
		PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical,
	})
}

// --- Semigroup & Monoid for Money ---

// TestMoneySemigroupLaws verifies that adding same-currency Money values is
// associative: (a + b) + c == a + (b + c). This is the fundamental property
// that makes it safe to split and recombine money totals (e.g., parallel
// aggregation of transaction amounts).
func TestMoneySemigroupLaws(t *testing.T) {
	laws.SemigroupLaws(t,
		genUSDMoney(),
		func(a, b Money) bool { return a == b },
		func(a, b Money) Money {
			return Money{Amount: a.Amount + b.Amount, Currency: a.Currency}
		},
	)
}

// TestMoneyMonoidLaws verifies that zero-amount money is a valid identity
// element: 0 + x == x and x + 0 == x. Together with associativity, this
// guarantees that folding an empty list of transactions yields zero, and that
// prepending/appending zero-value transactions is harmless.
func TestMoneyMonoidLaws(t *testing.T) {
	zero := Money{Amount: 0, Currency: "USD"}

	laws.MonoidLaws(t,
		genUSDMoney(),
		func(a, b Money) bool { return a == b },
		func(a, b Money) Money {
			return Money{Amount: a.Amount + b.Amount, Currency: a.Currency}
		},
		zero,
	)
}

// --- Eq for Email ---

// TestEmailEqLaws verifies that case-insensitive email equality is a valid
// equivalence relation: reflexive, symmetric, and transitive. Without these
// properties, using Email as a map key or deduplication criterion would be
// unreliable.
func TestEmailEqLaws(t *testing.T) {
	laws.EqLaws(t,
		genEmail(),
		func(a, b Email) bool {
			return strings.EqualFold(string(a), string(b))
		},
	)
}

// --- Ord for Priority ---

// TestPriorityOrdLaws verifies that Priority ordering satisfies antisymmetry,
// transitivity, and totality. These properties guarantee correct behavior when
// sorting tasks by priority or using priority queues.
func TestPriorityOrdLaws(t *testing.T) {
	laws.OrdLaws(t,
		genPriority(),
		func(a, b Priority) bool { return a == b },
		func(a, b Priority) int {
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
