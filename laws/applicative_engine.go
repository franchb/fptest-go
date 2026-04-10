package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// ApplicativeLawsEngine verifies the Applicative functor laws (identity and homomorphism)
// using the engine abstraction.
func ApplicativeLawsEngine[FA, FB, FAB, A, B any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	genFA engine.Generator[FA],
	genAB engine.Generator[func(A) B],
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

	runner.MakeCheck(t, "Applicative/Identity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")

		// Law: fmap(id)(v) == v (derived from ap(pure(id), v) == v)
		got := fmapAA(identity)(fa)
		if !eqFA(got, fa) {
			et.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	})

	runner.MakeCheck(t, "Applicative/Homomorphism", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")

		// Law: ap(pure(f), pure(x)) == pure(f(x))
		got := apAB(ofA(a))(ofAB(f))
		want := ofB(f(a))
		if !eqFB(got, want) {
			et.Fatalf("Applicative homomorphism violated:\n  ap(pure(f), pure(x)) = %v\n  pure(f(x))           = %v", got, want)
		}
	})
}

// ApplicativeInterchangeEngine verifies the Applicative interchange law
// using the engine abstraction.
func ApplicativeInterchangeEngine[FA, FB, FAB, FABB, A, B any](
	t *testing.T,
	runner engine.Runner,
	eqFB func(FB, FB) bool,
	apAB Applicative[A, B, FA, FB, FAB],
	ptdAB Pointed[func(A) B, FAB],
	apABB Apply[func(A) B, B, FAB, FB, FABB],
	ptdABB Pointed[func(func(A) B) B, FABB],
	genA engine.Generator[A],
	genAB engine.Generator[func(A) B],
) {
	t.Helper()

	runner.MakeCheck(t, "Applicative/Interchange", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")

		fab := ptdAB.Of(f) // F[func(A) B]

		// Left: Ap(Of(a))(fab)
		left := apAB.Ap(apAB.Of(a))(fab)

		// Right: Ap(fab)(Of(g => g(a)))
		callWithA := func(g func(A) B) B { return g(a) }
		right := apABB.Ap(fab)(ptdABB.Of(callWithA))

		if !eqFB(left, right) {
			et.Fatalf("Applicative interchange violated:\n  Ap(Of(a))(u)         = %v\n  Ap(u)(Of(f=>f(a)))   = %v", left, right)
		}
	})
}

// ApplicativeFullLawsEngine verifies all four Applicative functor laws using a
// pre-built [ApplicativeInstances] bundle and the engine abstraction.
func ApplicativeFullLawsEngine[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genFA engine.Generator[FA],
	genA engine.Generator[A],
	genAB engine.Generator[func(A) B],
	genBC engine.Generator[func(B) C],
	inst *ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C],
) {
	t.Helper()

	// Identity: fmap(id)(v) == v
	runner.MakeCheck(t, "Applicative/Identity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		got := inst.FmapAA(func(a A) A { return a })(fa)
		if !inst.EqFA(got, fa) {
			et.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	})

	// Homomorphism: ap(of(f))(of(a)) == of(f(a))
	runner.MakeCheck(t, "Applicative/Homomorphism", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")
		got := inst.ApAB.Ap(inst.ApAB.Of(a))(inst.PtdAB.Of(f))
		want := inst.PtdB.Of(f(a))
		if !inst.EqFB(got, want) {
			et.Fatalf("Applicative homomorphism violated:\n  ap(of(f))(of(a)) = %v\n  of(f(a))         = %v", got, want)
		}
	})

	// Interchange: ap(of(a))(u) == ap(u)(of(f => f(a)))
	runner.MakeCheck(t, "Applicative/Interchange", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")
		fab := inst.PtdAB.Of(f)
		left := inst.ApAB.Ap(inst.ApAB.Of(a))(fab)
		callWithA := func(g func(A) B) B { return g(a) }
		right := inst.ApABB.Ap(fab)(inst.PtdABB.Of(callWithA))
		if !inst.EqFB(left, right) {
			et.Fatalf("Applicative interchange violated:\n  ap(of(a))(u)         = %v\n  ap(u)(of(f=>f(a)))   = %v", left, right)
		}
	})

	// Composition: ap(ap(map(compose)(fbc))(fab))(fa) == ap(fbc)(ap(fab)(fa))
	runner.MakeCheck(t, "Apply/AssociativeComposition", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		ab := genAB.Draw(et, "ab")
		bc := genBC.Draw(et, "bc")
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
			et.Fatalf("Apply associative composition violated:\n  ap(ap(map(compose)(fbc))(fab))(fa) = %v\n  ap(fbc)(ap(fab)(fa))               = %v", left, right)
		}
	})
}
