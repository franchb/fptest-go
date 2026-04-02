package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/eq"
)

// AssertEq asserts that two values are equal according to the provided fp-go Eq instance.
func AssertEq[A any](t testing.TB, eqA eq.Eq[A], got, want A) {
	t.Helper()
	if !eqA.Equals(got, want) {
		t.Fatalf("Values not equal:\n  got  = %v\n  want = %v", got, want)
	}
}

// AssertEqFunc asserts that two values are equal according to the provided equality function.
func AssertEqFunc[A any](t testing.TB, equals func(A, A) bool, got, want A) {
	t.Helper()
	if !equals(got, want) {
		t.Fatalf("Values not equal:\n  got  = %v\n  want = %v", got, want)
	}
}

// AssertNotEq asserts that two values are not equal according to the provided Eq instance.
func AssertNotEq[A any](t testing.TB, eqA eq.Eq[A], got, notWant A) {
	t.Helper()
	if eqA.Equals(got, notWant) {
		t.Fatalf("Values should not be equal:\n  got = %v", got)
	}
}
