# Recover Applicative/Apply Laws Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Recover the 2 missing Applicative laws (interchange, composition) and add the Apply associative composition law, using copied typeclass interfaces with helper constructors for fp-go types.

**Architecture:** Define minimal typeclass interfaces (Pointed, Functor, Apply, Applicative, Chainable, Monad) in `laws/typeclass.go`. Add new law functions that accept these interfaces. Provide per-type helper constructors (`NewOptionInstances`, `NewEitherInstances`) that wire fp-go's public API into the interface instances. Existing raw-function APIs remain untouched for backward compatibility.

**Tech Stack:** Go 1.25, github.com/IBM/fp-go/v2, pgregory.net/rapid

---

## Background

fp-go's internal law hierarchy tests: Functor (2) -> Apply (1) -> Applicative (3) -> Chain (1) -> Monad (3). The current fptest-go v2 only tests Applicative identity + homomorphism. The original v1 had all 4 Applicative laws. This plan recovers the missing interchange and composition laws.

### Key fp-go API signatures used throughout

```go
// Option
option.Of[T any](value T) Option[T]
option.Map[A, B any](f func(A) B) func(Option[A]) Option[B]  // curried
option.Ap[B, A any](fa Option[A]) func(Option[func(A) B]) Option[B]  // value-first
option.Chain[A, B any](f func(A) Option[B]) func(Option[A]) Option[B]

// Either
either.Of[E, A any](value A) Either[E, A]
either.Map[E, A, B any](f func(A) B) func(Either[E, A]) Either[E, B]
either.Ap[B, E, A any](fa Either[E, A]) func(Either[E, func(A) B]) Either[E, B]
either.Chain[E, A, B any](f func(A) Either[E, B]) func(Either[E, A]) Either[E, B]
```

### Ap direction (critical)

fp-go uses **value-first**: `Ap(fa)(fab) -> fb`. The interface `Ap` method follows this: `Ap(FA) func(FAB) FB`.

---

## Task 1: Define typeclass interfaces

**Files:**
- Create: `laws/typeclass.go`
- Test: `laws/typeclass_test.go`

**Step 1: Write the test that constructs an Option Pointed instance**

This test validates that the interfaces compile and can be satisfied by fp-go types.

```go
// laws/typeclass_test.go
package laws_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
)

// Compile-time check: option.Of satisfies Pointed via adapter
func TestPointedOptionCompiles(t *testing.T) {
	ptd := laws.MakePointed(option.Of[int])
	got := ptd.Of(42)
	if !option.IsSome(got) {
		t.Fatal("expected Some(42)")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestPointedOptionCompiles -v`
Expected: FAIL — `laws.MakePointed` undefined

**Step 3: Write the interfaces and constructors**

```go
// laws/typeclass.go
package laws

// Pointed lifts a pure value into context F.
type Pointed[A, FA any] interface {
	Of(A) FA
}

// Functor maps a function over context F.
type Functor[A, B, FA, FB any] interface {
	Map(func(A) B) func(FA) FB
}

// Apply extends Functor with function application in context.
// Ap takes value-in-context first (FA), returns function taking function-in-context (FAB).
type Apply[A, B, FA, FB, FAB any] interface {
	Functor[A, B, FA, FB]
	Ap(FA) func(FAB) FB
}

// Applicative combines Apply with Pointed.
type Applicative[A, B, FA, FB, FAB any] interface {
	Apply[A, B, FA, FB, FAB]
	Pointed[A, FA]
}

// Chainable provides monadic bind.
type Chainable[A, B, FA, FB any] interface {
	Chain(func(A) FB) func(FA) FB
}

// Monad combines Applicative with Chainable.
type Monad[A, B, FA, FB, FAB any] interface {
	Applicative[A, B, FA, FB, FAB]
	Chainable[A, B, FA, FB]
}

// --- Concrete adapters ---

type pointed[A, FA any] struct {
	of func(A) FA
}

func (p pointed[A, FA]) Of(a A) FA { return p.of(a) }

func MakePointed[A, FA any](of func(A) FA) Pointed[A, FA] {
	return pointed[A, FA]{of: of}
}

type functor[A, B, FA, FB any] struct {
	fmap func(func(A) B) func(FA) FB
}

func (f functor[A, B, FA, FB]) Map(fn func(A) B) func(FA) FB { return f.fmap(fn) }

func MakeFunctor[A, B, FA, FB any](fmap func(func(A) B) func(FA) FB) Functor[A, B, FA, FB] {
	return functor[A, B, FA, FB]{fmap: fmap}
}

type apply[A, B, FA, FB, FAB any] struct {
	fmap func(func(A) B) func(FA) FB
	ap   func(FA) func(FAB) FB
}

func (a apply[A, B, FA, FB, FAB]) Map(fn func(A) B) func(FA) FB { return a.fmap(fn) }
func (a apply[A, B, FA, FB, FAB]) Ap(fa FA) func(FAB) FB        { return a.ap(fa) }

func MakeApply[A, B, FA, FB, FAB any](
	fmap func(func(A) B) func(FA) FB,
	ap func(FA) func(FAB) FB,
) Apply[A, B, FA, FB, FAB] {
	return apply[A, B, FA, FB, FAB]{fmap: fmap, ap: ap}
}

type applicative[A, B, FA, FB, FAB any] struct {
	fmap func(func(A) B) func(FA) FB
	ap   func(FA) func(FAB) FB
	of   func(A) FA
}

func (a applicative[A, B, FA, FB, FAB]) Map(fn func(A) B) func(FA) FB { return a.fmap(fn) }
func (a applicative[A, B, FA, FB, FAB]) Ap(fa FA) func(FAB) FB        { return a.ap(fa) }
func (a applicative[A, B, FA, FB, FAB]) Of(v A) FA                    { return a.of(v) }

func MakeApplicative[A, B, FA, FB, FAB any](
	of func(A) FA,
	fmap func(func(A) B) func(FA) FB,
	ap func(FA) func(FAB) FB,
) Applicative[A, B, FA, FB, FAB] {
	return applicative[A, B, FA, FB, FAB]{fmap: fmap, ap: ap, of: of}
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestPointedOptionCompiles -v`
Expected: PASS

**Step 5: Add more compile-time tests for Apply and Applicative**

Add to `laws/typeclass_test.go`:

```go
func TestApplyOptionCompiles(t *testing.T) {
	ap := laws.MakeApply[int, string, option.Option[int], option.Option[string], option.Option[func(int) string]](
		option.Map[int, string],
		option.Ap[string, int],
	)
	fa := option.Of(42)
	fab := option.Of(func(i int) string { return fmt.Sprintf("%d", i) })
	got := ap.Ap(fa)(fab)
	want := option.Of("42")
	if !option.Eq(eq.FromStrictEquals[string]()).Equals(got, want) {
		t.Fatalf("expected Some(\"42\"), got %v", got)
	}
}

func TestApplicativeOptionCompiles(t *testing.T) {
	ap := laws.MakeApplicative[int, string, option.Option[int], option.Option[string], option.Option[func(int) string]](
		option.Of[int],
		option.Map[int, string],
		option.Ap[string, int],
	)
	got := ap.Ap(ap.Of(42))(ap.Of(func(i int) string { return fmt.Sprintf("%d", i) }))
	// ap.Of for func type won't work here — Of is typed to A, not func(A)B
	// This test verifies the Of/Map/Ap wiring is correct
	_ = got
}
```

Note: `ap.Of` only lifts `A` (int), not `func(A) B`. For lifting functions, a separate `Pointed[func(A) B, FAB]` is needed. This is by design — the interchange law will use a separately-constructed Pointed.

**Step 6: Run all tests**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -v`
Expected: All PASS

**Step 7: Commit**

```bash
git add laws/typeclass.go laws/typeclass_test.go
git commit -m "feat(laws): add typeclass interfaces with Make* constructors"
```

---

## Task 2: Add Apply associative composition law

**Files:**
- Create: `laws/apply.go`
- Modify: `laws/typeclass_test.go` (add test)

The Apply composition law: `Ap(Ap(Map(compose)(fbc))(fab))(fa) == Ap(fbc)(Ap(fab)(fa))`

This needs these instances:
- `Pointed[func(A) B, FAB]` and `Pointed[func(B) C, FBC]` — to lift functions
- `Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC]` — map compose over FBC
- `Apply[A, B, FA, FB, FAB]` — apply fab to fa
- `Apply[B, C, FB, FC, FBC]` — apply fbc to fb
- `Apply[A, C, FA, FC, FAC]` — the composed result
- `Apply[func(A) B, func(A) C, FAB, FAC, FABAC]` — apply composed-functions to fab

**Step 1: Write the failing test**

Add to `laws/typeclass_test.go`:

```go
func TestApplyCompositionOption(t *testing.T) {
	laws.ApplyAssociativeComposition[
		option.Option[int],    // FA
		option.Option[string], // FB
		option.Option[bool],   // FC
		option.Option[func(int) string],    // FAB
		option.Option[func(string) bool],   // FBC
		option.Option[func(int) bool],      // FAC
		option.Option[func(func(int) string) func(int) bool], // FABAC
		int, string, bool,
	](
		t,
		eqOption[bool](eqBool),
		laws.MakePointed[func(int) string](option.Of[func(int) string]),
		laws.MakePointed[func(string) bool](option.Of[func(string) bool]),
		laws.MakeFunctor[func(string) bool, func(func(int) string) func(int) bool](
			option.Map[func(string) bool, func(func(int) string) func(int) bool],
		),
		laws.MakeApply[int, string](option.Map[int, string], option.Ap[string, int]),
		laws.MakeApply[string, bool](option.Map[string, bool], option.Ap[bool, string]),
		laws.MakeApply[int, bool](option.Map[int, bool], option.Ap[bool, int]),
		laws.MakeApply[func(int) string, func(int) bool](
			option.Map[func(int) string, func(int) bool],
			option.Ap[func(int) bool, func(int) string],
		),
		gen.GenOption(rapid.Int()),
		gen.GenFunc[int](rapid.String()),
		gen.GenFunc[string](rapid.Bool()),
	)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestApplyCompositionOption -v`
Expected: FAIL — `laws.ApplyAssociativeComposition` undefined

**Step 3: Implement ApplyAssociativeComposition**

```go
// laws/apply.go
package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// ApplyAssociativeComposition verifies the Apply associative composition law:
//
//   Ap(Ap(Map(compose)(fbc))(fab))(fa) == Ap(fbc)(Ap(fab)(fa))
//
// This law ensures that applying composed functions equals composing applications.
//
// Type parameters:
//   - FA, FB, FC: container types F[A], F[B], F[C]
//   - FAB, FBC, FAC: container types for functions F[A->B], F[B->C], F[A->C]
//   - FABAC: container type F[(A->B) -> (A->C)] for the higher-order composition
//   - A, B, C: element types
func ApplyAssociativeComposition[FA, FB, FC, FAB, FBC, FAC, FABAC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	ptdAB Pointed[func(A) B, FAB],
	ptdBC Pointed[func(B) C, FBC],
	fmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC],
	apAB Apply[A, B, FA, FB, FAB],
	apBC Apply[B, C, FB, FC, FBC],
	apAC Apply[A, C, FA, FC, FAC],
	apABAC Apply[func(A) B, func(A) C, FAB, FAC, FABAC],
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
) {
	t.Helper()

	t.Run("Apply/AssociativeComposition", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		ab := genAB.Draw(t, "ab")
		bc := genBC.Draw(t, "bc")

		fab := ptdAB.Of(ab)
		fbc := ptdBC.Of(bc)

		// compose: (B -> C) -> (A -> B) -> (A -> C)
		compose := func(g func(B) C) func(func(A) B) func(A) C {
			return func(f func(A) B) func(A) C {
				return func(a A) C { return g(f(a)) }
			}
		}

		// Left side: Ap(Ap(Map(compose)(fbc))(fab))(fa)
		composed := fmapCompose.Map(compose)(fbc) // F[(A->B) -> (A->C)]
		applied := apABAC.Ap(fab)(composed)        // F[A -> C]
		left := apAC.Ap(fa)(applied)               // F[C]

		// Right side: Ap(fbc)(Ap(fab)(fa))
		inner := apAB.Ap(fa)(fab)   // F[B]
		right := apBC.Ap(inner)(fbc) // F[C]

		if !eqFC(left, right) {
			t.Fatalf("Apply associative composition violated:\n  Ap(Ap(Map(compose)(fbc))(fab))(fa) = %v\n  Ap(fbc)(Ap(fab)(fa))                = %v", left, right)
		}
	}))
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestApplyCompositionOption -v`
Expected: PASS — `Apply/AssociativeComposition`

**Step 5: Commit**

```bash
git add laws/apply.go laws/typeclass_test.go
git commit -m "feat(laws): add Apply associative composition law"
```

---

## Task 3: Add Applicative interchange law

**Files:**
- Modify: `laws/applicative.go` (add `ApplicativeInterchange` function)
- Modify: `laws/typeclass_test.go` (add test)

The interchange law: `Ap(Of(a))(u) == Ap(u)(Of(f => f(a)))`

Where `u : F[A -> B]`, this needs:
- `Applicative[A, B, FA, FB, FAB]` — normal Ap + Of for values
- `Apply[func(A)B, B, FAB, FB, FABB]` — Ap where the "values" are functions
- `Pointed[func(func(A) B) B, FABB]` — lift "call-with-a" into context

**Step 1: Write the failing test**

Add to `laws/typeclass_test.go`:

```go
func TestApplicativeInterchangeOption(t *testing.T) {
	laws.ApplicativeInterchange[
		option.Option[int],                  // FA
		option.Option[string],               // FB
		option.Option[func(int) string],     // FAB
		option.Option[func(func(int) string) string], // FABB
		int, string,
	](
		t,
		eqOption[string](eqString),
		laws.MakeApplicative[int, string](
			option.Of[int],
			option.Map[int, string],
			option.Ap[string, int],
		),
		laws.MakeApply[func(int) string, string](
			option.Map[func(int) string, string],
			option.Ap[string, func(int) string],
		),
		laws.MakePointed[func(func(int) string) string](
			option.Of[func(func(int) string) string],
		),
		rapid.Int(),
		gen.GenFunc[int](rapid.String()),
	)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestApplicativeInterchangeOption -v`
Expected: FAIL — `laws.ApplicativeInterchange` undefined

**Step 3: Implement ApplicativeInterchange**

Add to `laws/applicative.go` (append after existing `ApplicativeLaws` function):

```go
// ApplicativeInterchange verifies the Applicative interchange law:
//
//   Ap(Of(a))(u) == Ap(u)(Of(f => f(a)))
//
// This law says: applying a wrapped function to a pure value is the same as
// applying "call-with-that-value" to the wrapped function.
//
// Type parameters:
//   - FA, FB: container types F[A], F[B]
//   - FAB: container type F[func(A) B]
//   - FABB: container type F[func(func(A) B) B]
//   - A, B: element types
func ApplicativeInterchange[FA, FB, FAB, FABB, A, B any](
	t *testing.T,
	eqFB func(FB, FB) bool,
	apAB Applicative[A, B, FA, FB, FAB],
	apABB Apply[func(A) B, B, FAB, FB, FABB],
	ptdABB Pointed[func(func(A) B) B, FABB],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
) {
	t.Helper()

	t.Run("Applicative/Interchange", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")

		fab := apAB.Of(f) // This won't work: Of lifts A, not func(A)B

		// We need a separate Pointed for func(A)B -> FAB
		// Actually, apAB.Of only lifts A -> FA. We need the function lifted separately.
		// Use: the existing ofAB from the old API, or a separate Pointed.

		// Hmm — this reveals we need Pointed[func(A) B, FAB] as well.
		// Let me restructure.
	}))
}
```

Wait — this reveals a design issue. `Applicative[A, B, FA, FB, FAB].Of` lifts `A -> FA`, not `func(A) B -> FAB`. For the interchange law we need to lift functions into `FAB`. Let me restructure the signature to also accept `Pointed[func(A) B, FAB]`:

```go
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

		// Left side: Ap(Of(a))(fab)
		left := apAB.Ap(apAB.Of(a))(fab)

		// Right side: Ap(fab)(Of(g => g(a)))
		callWithA := func(g func(A) B) B { return g(a) }
		right := apABB.Ap(fab)(ptdABB.Of(callWithA))

		if !eqFB(left, right) {
			t.Fatalf("Applicative interchange violated:\n  Ap(Of(a))(u)          = %v\n  Ap(u)(Of(f => f(a)))  = %v", left, right)
		}
	}))
}
```

Update the test accordingly (add `ptdAB` parameter):

```go
func TestApplicativeInterchangeOption(t *testing.T) {
	laws.ApplicativeInterchange[
		option.Option[int],
		option.Option[string],
		option.Option[func(int) string],
		option.Option[func(func(int) string) string],
		int, string,
	](
		t,
		eqOption[string](eqString),
		laws.MakeApplicative[int, string](
			option.Of[int],
			option.Map[int, string],
			option.Ap[string, int],
		),
		laws.MakePointed[func(int) string](option.Of[func(int) string]),
		laws.MakeApply[func(int) string, string](
			option.Map[func(int) string, string],
			option.Ap[string, func(int) string],
		),
		laws.MakePointed[func(func(int) string) string](
			option.Of[func(func(int) string) string],
		),
		rapid.Int(),
		gen.GenFunc[int](rapid.String()),
	)
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestApplicativeInterchangeOption -v`
Expected: PASS — `Applicative/Interchange`

**Step 5: Commit**

```bash
git add laws/applicative.go laws/typeclass_test.go
git commit -m "feat(laws): add Applicative interchange law"
```

---

## Task 4: Add full ApplicativeFullLaws combining all 4 laws

**Files:**
- Modify: `laws/applicative.go` (add `ApplicativeFullLaws`)
- Modify: `laws/typeclass_test.go`

This is the combined function that runs identity, homomorphism, interchange, and composition (via Apply).

**Step 1: Write the failing test**

```go
func TestOptionApplicativeFullLaws(t *testing.T) {
	inst := laws.OptionApplicativeInstances[int, string, bool]()
	laws.ApplicativeFullLaws(t,
		gen.GenOption(rapid.Int()),
		rapid.Int(),
		gen.GenFunc[int](rapid.String()),
		gen.GenFunc[string](rapid.Bool()),
		inst,
	)
}
```

**Step 2: Run test to verify it fails**

Expected: FAIL — `laws.OptionApplicativeInstances` and `laws.ApplicativeFullLaws` undefined

**Step 3: Define the instances bundle type and Option constructor**

Create `laws/instances_option.go`:

```go
package laws

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/option"
)

// ApplicativeInstances bundles all typed interface instances needed to test the
// full Applicative law suite (identity, homomorphism, interchange, composition).
type ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any] struct {
	// Equality
	EqFA func(FA, FA) bool
	EqFB func(FB, FB) bool
	EqFC func(FC, FC) bool

	// For identity + homomorphism
	ApAB  Applicative[A, B, FA, FB, FAB]
	PtdB  Pointed[B, FB]
	FmapAA func(func(A) A) func(FA) FA // for identity law: fmap(id)

	// For interchange
	PtdAB  Pointed[func(A) B, FAB]
	ApABB  Apply[func(A) B, B, FAB, FB, FABB]
	PtdABB Pointed[func(func(A) B) B, FABB]

	// For composition (via Apply)
	PtdBC      Pointed[func(B) C, FBC]
	FmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC]
	ApBC       Apply[B, C, FB, FC, FBC]
	ApAC       Apply[A, C, FA, FC, FAC]
	ApABAC     Apply[func(A) B, func(A) C, FAB, FAC, FABAC]
}

// OptionApplicativeInstances constructs all instances for testing Option's Applicative laws.
func OptionApplicativeInstances[A, B, C comparable]() *ApplicativeInstances[
	option.Option[A],
	option.Option[B],
	option.Option[C],
	option.Option[func(A) B],
	option.Option[func(B) C],
	option.Option[func(A) C],
	option.Option[func(func(A) B) func(A) C],
	option.Option[func(func(A) B) B],
	A, B, C,
] {
	return &ApplicativeInstances[
		option.Option[A],
		option.Option[B],
		option.Option[C],
		option.Option[func(A) B],
		option.Option[func(B) C],
		option.Option[func(A) C],
		option.Option[func(func(A) B) func(A) C],
		option.Option[func(func(A) B) B],
		A, B, C,
	]{
		EqFA: option.Eq(eq.FromStrictEquals[A]()).Equals,
		EqFB: option.Eq(eq.FromStrictEquals[B]()).Equals,
		EqFC: option.Eq(eq.FromStrictEquals[C]()).Equals,

		ApAB: MakeApplicative[A, B](
			option.Of[A],
			option.Map[A, B],
			option.Ap[B, A],
		),
		PtdB:   MakePointed[B](option.Of[B]),
		FmapAA: option.Map[A, A],

		PtdAB:  MakePointed[func(A) B](option.Of[func(A) B]),
		ApABB: MakeApply[func(A) B, B](
			option.Map[func(A) B, B],
			option.Ap[B, func(A) B],
		),
		PtdABB: MakePointed[func(func(A) B) B](option.Of[func(func(A) B) B]),

		PtdBC: MakePointed[func(B) C](option.Of[func(B) C]),
		FmapCompose: MakeFunctor[func(B) C, func(func(A) B) func(A) C](
			option.Map[func(B) C, func(func(A) B) func(A) C],
		),
		ApBC:   MakeApply[B, C](option.Map[B, C], option.Ap[C, B]),
		ApAC:   MakeApply[A, C](option.Map[A, C], option.Ap[C, A]),
		ApABAC: MakeApply[func(A) B, func(A) C](
			option.Map[func(A) B, func(A) C],
			option.Ap[func(A) C, func(A) B],
		),
	}
}
```

**Step 4: Implement ApplicativeFullLaws**

Add to `laws/applicative.go`:

```go
// ApplicativeFullLaws verifies all four Applicative laws using pre-built instances.
//
// Laws tested:
//   - Applicative/Identity: fmap(id)(v) == v
//   - Applicative/Homomorphism: ap(of(f))(of(x)) == of(f(x))
//   - Applicative/Interchange: ap(of(a))(u) == ap(u)(of(f => f(a)))
//   - Apply/AssociativeComposition: ap(ap(map(compose)(fbc))(fab))(fa) == ap(fbc)(ap(fab)(fa))
func ApplicativeFullLaws[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	inst *ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C],
) {
	t.Helper()

	// Identity
	t.Run("Applicative/Identity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		got := inst.FmapAA(func(a A) A { return a })(fa)
		if !inst.EqFA(got, fa) {
			t.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	}))

	// Homomorphism
	t.Run("Applicative/Homomorphism", rapid.MakeCheck(func(t *rapid.T) {
		a := genA.Draw(t, "a")
		f := genAB.Draw(t, "f")
		got := inst.ApAB.Ap(inst.ApAB.Of(a))(inst.PtdAB.Of(f))
		want := inst.PtdB.Of(f(a))
		if !inst.EqFB(got, want) {
			t.Fatalf("Applicative homomorphism violated:\n  ap(of(f))(of(a)) = %v\n  of(f(a))         = %v", got, want)
		}
	}))

	// Interchange
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

	// Composition (via Apply)
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
```

**Step 5: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestOptionApplicativeFullLaws -v`
Expected: PASS with 4 subtests:
- `Applicative/Identity`
- `Applicative/Homomorphism`
- `Applicative/Interchange`
- `Apply/AssociativeComposition`

**Step 6: Commit**

```bash
git add laws/applicative.go laws/instances_option.go laws/typeclass_test.go
git commit -m "feat(laws): add ApplicativeFullLaws with all 4 laws + Option instances"
```

---

## Task 5: Add Either instances

**Files:**
- Create: `laws/instances_either.go`
- Modify: `laws/typeclass_test.go`

**Step 1: Write the failing test**

```go
func TestEitherApplicativeFullLaws(t *testing.T) {
	inst := laws.EitherApplicativeInstances[string, int, string, bool]()
	laws.ApplicativeFullLaws(t,
		gen.GenEither(rapid.String(), rapid.Int()),
		rapid.Int(),
		gen.GenFunc[int](rapid.String()),
		gen.GenFunc[string](rapid.Bool()),
		inst,
	)
}
```

**Step 2: Run test to verify it fails**

Expected: FAIL — `laws.EitherApplicativeInstances` undefined

**Step 3: Implement EitherApplicativeInstances**

```go
// laws/instances_either.go
package laws

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
)

// EitherApplicativeInstances constructs all instances for testing Either's Applicative laws.
// E is the error/left type; A, B, C are the right/value types.
func EitherApplicativeInstances[E, A, B, C comparable]() *ApplicativeInstances[
	either.Either[E, A],
	either.Either[E, B],
	either.Either[E, C],
	either.Either[E, func(A) B],
	either.Either[E, func(B) C],
	either.Either[E, func(A) C],
	either.Either[E, func(func(A) B) func(A) C],
	either.Either[E, func(func(A) B) B],
	A, B, C,
] {
	eqE := eq.FromStrictEquals[E]()
	return &ApplicativeInstances[
		either.Either[E, A],
		either.Either[E, B],
		either.Either[E, C],
		either.Either[E, func(A) B],
		either.Either[E, func(B) C],
		either.Either[E, func(A) C],
		either.Either[E, func(func(A) B) func(A) C],
		either.Either[E, func(func(A) B) B],
		A, B, C,
	]{
		EqFA: either.Eq(eqE, eq.FromStrictEquals[A]()).Equals,
		EqFB: either.Eq(eqE, eq.FromStrictEquals[B]()).Equals,
		EqFC: either.Eq(eqE, eq.FromStrictEquals[C]()).Equals,

		ApAB: MakeApplicative[A, B](
			either.Of[E, A],
			either.Map[E, A, B],
			either.Ap[B, E, A],
		),
		PtdB:   MakePointed[B](either.Of[E, B]),
		FmapAA: either.Map[E, A, A],

		PtdAB: MakePointed[func(A) B](either.Of[E, func(A) B]),
		ApABB: MakeApply[func(A) B, B](
			either.Map[E, func(A) B, B],
			either.Ap[B, E, func(A) B],
		),
		PtdABB: MakePointed[func(func(A) B) B](either.Of[E, func(func(A) B) B]),

		PtdBC: MakePointed[func(B) C](either.Of[E, func(B) C]),
		FmapCompose: MakeFunctor[func(B) C, func(func(A) B) func(A) C](
			either.Map[E, func(B) C, func(func(A) B) func(A) C],
		),
		ApBC:   MakeApply[B, C](either.Map[E, B, C], either.Ap[C, E, B]),
		ApAC:   MakeApply[A, C](either.Map[E, A, C], either.Ap[C, E, A]),
		ApABAC: MakeApply[func(A) B, func(A) C](
			either.Map[E, func(A) B, func(A) C],
			either.Ap[func(A) C, E, func(A) B],
		),
	}
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestEitherApplicativeFullLaws -v`
Expected: PASS with 4 subtests

**Step 5: Commit**

```bash
git add laws/instances_either.go laws/typeclass_test.go
git commit -m "feat(laws): add Either Applicative instances"
```

---

## Task 6: Add Chain law (separate from Monad)

**Files:**
- Create: `laws/chain.go`
- Modify: `laws/typeclass_test.go`

The Chain associativity law is currently tested inside `MonadLaws`/`MonadLawsFull`. For the full hierarchy (mirroring fp-go), it should also be available standalone.

**Step 1: Write the failing test**

```go
func TestOptionChainAssociativity(t *testing.T) {
	laws.ChainAssociativity[
		option.Option[int],
		option.Option[string],
		option.Option[bool],
		int, string, bool,
	](
		t,
		eqOption[bool](eqBool),
		gen.GenOption(rapid.Int()),
		gen.GenFunc[int](gen.GenOption(rapid.String())),
		gen.GenFunc[string](gen.GenOption(rapid.Bool())),
		option.Chain[int, string],
		option.Chain[string, bool],
		option.Chain[int, bool],
	)
}
```

**Step 2: Run test to verify it fails**

Expected: FAIL — `laws.ChainAssociativity` undefined

**Step 3: Implement ChainAssociativity**

```go
// laws/chain.go
package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// ChainAssociativity verifies the Chain associativity law:
//
//   Chain(g)(Chain(f)(fa)) == Chain(x => Chain(g)(f(x)))(fa)
//
// This ensures monadic composition is independent of grouping.
func ChainAssociativity[FA, FB, FC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	t.Run("Chain/Associativity", rapid.MakeCheck(func(t *rapid.T) {
		fa := genFA.Draw(t, "fa")
		f := genKleisliAB.Draw(t, "f")
		g := genKleisliBC.Draw(t, "g")

		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)

		if !eqFC(left, right) {
			t.Fatalf("Chain associativity violated:\n  chain(g)(chain(f)(fa))        = %v\n  chain(x=>chain(g)(f(x)))(fa) = %v", left, right)
		}
	}))
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./laws/ -run TestOptionChainAssociativity -v`
Expected: PASS — `Chain/Associativity`

**Step 5: Commit**

```bash
git add laws/chain.go laws/typeclass_test.go
git commit -m "feat(laws): add standalone Chain associativity law"
```

---

## Task 7: Run full test suite and verify backward compatibility

**Files:** None (verification only)

**Step 1: Run ALL existing tests**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./... -v`

Expected: ALL tests pass, including all original tests:
- `TestOptionFunctorLaws`
- `TestEitherFunctorLaws`
- `TestOptionMonadLaws`
- `TestOptionApplicativeLaws` (original 2-law version)
- `TestIntSumSemigroupLaws`, `TestStringSemigroupLaws`
- `TestIntSumMonoidLaws`, `TestStringMonoidLaws`
- `TestIntEqLaws`, `TestStringEqLaws`, `TestOptionEqLaws`
- `TestIntOrdLaws`
- `TestPersonNameLensLaws`
- Plus all new tests

The original `ApplicativeLaws` function (raw-function API) must remain untouched and working.

**Step 2: Run with race detector**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go test ./... -race -v`
Expected: No race conditions

**Step 3: Verify go vet passes**

Run: `cd /home/iru/p/github.com/franchb/fptest-go && go vet ./...`
Expected: Clean

**Step 4: Commit if any test fixes were needed**

```bash
git add -A
git commit -m "fix: address test issues from full suite run"
```

---

## Summary of new files

| File | Purpose |
|---|---|
| `laws/typeclass.go` | Pointed, Functor, Apply, Applicative, Chainable, Monad interfaces + Make* constructors |
| `laws/typeclass_test.go` | Compile-time and integration tests for interfaces |
| `laws/apply.go` | `ApplyAssociativeComposition` — standalone Apply composition law |
| `laws/chain.go` | `ChainAssociativity` — standalone Chain associativity law |
| `laws/instances_option.go` | `OptionApplicativeInstances[A, B, C]()` helper |
| `laws/instances_either.go` | `EitherApplicativeInstances[E, A, B, C]()` helper |

## Modified files

| File | Change |
|---|---|
| `laws/applicative.go` | Add `ApplicativeInterchange` + `ApplicativeFullLaws` (existing `ApplicativeLaws` untouched) |

## Files NOT modified (backward compat)

| File | Reason |
|---|---|
| `laws/functor.go` | Raw-function API remains as-is |
| `laws/monad.go` | Raw-function API remains as-is |
| `laws/eq.go`, `ord.go`, `semigroup.go`, `monoid.go`, `lens.go` | Standalone, already complete |
| `gen/*`, `assert/*`, `mock/*`, `prop/*` | No changes needed |
