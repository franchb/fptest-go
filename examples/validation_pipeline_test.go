package examples_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
	"github.com/franchb/fptest-go/assert"
	"github.com/franchb/fptest-go/prop"
	"pgregory.net/rapid"
)

// --- Validation functions ---

// normalizeEmail lowercases and trims an email address.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// validateUsername checks that a username is 3-20 alphanumeric characters.
func validateUsername(name string) either.Either[string, string] {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < 3 || len(trimmed) > 20 {
		return either.Left[string, string]("username must be 3-20 characters")
	}
	for _, r := range trimmed {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return either.Left[string, string]("username must be alphanumeric or underscore")
		}
	}
	return either.Right[string](trimmed)
}

// validateAge checks that an age is between 0 and 150.
func validateAge(age int) either.Either[string, int] {
	if age < 0 || age > 150 {
		return either.Left[int, string]("age must be between 0 and 150")
	}
	return either.Right[string](age)
}

// validateAgeOptimized is a "new" implementation that should behave identically.
func validateAgeOptimized(age int) either.Either[string, int] {
	if uint(age) > 150 {
		return either.Left[int, string]("age must be between 0 and 150")
	}
	return either.Right[string](age)
}

// maxAge returns the larger of two ages.
func maxAge(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- Tests ---

// TestNormalizeEmailIdempotent verifies that email normalization is idempotent:
// normalizing an already-normalized email produces the same result. This property
// is critical for deduplication — without it, storing the "normalized" form
// could still produce duplicates.
func TestNormalizeEmailIdempotent(t *testing.T) {
	prop.Idempotent(t, "NormalizeEmail",
		rapid.String(),
		func(a, b string) bool { return a == b },
		normalizeEmail,
	)
}

// TestMaxAgeCommutative verifies that max(a, b) == max(b, a). This property
// ensures that the order of comparison doesn't affect the result — important
// when aggregating ages from concurrent data sources.
func TestMaxAgeCommutative(t *testing.T) {
	prop.Commutative(t, "MaxAge",
		rapid.IntRange(0, 150),
		func(a, b int) bool { return a == b },
		maxAge,
	)
}

// TestValidateAgeOracle verifies that an optimized age validator produces the
// same results as the reference implementation. The "oracle" pattern is useful
// when rewriting performance-critical code — it proves behavioral equivalence
// across all inputs, not just a handful of examples.
func TestValidateAgeOracle(t *testing.T) {
	eqResult := either.Eq(eq.FromStrictEquals[string](), eq.FromStrictEquals[int]()).Equals
	prop.Oracle(t, "ValidateAge",
		rapid.IntRange(-1000, 1000),
		eqResult,
		validateAge,
		validateAgeOptimized,
	)
}

// TestValidUsernameInvariant verifies that all alphanumeric strings of valid
// length pass username validation. This is a domain invariant: if input matches
// the allowed character set and length, validation must succeed.
func TestValidUsernameInvariant(t *testing.T) {
	genValidUsername := rapid.StringMatching(`[a-zA-Z0-9_]{3,20}`)
	prop.Invariant(t, "ValidUsername",
		genValidUsername,
		func(name string) bool {
			return either.IsRight(validateUsername(name))
		},
	)
}

// TestValidationAssertions demonstrates using assert helpers to verify specific
// validation outcomes in scenario-based tests.
func TestValidationAssertions(t *testing.T) {
	t.Run("valid username returns Right", func(t *testing.T) {
		result := validateUsername("alice")
		name := assert.AssertRight[string](t, result)
		if name != "alice" {
			t.Fatalf("expected 'alice', got %q", name)
		}
	})

	t.Run("short username returns Left", func(t *testing.T) {
		result := validateUsername("ab")
		errMsg := assert.AssertLeft[string, string](t, result)
		if !strings.Contains(errMsg, "3-20") {
			t.Fatalf("unexpected error: %q", errMsg)
		}
	})

	t.Run("valid age returns Right", func(t *testing.T) {
		result := validateAge(25)
		age := assert.AssertRight[string](t, result)
		if age != 25 {
			t.Fatalf("expected 25, got %d", age)
		}
	})

	t.Run("negative age returns Left", func(t *testing.T) {
		result := validateAge(-1)
		assert.AssertLeft[string, int](t, result)
	})
}
