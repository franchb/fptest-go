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

// ApplicativeInterchange verifies the Applicative interchange law:
//
//	Ap(Of(a))(u) == Ap(u)(Of(f => f(a)))
//
// Where u : F[func(A) B]. This says: applying a wrapped function to a pure value
// is the same as applying "call-with-that-value" to the wrapped function.
//
// Type parameters:
//   - FA, FB, FAB: container types F[A], F[B], F[func(A) B]
//   - FABB: container type F[func(func(A) B) B]
//   - A, B: element types
func ApplicativeInterchange[FA, FB, FAB, FABB, A, B any](
	t *testing.T,
	eqFB func(FB, FB) bool,
	apAB Applicative[A, B, FA, FB, FAB],
	ptdAB Pointed[func(A) B, FAB],
	apABB Apply[func(A) B, B, FAB, FB, FABB],
	ptdABB Pointed[func(func(A) B) B, FABB],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
) {
	t.Helper()

	t.Run("Applicative/Interchange", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")

		fab := ptdAB.Of(f) // F[func(A) B]

		// Left: Ap(Of(a))(fab)
		left := apAB.Ap(apAB.Of(a))(fab)

		// Right: Ap(fab)(Of(g => g(a)))
		callWithA := func(g func(A) B) B { return g(a) }
		right := apABB.Ap(fab)(ptdABB.Of(callWithA))

		if !eqFB(left, right) {
			t.Fatalf("Applicative interchange violated:\n  Ap(Of(a))(u)         = %v\n  Ap(u)(Of(f=>f(a)))   = %v", left, right)
		}
	}))
}

// ApplicativeFullLaws verifies all four Applicative functor laws using a
// pre-built [ApplicativeInstances] bundle:
//
//  1. Identity:     fmap(id)(v) == v
//  2. Homomorphism: ap(of(f))(of(a)) == of(f(a))
//  3. Interchange:  ap(of(a))(u) == ap(u)(of(f => f(a)))
//  4. Composition:  ap(ap(map(compose)(fbc))(fab))(fa) == ap(fbc)(ap(fab)(fa))
//
// The instances bundle avoids forcing callers to construct 12+ interface values
// manually.
func ApplicativeFullLaws[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	inst *ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C],
) {
	t.Helper()

	// Identity: fmap(id)(v) == v
	t.Run("Applicative/Identity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		got := inst.FmapAA(func(a A) A { return a })(fa)
		if !inst.EqFA(got, fa) {
			t.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	}))

	// Homomorphism: ap(of(f))(of(a)) == of(f(a))
	t.Run("Applicative/Homomorphism", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")
		got := inst.ApAB.Ap(inst.ApAB.Of(a))(inst.PtdAB.Of(f))
		want := inst.PtdB.Of(f(a))
		if !inst.EqFB(got, want) {
			t.Fatalf("Applicative homomorphism violated:\n  ap(of(f))(of(a)) = %v\n  of(f(a))         = %v", got, want)
		}
	}))

	// Interchange: ap(of(a))(u) == ap(u)(of(f => f(a)))
	t.Run("Applicative/Interchange", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")
		fab := inst.PtdAB.Of(f)
		left := inst.ApAB.Ap(inst.ApAB.Of(a))(fab)
		callWithA := func(g func(A) B) B { return g(a) }
		right := inst.ApABB.Ap(fab)(inst.PtdABB.Of(callWithA))
		if !inst.EqFB(left, right) {
			t.Fatalf("Applicative interchange violated:\n  ap(of(a))(u)         = %v\n  ap(u)(of(f=>f(a)))   = %v", left, right)
		}
	}))

	// Composition: ap(ap(map(compose)(fbc))(fab))(fa) == ap(fbc)(ap(fab)(fa))
	t.Run("Apply/AssociativeComposition", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		ab := genAB.Draw(t, "ab")
		bc := genBC.Draw(t, "bc")
		fab := inst.PtdAB.Of(ab)
		fbc := inst.PtdBC.Of(bc)
		compose := func(g func(B) C) func(func(A) B) func(A) C {
			return func(f func(A) B) func(A) C {
				return func(a A) C { return g(f(a)) }
			}
		}
		composed := inst.FmapCompose.Map(compose)(fbc)
		applied := inst.ApABAC.Ap(fab)(composed)
		left := inst.ApAC.Ap(fa)(applied)
		inner := inst.ApAB.Ap(fa)(fab)
		right := inst.ApBC.Ap(inner)(fbc)
		if !inst.EqFC(left, right) {
			t.Fatalf("Apply associative composition violated:\n  ap(ap(map(compose)(fbc))(fab))(fa) = %v\n  ap(fbc)(ap(fab)(fa))               = %v", left, right)
		}
	}))
}
