package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// ApplicativeLaws verifies the Applicative functor laws.
//
// Due to Go's type system limitations (no HKTs), this tests the most important
// Applicative properties: identity via fmap and homomorphism via ap.
//
// Type parameters:
//   - FA, FB, FAB: container types F[A], F[B], F[func(A) B]
//   - A, B: element types
func ApplicativeLaws[FA, FB, FAB, A, B any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	eqFA func(FA, FA) bool,
	eqFB func(FB, FB) bool,
	ofA func(A) FA,
	ofB func(B) FB,
	ofAB func(func(A) B) FAB,
	fmapAA func(func(A) A) func(FA) FA,
	apAB func(FA) func(FAB) FB,
	identity func(A) A,
) {
	t.Helper()

	t.Run("Applicative/Identity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")

		// Law: fmap(id)(v) == v (derived from ap(pure(id), v) == v)
		got := fmapAA(identity)(fa)
		if !eqFA(got, fa) {
			t.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	}))

	t.Run("Applicative/Homomorphism", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")

		// Law: ap(pure(f), pure(x)) == pure(f(x))
		got := apAB(ofA(a))(ofAB(f))
		want := ofB(f(a))
		if !eqFB(got, want) {
			t.Fatalf("Applicative homomorphism violated:\n  ap(pure(f), pure(x)) = %v\n  pure(f(x))           = %v", got, want)
		}
	}))
}
