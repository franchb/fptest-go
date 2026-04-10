# hegel-go Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add hegel-go as an alternative PBT engine alongside rapid, with full interoperability across all fptest-go packages.

**Architecture:** Engine abstraction layer (`engine/` package) defines `T`, `Generator[A]`, and `Runner` interfaces. Core `laws/` and `prop/` packages gain optional `...Option` variadic params to accept a custom Runner (default: rapid). A separate `hegel/` sub-module provides hegel-native adapters and convenience wrappers. Monorepo with Go workspace.

**Tech Stack:** Go 1.25, `pgregory.net/rapid`, `hegel.dev/go/hegel`, `github.com/IBM/fp-go/v2`

**Spec:** `docs/superpowers/specs/2026-04-10-hegel-go-integration-design.md`

---

## File Map

### New files (core module)

| File | Responsibility |
|------|---------------|
| `engine/engine.go` | `T`, `Generator[A]`, `Runner` interfaces |
| `engine/rapid/rapid.go` | `RapidRunner`, `RapidGen[A]`, `Wrap[A]()` |
| `engine/rapid/rapid_test.go` | Tests for rapid adapter |
| `gen/adapt.go` | `ToEngine[A]()` — converts `Gen[A]` to `engine.Generator[A]` |
| `gen/adapt_test.go` | Tests for Gen adapter |
| `laws/config.go` | `Option`, `config`, `WithRunner()`, `resolveConfig()` |
| `laws/semigroup_engine.go` | `SemigroupLawsEngine[A]()` — engine-generic |
| `laws/monoid_engine.go` | `MonoidLawsEngine[A]()` — engine-generic |
| `laws/eq_engine.go` | `EqLawsEngine[A]()` — engine-generic |
| `laws/ord_engine.go` | `OrdLawsEngine[A]()` — engine-generic |
| `laws/lens_engine.go` | `LensLawsEngine[S,A]()` — engine-generic |
| `laws/functor_engine.go` | `FunctorLawsEngine[...]()` — engine-generic |
| `laws/chain_engine.go` | `ChainAssociativityEngine[...]()` — engine-generic |
| `laws/monad_engine.go` | `MonadLawsEngine[...]()`, `MonadLawsFullEngine[...]()` — engine-generic |
| `laws/apply_engine.go` | `ApplyAssociativeCompositionEngine[...]()` — engine-generic |
| `laws/applicative_engine.go` | `ApplicativeLawsEngine[...]()`, `ApplicativeFullLawsEngine[...]()`, `ApplicativeInterchangeEngine[...]()` — engine-generic |
| `prop/config.go` | `Option`, `config`, `WithRunner()`, `resolveConfig()` (prop's own copy) |
| `prop/roundtrip_engine.go` | `RoundTripEngine[A,B]()` etc — engine-generic |
| `prop/oracle_engine.go` | `OracleEngine[A,B]()`, `IdempotentEngine[A]()`, `CommutativeEngine[A,B]()`, `InvariantEngine[A]()` — engine-generic |
| `go.work` | Go workspace linking `.` and `./hegel` |

### Modified files (core module)

| File | Change |
|------|--------|
| `laws/semigroup.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/monoid.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/eq.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/ord.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/lens.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/functor.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/chain.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/monad.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/apply.go` | Add `opts ...Option`, delegate to engine variant |
| `laws/applicative.go` | Add `opts ...Option`, delegate to engine variant |
| `prop/roundtrip.go` | Add `opts ...Option`, delegate to engine variant |
| `prop/oracle.go` | Add `opts ...Option`, delegate to engine variant |

### New files (hegel sub-module)

| File | Responsibility |
|------|---------------|
| `hegel/go.mod` | Module `github.com/franchb/fptest/hegel` |
| `hegel/runner.go` | `HegelRunner` implements `engine.Runner` via `hegel.Case` |
| `hegel/runner_test.go` | Tests for HegelRunner |
| `hegel/gen.go` | `HegelGen[A]` wraps `hegel.Generator[T]` as `engine.Generator[A]` |
| `hegel/gen_test.go` | Tests for HegelGen |
| `hegel/hegelgen/domain.go` | `Emails()`, `URLs()`, `Dates()`, etc. |
| `hegel/hegelgen/fp.go` | `Option[A]()`, `Either[E,A]()`, `Right[E,A]()`, etc. |
| `hegel/hegelgen/domain_test.go` | Tests for domain generators |
| `hegel/hegelgen/fp_test.go` | Tests for FP generators |
| `hegel/laws/laws.go` | Hegel convenience wrappers for all law functions |
| `hegel/laws/laws_test.go` | Tests: same laws run with hegel engine |
| `hegel/prop/prop.go` | Hegel convenience wrappers for all prop functions |
| `hegel/prop/prop_test.go` | Tests for hegel prop wrappers |

---

### Task 1: Engine Abstraction Interfaces

**Files:**
- Create: `engine/engine.go`

- [ ] **Step 1: Create engine/engine.go with T, Generator, Runner interfaces**

```go
// Package engine defines the PBT engine abstraction for fptest-go.
//
// This package contains only interfaces — no external dependencies.
// Concrete implementations live in engine/rapid/ (for pgregory.net/rapid)
// and hegel/ (for hegel.dev/go/hegel).
package engine

import "testing"

// T is the test context provided by a PBT engine during property execution.
// Both rapid.T and hegel.T satisfy this interface.
type T interface {
	testing.TB
}

// Generator draws values of type A from the PBT engine's search space.
type Generator[A any] interface {
	Draw(t T, label string) A
}

// Runner executes a property check as a named subtest.
type Runner interface {
	MakeCheck(t *testing.T, name string, prop func(T))
}
```

- [ ] **Step 2: Verify the package compiles**

Run: `go build ./engine/...`
Expected: success (no output)

- [ ] **Step 3: Commit**

```bash
git add engine/engine.go
git commit -m "feat(engine): add PBT engine abstraction interfaces"
```

---

### Task 2: Rapid Adapter

**Files:**
- Create: `engine/rapid/rapid.go`
- Create: `engine/rapid/rapid_test.go`

- [ ] **Step 1: Write the failing test for RapidRunner and RapidGen**

```go
package rapid_test

import (
	"testing"

	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
	rapidlib "pgregory.net/rapid"
)

func TestRapidRunnerExecutesProperty(t *testing.T) {
	var runner engine.Runner = enginerapid.RapidRunner{}
	executed := false
	runner.MakeCheck(t, "test/property", func(et engine.T) {
		executed = true
	})
	if !executed {
		t.Fatal("property was never executed")
	}
}

func TestRapidGenDrawsValues(t *testing.T) {
	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/draw", func(et engine.T) {
		gen := enginerapid.Wrap(rapidlib.IntRange(1, 100))
		val := gen.Draw(et, "n")
		if val < 1 || val > 100 {
			t.Fatalf("expected value in [1,100], got %d", val)
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./engine/rapid/ -v -run TestRapid`
Expected: FAIL — package or types not found

- [ ] **Step 3: Write the implementation**

```go
// Package rapid provides the rapid PBT engine adapter for fptest-go.
package rapid

import (
	"testing"

	"github.com/franchb/fptest/engine"
	rapidlib "pgregory.net/rapid"
)

// RapidRunner implements engine.Runner using rapid.MakeCheck.
type RapidRunner struct{}

// MakeCheck runs the property as a named subtest using rapid's property checker.
func (RapidRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
	t.Helper()
	t.Run(name, rapidlib.MakeCheck(func(rt *rapidlib.T) {
		prop(rt)
	}))
}

// RapidGen wraps a *rapid.Generator[A] as an engine.Generator[A].
type RapidGen[A any] struct {
	G *rapidlib.Generator[A]
}

// Draw produces a value from the wrapped rapid generator.
func (rg RapidGen[A]) Draw(t engine.T, label string) A {
	return rg.G.Draw(t.(*rapidlib.T), label)
}

// Wrap converts a *rapid.Generator[A] into an engine.Generator[A].
func Wrap[A any](g *rapidlib.Generator[A]) engine.Generator[A] {
	return RapidGen[A]{G: g}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./engine/rapid/ -v -run TestRapid`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add engine/rapid/rapid.go engine/rapid/rapid_test.go
git commit -m "feat(engine/rapid): add rapid PBT engine adapter"
```

---

### Task 3: Gen[A] to engine.Generator Adapter

**Files:**
- Create: `gen/adapt.go`
- Create: `gen/adapt_test.go`

- [ ] **Step 1: Write the failing test**

```go
package gen_test

import (
	"testing"

	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
	"github.com/franchb/fptest/gen"
	"pgregory.net/rapid"
)

func TestToEngineDrawsValues(t *testing.T) {
	g := gen.Of(42)
	var eg engine.Generator[int] = gen.ToEngine(g)

	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/ToEngine", func(et engine.T) {
		val := eg.Draw(et, "n")
		if val != 42 {
			t.Fatalf("expected 42, got %d", val)
		}
	})
}

func TestToEngineWithRapidGen(t *testing.T) {
	g := gen.FromRapid(rapid.IntRange(10, 20))
	eg := gen.ToEngine(g)

	runner := enginerapid.RapidRunner{}
	runner.MakeCheck(t, "test/ToEngineRapid", func(et engine.T) {
		val := eg.Draw(et, "n")
		if val < 10 || val > 20 {
			t.Fatalf("expected value in [10,20], got %d", val)
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./gen/ -v -run TestToEngine`
Expected: FAIL — `gen.ToEngine` not found

- [ ] **Step 3: Write the implementation**

```go
package gen

import (
	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
)

// ToEngine converts a Gen[A] to an engine.Generator[A] via the rapid adapter.
// This allows Gen values to be used with any PBT engine through the engine abstraction.
func ToEngine[A any](g Gen[A]) engine.Generator[A] {
	return enginerapid.Wrap(ToRapid(g))
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./gen/ -v -run TestToEngine`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add gen/adapt.go gen/adapt_test.go
git commit -m "feat(gen): add ToEngine adapter for Gen[A] to engine.Generator"
```

---

### Task 4: Laws Config System

**Files:**
- Create: `laws/config.go`

- [ ] **Step 1: Create the config and Option types**

```go
package laws

import (
	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
)

// Option configures law verification behavior.
type Option func(*config)

type config struct {
	runner engine.Runner
}

// WithRunner sets the PBT engine runner for law verification.
// If not set, the default rapid runner is used.
func WithRunner(r engine.Runner) Option {
	return func(c *config) { c.runner = r }
}

func resolveConfig(opts []Option) config {
	cfg := config{runner: enginerapid.RapidRunner{}}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
```

- [ ] **Step 2: Verify the package compiles**

Run: `go build ./laws/...`
Expected: success

- [ ] **Step 3: Commit**

```bash
git add laws/config.go
git commit -m "feat(laws): add config system for pluggable PBT engine runner"
```

---

### Task 5: Refactor Semigroup + Monoid Laws (Engine-Generic)

**Files:**
- Create: `laws/semigroup_engine.go`
- Create: `laws/monoid_engine.go`
- Modify: `laws/semigroup.go`
- Modify: `laws/monoid.go`

- [ ] **Step 1: Verify existing tests pass before refactoring**

Run: `go test ./laws/ -v -run "TestIntSumSemigroup|TestStringSemigroup|TestIntSumMonoid|TestStringMonoid"`
Expected: PASS (4 tests)

- [ ] **Step 2: Create semigroup_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// SemigroupLawsEngine verifies the Semigroup law (associativity) using the given engine.
func SemigroupLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
) {
	t.Helper()

	runner.MakeCheck(t, "Semigroup/Associativity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")

		left := concat(concat(a, b), c)
		right := concat(a, concat(b, c))
		if !eqA(left, right) {
			t.Fatalf("Semigroup associativity violated:\n  (a <> b) <> c = %v\n  a <> (b <> c) = %v", left, right)
		}
	})
}
```

- [ ] **Step 3: Update semigroup.go to delegate**

Replace the body of `SemigroupLaws` with delegation to the engine variant. Add `opts ...Option` parameter:

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// SemigroupLaws verifies the Semigroup law (associativity).
func SemigroupLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	SemigroupLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), eqA, concat)
}

// SemigroupInterfaceLaws verifies Semigroup laws using fp-go's Semigroup interface.
func SemigroupInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	sg interface{ Concat(A, A) A },
	opts ...Option,
) {
	t.Helper()
	SemigroupLaws(t, genA, eqA, sg.Concat, opts...)
}
```

- [ ] **Step 4: Create monoid_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// MonoidLawsEngine verifies the Monoid laws using the given engine.
func MonoidLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
) {
	t.Helper()

	SemigroupLawsEngine(t, runner, genA, eqA, concat)

	runner.MakeCheck(t, "Monoid/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		got := concat(empty, a)
		if !eqA(got, a) {
			t.Fatalf("Monoid left identity violated:\n  empty <> a = %v\n  a          = %v", got, a)
		}
	})

	runner.MakeCheck(t, "Monoid/RightIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		got := concat(a, empty)
		if !eqA(got, a) {
			t.Fatalf("Monoid right identity violated:\n  a <> empty = %v\n  a          = %v", got, a)
		}
	})
}
```

- [ ] **Step 5: Update monoid.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// MonoidLaws verifies the Monoid laws (left identity, right identity, associativity).
func MonoidLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonoidLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), eqA, concat, empty)
}

// MonoidInterfaceLaws verifies Monoid laws using fp-go's Monoid interface.
func MonoidInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	m interface {
		Concat(A, A) A
		Empty() A
	},
	opts ...Option,
) {
	t.Helper()
	MonoidLaws(t, genA, eqA, m.Concat, m.Empty(), opts...)
}
```

- [ ] **Step 6: Run existing tests to verify backwards compatibility**

Run: `go test ./laws/ -v -run "TestIntSumSemigroup|TestStringSemigroup|TestIntSumMonoid|TestStringMonoid"`
Expected: PASS (4 tests, identical behavior)

- [ ] **Step 7: Commit**

```bash
git add laws/semigroup_engine.go laws/monoid_engine.go laws/semigroup.go laws/monoid.go
git commit -m "feat(laws): add engine-generic Semigroup and Monoid law verification"
```

---

### Task 6: Refactor Eq + Ord + Lens Laws (Engine-Generic)

**Files:**
- Create: `laws/eq_engine.go`
- Create: `laws/ord_engine.go`
- Create: `laws/lens_engine.go`
- Modify: `laws/eq.go`
- Modify: `laws/ord.go`
- Modify: `laws/lens.go`

- [ ] **Step 1: Verify existing tests pass**

Run: `go test ./laws/ -v -run "TestIntEq|TestStringEq|TestOptionEq|TestIntOrd|TestPersonName"`
Expected: PASS

- [ ] **Step 2: Create eq_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// EqLawsEngine verifies the Eq laws using the given engine.
func EqLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	equals func(A, A) bool,
) {
	t.Helper()

	runner.MakeCheck(t, "Eq/Reflexivity", func(et engine.T) {
		a := genA.Draw(et, "a")
		if !equals(a, a) {
			t.Fatalf("Eq reflexivity violated: equals(%v, %v) = false", a, a)
		}
	})

	runner.MakeCheck(t, "Eq/Symmetry", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		ab := equals(a, b)
		ba := equals(b, a)
		if ab != ba {
			t.Fatalf("Eq symmetry violated:\n  equals(%v, %v) = %v\n  equals(%v, %v) = %v", a, b, ab, b, a, ba)
		}
	})

	runner.MakeCheck(t, "Eq/Transitivity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")
		if equals(a, b) && equals(b, c) && !equals(a, c) {
			t.Fatalf("Eq transitivity violated:\n  equals(%v, %v) = true\n  equals(%v, %v) = true\n  equals(%v, %v) = false", a, b, b, c, a, c)
		}
	})
}
```

- [ ] **Step 3: Update eq.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// EqLaws verifies the Eq laws (reflexivity, symmetry, transitivity).
func EqLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	EqLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), equals)
}

// EqInterfaceLaws verifies Eq laws using fp-go's Eq interface.
func EqInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	eqA interface{ Equals(A, A) bool },
	opts ...Option,
) {
	t.Helper()
	EqLaws(t, genA, eqA.Equals, opts...)
}
```

- [ ] **Step 4: Create ord_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// OrdLawsEngine verifies the Ord laws using the given engine.
func OrdLawsEngine[A any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
) {
	t.Helper()

	EqLawsEngine(t, runner, genA, equals)

	runner.MakeCheck(t, "Ord/Antisymmetry", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		if compare(a, b) <= 0 && compare(b, a) <= 0 && !equals(a, b) {
			t.Fatalf("Ord antisymmetry violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  but not equal", a, b, compare(a, b), b, a, compare(b, a))
		}
	})

	runner.MakeCheck(t, "Ord/Transitivity", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		c := genA.Draw(et, "c")
		if compare(a, b) <= 0 && compare(b, c) <= 0 && compare(a, c) > 0 {
			t.Fatalf("Ord transitivity violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, c, compare(b, c), a, c, compare(a, c))
		}
	})

	runner.MakeCheck(t, "Ord/Totality", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		if compare(a, b) > 0 && compare(b, a) > 0 {
			t.Fatalf("Ord totality violated:\n  compare(%v, %v) = %d\n  compare(%v, %v) = %d",
				a, b, compare(a, b), b, a, compare(b, a))
		}
	})

	runner.MakeCheck(t, "Ord/Consistency", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		cmp := compare(a, b)
		eq := equals(a, b)
		if (cmp == 0) != eq {
			t.Fatalf("Ord consistency violated:\n  compare(%v, %v) = %d\n  equals(%v, %v) = %v",
				a, b, cmp, a, b, eq)
		}
	})
}
```

- [ ] **Step 5: Update ord.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// OrdLaws verifies the Ord laws (antisymmetry, transitivity, totality).
func OrdLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	OrdLawsEngine(t, cfg.runner, enginerapid.Wrap(genA), equals, compare)
}

// OrdInterfaceLaws verifies Ord laws using fp-go's Ord interface.
func OrdInterfaceLaws[A any](
	t *testing.T,
	genA *rapid.Generator[A],
	ordA interface {
		Equals(A, A) bool
		Compare(A, A) int
	},
	opts ...Option,
) {
	t.Helper()
	OrdLaws(t, genA, ordA.Equals, ordA.Compare, opts...)
}
```

- [ ] **Step 6: Create lens_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// LensLawsEngine verifies the Lens laws using the given engine.
func LensLawsEngine[S, A any](
	t *testing.T,
	runner engine.Runner,
	genS engine.Generator[S],
	genA engine.Generator[A],
	eqS func(S, S) bool,
	eqA func(A, A) bool,
	get func(S) A,
	set func(A) func(S) S,
) {
	t.Helper()

	runner.MakeCheck(t, "Lens/GetSet", func(et engine.T) {
		s := genS.Draw(et, "s")
		got := set(get(s))(s)
		if !eqS(got, s) {
			t.Fatalf("Lens GetSet violated:\n  set(get(s))(s) = %v\n  s              = %v", got, s)
		}
	})

	runner.MakeCheck(t, "Lens/SetGet", func(et engine.T) {
		s := genS.Draw(et, "s")
		a := genA.Draw(et, "a")
		got := get(set(a)(s))
		if !eqA(got, a) {
			t.Fatalf("Lens SetGet violated:\n  get(set(a)(s)) = %v\n  a              = %v", got, a)
		}
	})

	runner.MakeCheck(t, "Lens/SetSet", func(et engine.T) {
		s := genS.Draw(et, "s")
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		left := set(b)(set(a)(s))
		right := set(b)(s)
		if !eqS(left, right) {
			t.Fatalf("Lens SetSet violated:\n  set(b)(set(a)(s)) = %v\n  set(b)(s)         = %v", left, right)
		}
	})
}
```

- [ ] **Step 7: Update lens.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// LensLaws verifies the Lens laws (get-set, set-get, set-set).
func LensLaws[S, A any](
	t *testing.T,
	genS *rapid.Generator[S],
	genA *rapid.Generator[A],
	eqS func(S, S) bool,
	eqA func(A, A) bool,
	get func(S) A,
	set func(A) func(S) S,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	LensLawsEngine(t, cfg.runner, enginerapid.Wrap(genS), enginerapid.Wrap(genA), eqS, eqA, get, set)
}
```

- [ ] **Step 8: Run all existing tests to verify backwards compatibility**

Run: `go test ./laws/ -v`
Expected: all existing tests PASS

- [ ] **Step 9: Commit**

```bash
git add laws/eq_engine.go laws/ord_engine.go laws/lens_engine.go laws/eq.go laws/ord.go laws/lens.go
git commit -m "feat(laws): add engine-generic Eq, Ord, and Lens law verification"
```

---

### Task 7: Refactor Functor + Chain Laws (Engine-Generic)

**Files:**
- Create: `laws/functor_engine.go`
- Create: `laws/chain_engine.go`
- Modify: `laws/functor.go`
- Modify: `laws/chain.go`

- [ ] **Step 1: Verify existing tests pass**

Run: `go test ./laws/ -v -run "TestOptionFunctor|TestEitherFunctor"`
Expected: PASS

- [ ] **Step 2: Create functor_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// FunctorLawsEngine verifies the Functor laws using the given engine.
func FunctorLawsEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genFA engine.Generator[FA],
	genAB engine.Generator[func(A) B],
	genBC engine.Generator[func(B) C],
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	fmapAA func(func(A) A) func(FA) FA,
	fmapAB func(func(A) B) func(FA) FB,
	fmapBC func(func(B) C) func(FB) FC,
	fmapAC func(func(A) C) func(FA) FC,
	identity func(A) A,
	compose func(func(A) B, func(B) C) func(A) C,
) {
	t.Helper()

	runner.MakeCheck(t, "Functor/Identity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		got := fmapAA(identity)(fa)
		if !eqFA(got, fa) {
			t.Fatalf("Functor identity law violated:\n  fa       = %v\n  fmap(id) = %v", fa, got)
		}
	})

	runner.MakeCheck(t, "Functor/Composition", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genAB.Draw(et, "f")
		g := genBC.Draw(et, "g")
		composed := fmapAC(compose(f, g))(fa)
		chained := fmapBC(g)(fmapAB(f)(fa))
		if !eqFC(composed, chained) {
			t.Fatalf("Functor composition law violated:\n  fmap(g.f)(fa) = %v\n  fmap(g)(fmap(f)(fa)) = %v", composed, chained)
		}
	})
}
```

- [ ] **Step 3: Update functor.go to delegate**

```go
// Package laws provides typeclass law verification using property-based testing.
//
// Each law test function runs as subtests via rapid.MakeCheck, producing output
// like TestOptionLaws/Functor/Identity that integrates with go test -run.
// Functions accept typeclass operations as parameters (not interfaces), making them
// work with any type that provides the right operations — including fp-go's Option,
// Either, IO, and user-defined types.
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// FunctorLaws verifies the Functor laws (identity and composition) for a type constructor.
//
// Type parameters:
//   - FA: the container type F[A]
//   - FB: the container type F[B]
//   - FC: the container type F[C]
//   - A, B, C: element types
//
// The fmap parameters correspond to the Map operation specialized for each type combination.
// In fp-go terms: fmapAB = option.Map[A, B], etc.
func FunctorLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	fmapAA func(func(A) A) func(FA) FA,
	fmapAB func(func(A) B) func(FA) FB,
	fmapBC func(func(B) C) func(FB) FC,
	fmapAC func(func(A) C) func(FA) FC,
	identity func(A) A,
	compose func(func(A) B, func(B) C) func(A) C,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	FunctorLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genAB), enginerapid.Wrap(genBC),
		eqFA, eqFC, fmapAA, fmapAB, fmapBC, fmapAC, identity, compose)
}
```

- [ ] **Step 4: Create chain_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// ChainAssociativityEngine verifies the Chain associativity law using the given engine.
func ChainAssociativityEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	eqFC func(FC, FC) bool,
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Chain/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")
		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			t.Fatalf("Chain associativity violated:\n  chain(g)(chain(f)(fa))        = %v\n  chain(x=>chain(g)(f(x)))(fa) = %v", left, right)
		}
	})
}
```

- [ ] **Step 5: Update chain.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// ChainAssociativity verifies the Chain associativity law.
func ChainAssociativity[FA, FB, FC, A, B, C any](
	t *testing.T,
	eqFC func(FC, FC) bool,
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ChainAssociativityEngine(t, cfg.runner, eqFC,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		chainAB, chainBC, chainAC)
}
```

- [ ] **Step 6: Run existing tests**

Run: `go test ./laws/ -v -run "TestOptionFunctor|TestEitherFunctor"`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add laws/functor_engine.go laws/chain_engine.go laws/functor.go laws/chain.go
git commit -m "feat(laws): add engine-generic Functor and Chain law verification"
```

---

### Task 8: Refactor Monad Laws (Engine-Generic)

**Files:**
- Create: `laws/monad_engine.go`
- Modify: `laws/monad.go`

- [ ] **Step 1: Verify existing tests pass**

Run: `go test ./laws/ -v -run TestOptionMonad`
Expected: PASS

- [ ] **Step 2: Create monad_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// MonadLawsEngine verifies the Monad laws using the given engine.
func MonadLawsEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Monad/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genKleisliAB.Draw(et, "f")
		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			t.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	})

	runner.MakeCheck(t, "Monad/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")
		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			t.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	})
}

// MonadLawsFullEngine verifies all Monad laws including right identity.
func MonadLawsFullEngine[FA, FB, FC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	genA engine.Generator[A],
	genFA engine.Generator[FA],
	genKleisliAB engine.Generator[func(A) FB],
	genKleisliBC engine.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAA func(func(A) FA) func(FA) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()

	runner.MakeCheck(t, "Monad/LeftIdentity", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genKleisliAB.Draw(et, "f")
		got := chainAB(f)(of(a))
		want := f(a)
		if !eqFB(got, want) {
			t.Fatalf("Monad left identity violated:\n  chain(f)(of(a)) = %v\n  f(a)            = %v", got, want)
		}
	})

	runner.MakeCheck(t, "Monad/RightIdentity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		got := chainAA(of)(fa)
		if !eqFA(got, fa) {
			t.Fatalf("Monad right identity violated:\n  chain(of)(fa) = %v\n  fa            = %v", got, fa)
		}
	})

	runner.MakeCheck(t, "Monad/Associativity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		f := genKleisliAB.Draw(et, "f")
		g := genKleisliBC.Draw(et, "g")
		left := chainBC(g)(chainAB(f)(fa))
		right := chainAC(func(a A) FC {
			return chainBC(g)(f(a))
		})(fa)
		if !eqFC(left, right) {
			t.Fatalf("Monad associativity violated:\n  chain(g)(chain(f)(fa)) = %v\n  chain(g.f)(fa)         = %v", left, right)
		}
	})
}
```

- [ ] **Step 3: Update monad.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// MonadLaws verifies the Monad laws (left identity, associativity).
func MonadLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonadLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA),
		enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAB, chainBC, chainAC)
}

// MonadLawsFull verifies all Monad laws including right identity.
func MonadLawsFull[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA *rapid.Generator[A],
	genFA *rapid.Generator[FA],
	genKleisliAB *rapid.Generator[func(A) FB],
	genKleisliBC *rapid.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAA func(func(A) FA) func(FA) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	MonadLawsFullEngine(t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA),
		enginerapid.Wrap(genKleisliAB), enginerapid.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAA, chainAB, chainBC, chainAC)
}
```

- [ ] **Step 4: Run existing tests**

Run: `go test ./laws/ -v -run TestOptionMonad`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add laws/monad_engine.go laws/monad.go
git commit -m "feat(laws): add engine-generic Monad law verification"
```

---

### Task 9: Refactor Apply + Applicative Laws (Engine-Generic)

**Files:**
- Create: `laws/apply_engine.go`
- Create: `laws/applicative_engine.go`
- Modify: `laws/apply.go`
- Modify: `laws/applicative.go`

- [ ] **Step 1: Verify existing tests pass**

Run: `go test ./laws/ -v -run TestOptionApplicative`
Expected: PASS

- [ ] **Step 2: Create apply_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// ApplyAssociativeCompositionEngine verifies the Apply associative composition law using the given engine.
func ApplyAssociativeCompositionEngine[FA, FB, FC, FAB, FBC, FAC, FABAC, A, B, C any](
	t *testing.T,
	runner engine.Runner,
	eqFC func(FC, FC) bool,
	ptdAB Pointed[func(A) B, FAB],
	ptdBC Pointed[func(B) C, FBC],
	fmapCompose Functor[func(B) C, func(func(A) B) func(A) C, FBC, FABAC],
	apAB Apply[A, B, FA, FB, FAB],
	apBC Apply[B, C, FB, FC, FBC],
	apAC Apply[A, C, FA, FC, FAC],
	apABAC Apply[func(A) B, func(A) C, FAB, FAC, FABAC],
	genFA engine.Generator[FA],
	genAB engine.Generator[func(A) B],
	genBC engine.Generator[func(B) C],
) {
	t.Helper()

	runner.MakeCheck(t, "Apply/AssociativeComposition", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		ab := genAB.Draw(et, "ab")
		bc := genBC.Draw(et, "bc")

		fab := ptdAB.Of(ab)
		fbc := ptdBC.Of(bc)

		compose := func(g func(B) C) func(func(A) B) func(A) C {
			return func(f func(A) B) func(A) C {
				return func(a A) C { return g(f(a)) }
			}
		}

		composed := fmapCompose.Map(compose)(fbc)
		applied := apABAC.Ap(fab)(composed)
		left := apAC.Ap(fa)(applied)

		inner := apAB.Ap(fa)(fab)
		right := apBC.Ap(inner)(fbc)

		if !eqFC(left, right) {
			t.Fatalf("Apply associative composition violated:\n  Ap(Ap(Map(compose)(fbc))(fab))(fa) = %v\n  Ap(fbc)(Ap(fab)(fa))               = %v", left, right)
		}
	})
}
```

- [ ] **Step 3: Update apply.go to delegate**

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// ApplyAssociativeComposition verifies the Apply associative composition law.
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplyAssociativeCompositionEngine(t, cfg.runner, eqFC,
		ptdAB, ptdBC, fmapCompose, apAB, apBC, apAC, apABAC,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genAB), enginerapid.Wrap(genBC))
}
```

- [ ] **Step 4: Create applicative_engine.go**

```go
package laws

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// ApplicativeLawsEngine verifies the Applicative laws using the given engine.
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
		got := fmapAA(identity)(fa)
		if !eqFA(got, fa) {
			t.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	})

	runner.MakeCheck(t, "Applicative/Homomorphism", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")
		got := apAB(ofA(a))(ofAB(f))
		want := ofB(f(a))
		if !eqFB(got, want) {
			t.Fatalf("Applicative homomorphism violated:\n  ap(pure(f), pure(x)) = %v\n  pure(f(x))           = %v", got, want)
		}
	})
}

// ApplicativeInterchangeEngine verifies the Applicative interchange law using the given engine.
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
		fab := ptdAB.Of(f)
		left := apAB.Ap(apAB.Of(a))(fab)
		callWithA := func(g func(A) B) B { return g(a) }
		right := apABB.Ap(fab)(ptdABB.Of(callWithA))
		if !eqFB(left, right) {
			t.Fatalf("Applicative interchange violated:\n  Ap(Of(a))(u)         = %v\n  Ap(u)(Of(f=>f(a)))   = %v", left, right)
		}
	})
}

// ApplicativeFullLawsEngine verifies all four Applicative laws using the given engine.
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

	runner.MakeCheck(t, "Applicative/Identity", func(et engine.T) {
		fa := genFA.Draw(et, "fa")
		got := inst.FmapAA(func(a A) A { return a })(fa)
		if !inst.EqFA(got, fa) {
			t.Fatalf("Applicative identity violated:\n  fmap(id)(v) = %v\n  v           = %v", got, fa)
		}
	})

	runner.MakeCheck(t, "Applicative/Homomorphism", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")
		got := inst.ApAB.Ap(inst.ApAB.Of(a))(inst.PtdAB.Of(f))
		want := inst.PtdB.Of(f(a))
		if !inst.EqFB(got, want) {
			t.Fatalf("Applicative homomorphism violated:\n  ap(of(f))(of(a)) = %v\n  of(f(a))         = %v", got, want)
		}
	})

	runner.MakeCheck(t, "Applicative/Interchange", func(et engine.T) {
		a := genA.Draw(et, "a")
		f := genAB.Draw(et, "f")
		fab := inst.PtdAB.Of(f)
		left := inst.ApAB.Ap(inst.ApAB.Of(a))(fab)
		callWithA := func(g func(A) B) B { return g(a) }
		right := inst.ApABB.Ap(fab)(inst.PtdABB.Of(callWithA))
		if !inst.EqFB(left, right) {
			t.Fatalf("Applicative interchange violated:\n  ap(of(a))(u)         = %v\n  ap(u)(of(f=>f(a)))   = %v", left, right)
		}
	})

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
			t.Fatalf("Apply associative composition violated:\n  ap(ap(map(compose)(fbc))(fab))(fa) = %v\n  ap(fbc)(ap(fab)(fa))               = %v", left, right)
		}
	})
}
```

- [ ] **Step 5: Update applicative.go to delegate**

Replace the contents of `applicative.go` with delegation. The function signatures stay the same but gain `opts ...Option`:

- `ApplicativeLaws` delegates to `ApplicativeLawsEngine`
- `ApplicativeInterchange` delegates to `ApplicativeInterchangeEngine`
- `ApplicativeFullLaws` delegates to `ApplicativeFullLawsEngine`

Each wraps its `*rapid.Generator` params via `enginerapid.Wrap(...)`.

```go
package laws

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// ApplicativeLaws verifies the Applicative functor laws.
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
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genA), enginerapid.Wrap(genFA), enginerapid.Wrap(genAB),
		eqFA, eqFB, ofA, ofB, ofAB, fmapAA, apAB, identity)
}

// ApplicativeInterchange verifies the Applicative interchange law.
func ApplicativeInterchange[FA, FB, FAB, FABB, A, B any](
	t *testing.T,
	eqFB func(FB, FB) bool,
	apAB Applicative[A, B, FA, FB, FAB],
	ptdAB Pointed[func(A) B, FAB],
	apABB Apply[func(A) B, B, FAB, FB, FABB],
	ptdABB Pointed[func(func(A) B) B, FABB],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeInterchangeEngine(t, cfg.runner, eqFB, apAB, ptdAB, apABB, ptdABB,
		enginerapid.Wrap(genA), enginerapid.Wrap(genAB))
}

// ApplicativeFullLaws verifies all four Applicative functor laws.
func ApplicativeFullLaws[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C any](
	t *testing.T,
	genFA *rapid.Generator[FA],
	genA *rapid.Generator[A],
	genAB *rapid.Generator[func(A) B],
	genBC *rapid.Generator[func(B) C],
	inst *ApplicativeInstances[FA, FB, FC, FAB, FBC, FAC, FABAC, FABB, A, B, C],
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	ApplicativeFullLawsEngine(t, cfg.runner,
		enginerapid.Wrap(genFA), enginerapid.Wrap(genA),
		enginerapid.Wrap(genAB), enginerapid.Wrap(genBC), inst)
}
```

- [ ] **Step 6: Run all existing tests**

Run: `go test ./laws/ -v`
Expected: all tests PASS

- [ ] **Step 7: Commit**

```bash
git add laws/apply_engine.go laws/applicative_engine.go laws/apply.go laws/applicative.go
git commit -m "feat(laws): add engine-generic Apply and Applicative law verification"
```

---

### Task 10: Refactor prop/ Package (Engine-Generic)

**Files:**
- Create: `prop/config.go`
- Create: `prop/roundtrip_engine.go`
- Create: `prop/oracle_engine.go`
- Modify: `prop/roundtrip.go`
- Modify: `prop/oracle.go`

- [ ] **Step 1: Create prop/config.go**

```go
package prop

import (
	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
)

// Option configures property testing behavior.
type Option func(*config)

type config struct {
	runner engine.Runner
}

// WithRunner sets the PBT engine runner for property verification.
func WithRunner(r engine.Runner) Option {
	return func(c *config) { c.runner = r }
}

func resolveConfig(opts []Option) config {
	cfg := config{runner: enginerapid.RapidRunner{}}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
```

- [ ] **Step 2: Create prop/roundtrip_engine.go**

```go
package prop

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// RoundTripEngine verifies encode/decode round-trip using the given engine.
func RoundTripEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTrip", func(et engine.T) {
		a := genA.Draw(et, "a")
		got := decode(encode(a))
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}

// RoundTripPartialEngine verifies a partial round-trip using the given engine.
func RoundTripPartialEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTripPartial", func(et engine.T) {
		a := genA.Draw(et, "a")
		got, ok := decode(encode(a))
		if !ok {
			t.Fatalf("Round-trip decode failed for input: %v", a)
		}
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}

// RoundTripErrorEngine verifies a round-trip where decode returns error using the given engine.
func RoundTripErrorEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
) {
	t.Helper()
	runner.MakeCheck(t, name+"/RoundTripError", func(et engine.T) {
		a := genA.Draw(et, "a")
		got, err := decode(encode(a))
		if err != nil {
			t.Fatalf("Round-trip decode error for input %v: %v", a, err)
		}
		if !eqA(got, a) {
			t.Fatalf("Round-trip violated:\n  original          = %v\n  decode(encode(a)) = %v", a, got)
		}
	})
}
```

- [ ] **Step 3: Create prop/oracle_engine.go**

```go
package prop

import (
	"testing"

	"github.com/franchb/fptest/engine"
)

// OracleEngine verifies implementation matches reference using the given engine.
func OracleEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqB func(B, B) bool,
	impl func(A) B,
	reference func(A) B,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Oracle", func(et engine.T) {
		a := genA.Draw(et, "a")
		got := impl(a)
		want := reference(a)
		if !eqB(got, want) {
			t.Fatalf("Oracle mismatch for input %v:\n  impl      = %v\n  reference = %v", a, got, want)
		}
	})
}

// IdempotentEngine verifies f(f(x)) == f(x) using the given engine.
func IdempotentEngine[A any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Idempotent", func(et engine.T) {
		a := genA.Draw(et, "a")
		once := f(a)
		twice := f(once)
		if !eqA(once, twice) {
			t.Fatalf("Idempotency violated for input %v:\n  f(a)    = %v\n  f(f(a)) = %v", a, once, twice)
		}
	})
}

// CommutativeEngine verifies f(a, b) == f(b, a) using the given engine.
func CommutativeEngine[A, B any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Commutative", func(et engine.T) {
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")
		left := f(a, b)
		right := f(b, a)
		if !eqB(left, right) {
			t.Fatalf("Commutativity violated:\n  f(%v, %v) = %v\n  f(%v, %v) = %v", a, b, left, b, a, right)
		}
	})
}

// InvariantEngine verifies a predicate holds for all inputs using the given engine.
func InvariantEngine[A any](
	t *testing.T,
	runner engine.Runner,
	name string,
	genA engine.Generator[A],
	predicate func(A) bool,
) {
	t.Helper()
	runner.MakeCheck(t, name+"/Invariant", func(et engine.T) {
		a := genA.Draw(et, "a")
		if !predicate(a) {
			t.Fatalf("Invariant violated for input: %v", a)
		}
	})
}
```

- [ ] **Step 4: Update prop/roundtrip.go to delegate**

```go
// Package prop provides higher-level property testing utilities built on rapid.
package prop

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// RoundTrip verifies that encode and decode are inverses: decode(encode(a)) == a.
func RoundTrip[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}

// RoundTripPartial verifies a round-trip where decode may fail, returning (A, bool).
func RoundTripPartial[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripPartialEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}

// RoundTripError verifies a round-trip where decode may return an error.
func RoundTripError[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	RoundTripErrorEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, encode, decode)
}
```

- [ ] **Step 5: Update prop/oracle.go to delegate**

```go
package prop

import (
	"testing"

	enginerapid "github.com/franchb/fptest/engine/rapid"
	"pgregory.net/rapid"
)

// Oracle verifies that an implementation matches a reference implementation.
func Oracle[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqB func(B, B) bool,
	impl func(A) B,
	reference func(A) B,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	OracleEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqB, impl, reference)
}

// Idempotent verifies that applying a function twice yields the same result as once.
func Idempotent[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	IdempotentEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqA, f)
}

// Commutative verifies that f(a, b) == f(b, a) for all a, b.
func Commutative[A, B any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	CommutativeEngine(t, cfg.runner, name, enginerapid.Wrap(genA), eqB, f)
}

// Invariant verifies that a predicate holds for all generated inputs.
func Invariant[A any](
	t *testing.T,
	name string,
	genA *rapid.Generator[A],
	predicate func(A) bool,
	opts ...Option,
) {
	t.Helper()
	cfg := resolveConfig(opts)
	InvariantEngine(t, cfg.runner, name, enginerapid.Wrap(genA), predicate)
}
```

- [ ] **Step 6: Run all tests (laws + prop + examples)**

Run: `go test ./... -v -count=1`
Expected: all tests PASS

- [ ] **Step 7: Commit**

```bash
git add prop/config.go prop/roundtrip_engine.go prop/oracle_engine.go prop/roundtrip.go prop/oracle.go
git commit -m "feat(prop): add engine-generic property testing functions"
```

---

### Task 11: Go Workspace + Hegel Sub-Module Setup

**Files:**
- Create: `go.work`
- Create: `hegel/go.mod`

- [ ] **Step 1: Create go.work**

```
go 1.25.0

use (
	.
	./hegel
)
```

- [ ] **Step 2: Create hegel/go.mod**

Run:
```bash
mkdir -p hegel
cd hegel && go mod init github.com/franchb/fptest/hegel && cd ..
```

Then add dependencies:
```bash
cd hegel && go get github.com/franchb/fptest@latest && go get hegel.dev/go/hegel@latest && cd ..
```

**Note:** During development, the workspace `go.work` will resolve `github.com/franchb/fptest` to the local `./` module automatically.

- [ ] **Step 3: Verify workspace resolves**

Run: `go work sync`
Expected: success

- [ ] **Step 4: Commit**

```bash
git add go.work hegel/go.mod hegel/go.sum
git commit -m "feat: add Go workspace and hegel sub-module"
```

---

### Task 12: Hegel Runner + Generator Adapter

**Files:**
- Create: `hegel/runner.go`
- Create: `hegel/runner_test.go`
- Create: `hegel/gen.go`
- Create: `hegel/gen_test.go`

- [ ] **Step 1: Write failing test for HegelRunner**

```go
package hegel_test

import (
	"testing"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
)

func TestHegelRunnerExecutesProperty(t *testing.T) {
	var runner engine.Runner = fpthegel.HegelRunner{}
	executed := false
	runner.MakeCheck(t, "test/property", func(et engine.T) {
		executed = true
	})
	if !executed {
		t.Fatal("property was never executed")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./hegel/ -v -run TestHegelRunner`
Expected: FAIL — types not found

- [ ] **Step 3: Implement hegel/runner.go**

```go
// Package hegel provides the hegel PBT engine adapter for fptest-go.
package hegel

import (
	"os/exec"
	"testing"

	"github.com/franchb/fptest/engine"
	hegellib "hegel.dev/go/hegel"
)

// HegelRunner implements engine.Runner using hegel.Case.
type HegelRunner struct{}

// MakeCheck runs the property as a named subtest using hegel's property checker.
func (HegelRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
	t.Helper()
	if !hegelAvailable() {
		t.Skip("hegel-core not available: install with 'uv tool install hegel-core'")
		return
	}
	t.Run(name, hegellib.Case(func(ht *hegellib.T) {
		prop(ht)
	}))
}

func hegelAvailable() bool {
	_, err := exec.LookPath("hegel-core")
	return err == nil
}
```

- [ ] **Step 4: Run test to verify it passes (requires hegel-core installed)**

Run: `go test ./hegel/ -v -run TestHegelRunner`
Expected: PASS (or SKIP if hegel-core not installed)

- [ ] **Step 5: Write failing test for HegelGen**

```go
package hegel_test

import (
	"testing"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

func TestHegelGenDrawsValues(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/draw", func(et engine.T) {
		gen := fpthegel.Wrap(hegellib.Integers[int](1, 100))
		val := gen.Draw(et, "n")
		if val < 1 || val > 100 {
			t.Fatalf("expected value in [1,100], got %d", val)
		}
	})
}
```

- [ ] **Step 6: Implement hegel/gen.go**

```go
package hegel

import (
	"github.com/franchb/fptest/engine"
	hegellib "hegel.dev/go/hegel"
)

// HegelGen wraps a hegel.Generator[A] as an engine.Generator[A].
type HegelGen[A any] struct {
	G hegellib.Generator[A]
}

// Draw produces a value from the wrapped hegel generator.
func (hg HegelGen[A]) Draw(t engine.T, _ string) A {
	return hegellib.Draw(t.(*hegellib.T), hg.G)
}

// Wrap converts a hegel.Generator[A] into an engine.Generator[A].
func Wrap[A any](g hegellib.Generator[A]) engine.Generator[A] {
	return HegelGen[A]{G: g}
}
```

**Note:** The `label` parameter is ignored because hegel's `Draw` does not take a label.

- [ ] **Step 7: Run tests**

Run: `go test ./hegel/ -v`
Expected: PASS (or SKIP)

- [ ] **Step 8: Commit**

```bash
git add hegel/runner.go hegel/runner_test.go hegel/gen.go hegel/gen_test.go
git commit -m "feat(hegel): add HegelRunner and HegelGen engine adapters"
```

---

### Task 13: Hegel Domain Generators

**Files:**
- Create: `hegel/hegelgen/domain.go`
- Create: `hegel/hegelgen/domain_test.go`

- [ ] **Step 1: Write failing tests for domain generators**

```go
package hegelgen_test

import (
	"strings"
	"testing"
	"time"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	"github.com/franchb/fptest/hegel/hegelgen"
)

func TestEmailsGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/emails", func(et engine.T) {
		email := hegelgen.Emails().Draw(et, "email")
		if !strings.Contains(email, "@") {
			t.Fatalf("expected email to contain @, got %q", email)
		}
	})
}

func TestURLsGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/urls", func(et engine.T) {
		url := hegelgen.URLs().Draw(et, "url")
		if !strings.HasPrefix(url, "http") {
			t.Fatalf("expected URL to start with http, got %q", url)
		}
	})
}

func TestDatesGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/dates", func(et engine.T) {
		d := hegelgen.Dates().Draw(et, "date")
		var zero time.Time
		if d == zero {
			t.Fatal("expected non-zero date")
		}
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./hegel/hegelgen/ -v -run "TestEmails|TestURLs|TestDates"`
Expected: FAIL — package not found

- [ ] **Step 3: Implement hegel/hegelgen/domain.go**

```go
// Package hegelgen provides hegel-backed generators wrapped as engine.Generator.
package hegelgen

import (
	"time"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

// Emails returns a generator that produces email address strings.
func Emails() engine.Generator[string] {
	return fpthegel.Wrap(hegellib.Emails())
}

// URLs returns a generator that produces RFC3986 URL strings.
func URLs() engine.Generator[string] {
	return fpthegel.Wrap(hegellib.URLs())
}

// Dates returns a generator that produces time.Time values from ISO 8601 dates.
func Dates() engine.Generator[time.Time] {
	return fpthegel.Wrap(hegellib.Dates())
}

// Datetimes returns a generator that produces time.Time values from ISO 8601 datetimes.
func Datetimes() engine.Generator[time.Time] {
	return fpthegel.Wrap(hegellib.Datetimes())
}

// FromRegex returns a generator that produces strings matching the given regex pattern.
func FromRegex(pattern string, fullmatch bool) engine.Generator[string] {
	return fpthegel.Wrap(hegellib.FromRegex(pattern, fullmatch))
}

// Text returns a generator that produces strings with codepoint count in [minSize, maxSize].
func Text(minSize, maxSize int) engine.Generator[string] {
	return fpthegel.Wrap(hegellib.Text(minSize, maxSize))
}

// Booleans returns a generator that produces boolean values.
func Booleans() engine.Generator[bool] {
	return fpthegel.Wrap(hegellib.Booleans())
}

// Integers returns a generator that produces integers in [minVal, maxVal].
func Integers[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 }](minVal, maxVal T) engine.Generator[T] {
	return fpthegel.Wrap(hegellib.Integers(minVal, maxVal))
}
```

**Note:** The exact integer constraint type should match what hegel exports. Adjust the constraint to match `hegel`'s `integer` type constraint.

- [ ] **Step 4: Run tests**

Run: `go test ./hegel/hegelgen/ -v -run "TestEmails|TestURLs|TestDates"`
Expected: PASS (or SKIP if hegel-core not available)

- [ ] **Step 5: Commit**

```bash
git add hegel/hegelgen/domain.go hegel/hegelgen/domain_test.go
git commit -m "feat(hegel/hegelgen): add domain generators (emails, URLs, dates, regex)"
```

---

### Task 14: Hegel FP-Type Generators

**Files:**
- Create: `hegel/hegelgen/fp.go`
- Create: `hegel/hegelgen/fp_test.go`

- [ ] **Step 1: Write failing tests**

```go
package hegelgen_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	"github.com/franchb/fptest/hegel/hegelgen"
	hegellib "hegel.dev/go/hegel"
)

func TestOptionGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/option", func(et engine.T) {
		opt := hegelgen.Option(hegellib.Integers[int](1, 100)).Draw(et, "opt")
		option.Fold(func() {
			// None is valid
		}, func(v int) {
			if v < 1 || v > 100 {
				t.Fatalf("expected value in [1,100], got %d", v)
			}
		})(opt)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./hegel/hegelgen/ -v -run TestOption`
Expected: FAIL

- [ ] **Step 3: Implement hegel/hegelgen/fp.go**

```go
package hegelgen

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

// Option returns a generator that produces Option[A] values (Some or None).
func Option[A any](genA hegellib.Generator[A]) engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Map(
		hegellib.Optional(genA),
		func(ptr *A) option.Option[A] {
			if ptr == nil {
				return option.None[A]()
			}
			return option.Some(*ptr)
		},
	))
}

// Some returns a generator that always produces Some[A].
func Some[A any](genA hegellib.Generator[A]) engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Map(genA, option.Some[A]))
}

// None returns a generator that always produces None[A].
func None[A any]() engine.Generator[option.Option[A]] {
	return fpthegel.Wrap(hegellib.Just(option.None[A]()))
}

// Either returns a generator that produces Either[E, A] (Left or Right).
func Either[E, A any](genE hegellib.Generator[E], genA hegellib.Generator[A]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.FlatMap(hegellib.Booleans(), func(isRight bool) hegellib.Generator[either.Either[E, A]] {
		if isRight {
			return hegellib.Map(genA, either.Right[E, A])
		}
		return hegellib.Map(genE, either.Left[A, E])
	}))
}

// Right returns a generator that always produces Right[E, A].
func Right[E, A any](genA hegellib.Generator[A]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.Map(genA, either.Right[E, A]))
}

// Left returns a generator that always produces Left[E, A].
func Left[E, A any](genE hegellib.Generator[E]) engine.Generator[either.Either[E, A]] {
	return fpthegel.Wrap(hegellib.Map(genE, either.Left[A, E]))
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./hegel/hegelgen/ -v`
Expected: PASS (or SKIP)

- [ ] **Step 5: Commit**

```bash
git add hegel/hegelgen/fp.go hegel/hegelgen/fp_test.go
git commit -m "feat(hegel/hegelgen): add FP-type generators (Option, Either)"
```

---

### Task 15: Hegel Law Convenience Wrappers

**Files:**
- Create: `hegel/laws/laws.go`
- Create: `hegel/laws/laws_test.go`

- [ ] **Step 1: Write failing test — MonoidLaws via hegel**

```go
package laws_test

import (
	"testing"

	hegellaws "github.com/franchb/fptest/hegel/laws"
	hegellib "hegel.dev/go/hegel"
	"math"
)

func TestStringMonoidLaws_Hegel(t *testing.T) {
	hegellaws.MonoidLaws(t,
		hegellib.Text(0, 50),
		func(a, b string) bool { return a == b },
		func(a, b string) string { return a + b },
		"",
	)
}

func TestIntSemigroupLaws_Hegel(t *testing.T) {
	hegellaws.SemigroupLaws(t,
		hegellib.Integers[int](math.MinInt, math.MaxInt),
		func(a, b int) bool { return a == b },
		func(a, b int) int { return a + b },
	)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./hegel/laws/ -v`
Expected: FAIL

- [ ] **Step 3: Implement hegel/laws/laws.go**

```go
// Package laws provides hegel-native convenience wrappers for fptest law verification.
package laws

import (
	"testing"

	fpthegel "github.com/franchb/fptest/hegel"
	corelaws "github.com/franchb/fptest/laws"
	hegellib "hegel.dev/go/hegel"
)

var hegelRunner = fpthegel.HegelRunner{}

// SemigroupLaws verifies Semigroup laws using hegel generators.
func SemigroupLaws[A any](
	t *testing.T,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
) {
	t.Helper()
	corelaws.SemigroupLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), eqA, concat)
}

// MonoidLaws verifies Monoid laws using hegel generators.
func MonoidLaws[A any](
	t *testing.T,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	concat func(A, A) A,
	empty A,
) {
	t.Helper()
	corelaws.MonoidLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), eqA, concat, empty)
}

// EqLaws verifies Eq laws using hegel generators.
func EqLaws[A any](
	t *testing.T,
	genA hegellib.Generator[A],
	equals func(A, A) bool,
) {
	t.Helper()
	corelaws.EqLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), equals)
}

// OrdLaws verifies Ord laws using hegel generators.
func OrdLaws[A any](
	t *testing.T,
	genA hegellib.Generator[A],
	equals func(A, A) bool,
	compare func(A, A) int,
) {
	t.Helper()
	corelaws.OrdLawsEngine(t, hegelRunner, fpthegel.Wrap(genA), equals, compare)
}

// LensLaws verifies Lens laws using hegel generators.
func LensLaws[S, A any](
	t *testing.T,
	genS hegellib.Generator[S],
	genA hegellib.Generator[A],
	eqS func(S, S) bool,
	eqA func(A, A) bool,
	get func(S) A,
	set func(A) func(S) S,
) {
	t.Helper()
	corelaws.LensLawsEngine(t, hegelRunner, fpthegel.Wrap(genS), fpthegel.Wrap(genA), eqS, eqA, get, set)
}

// FunctorLaws verifies Functor laws using hegel generators.
func FunctorLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genFA hegellib.Generator[FA],
	genAB hegellib.Generator[func(A) B],
	genBC hegellib.Generator[func(B) C],
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	fmapAA func(func(A) A) func(FA) FA,
	fmapAB func(func(A) B) func(FA) FB,
	fmapBC func(func(B) C) func(FB) FC,
	fmapAC func(func(A) C) func(FA) FC,
	identity func(A) A,
	compose func(func(A) B, func(B) C) func(A) C,
) {
	t.Helper()
	corelaws.FunctorLawsEngine(t, hegelRunner,
		fpthegel.Wrap(genFA), fpthegel.Wrap(genAB), fpthegel.Wrap(genBC),
		eqFA, eqFC, fmapAA, fmapAB, fmapBC, fmapAC, identity, compose)
}

// MonadLaws verifies Monad laws using hegel generators.
func MonadLaws[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA hegellib.Generator[A],
	genFA hegellib.Generator[FA],
	genKleisliAB hegellib.Generator[func(A) FB],
	genKleisliBC hegellib.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()
	corelaws.MonadLawsEngine(t, hegelRunner,
		fpthegel.Wrap(genA), fpthegel.Wrap(genFA),
		fpthegel.Wrap(genKleisliAB), fpthegel.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAB, chainBC, chainAC)
}

// MonadLawsFull verifies all Monad laws including right identity using hegel generators.
func MonadLawsFull[FA, FB, FC, A, B, C any](
	t *testing.T,
	genA hegellib.Generator[A],
	genFA hegellib.Generator[FA],
	genKleisliAB hegellib.Generator[func(A) FB],
	genKleisliBC hegellib.Generator[func(B) FC],
	eqFB func(FB, FB) bool,
	eqFA func(FA, FA) bool,
	eqFC func(FC, FC) bool,
	of func(A) FA,
	chainAA func(func(A) FA) func(FA) FA,
	chainAB func(func(A) FB) func(FA) FB,
	chainBC func(func(B) FC) func(FB) FC,
	chainAC func(func(A) FC) func(FA) FC,
) {
	t.Helper()
	corelaws.MonadLawsFullEngine(t, hegelRunner,
		fpthegel.Wrap(genA), fpthegel.Wrap(genFA),
		fpthegel.Wrap(genKleisliAB), fpthegel.Wrap(genKleisliBC),
		eqFB, eqFA, eqFC, of, chainAA, chainAB, chainBC, chainAC)
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./hegel/laws/ -v`
Expected: PASS (or SKIP)

- [ ] **Step 5: Commit**

```bash
git add hegel/laws/laws.go hegel/laws/laws_test.go
git commit -m "feat(hegel/laws): add hegel-native law verification convenience wrappers"
```

---

### Task 16: Hegel Prop Convenience Wrappers

**Files:**
- Create: `hegel/prop/prop.go`
- Create: `hegel/prop/prop_test.go`

- [ ] **Step 1: Write failing test**

```go
package prop_test

import (
	"strconv"
	"testing"

	fpthegel "github.com/franchb/fptest/hegel"
	hegelprop "github.com/franchb/fptest/hegel/prop"
	hegellib "hegel.dev/go/hegel"
)

func TestRoundTrip_Hegel(t *testing.T) {
	_ = fpthegel.HegelRunner{} // ensure hegel available
	hegelprop.RoundTrip(t, "itoa",
		hegellib.Integers[int](0, 999),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		func(s string) int { n, _ := strconv.Atoi(s); return n },
	)
}
```

- [ ] **Step 2: Implement hegel/prop/prop.go**

```go
// Package prop provides hegel-native convenience wrappers for fptest property testing.
package prop

import (
	"testing"

	fpthegel "github.com/franchb/fptest/hegel"
	coreprop "github.com/franchb/fptest/prop"
	hegellib "hegel.dev/go/hegel"
)

var hegelRunner = fpthegel.HegelRunner{}

// RoundTrip verifies encode/decode round-trip using hegel generators.
func RoundTrip[A, B any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) A,
) {
	t.Helper()
	coreprop.RoundTripEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// RoundTripPartial verifies a partial round-trip using hegel generators.
func RoundTripPartial[A, B any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, bool),
) {
	t.Helper()
	coreprop.RoundTripPartialEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// RoundTripError verifies a round-trip with error using hegel generators.
func RoundTripError[A, B any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	encode func(A) B,
	decode func(B) (A, error),
) {
	t.Helper()
	coreprop.RoundTripErrorEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, encode, decode)
}

// Oracle verifies implementation matches reference using hegel generators.
func Oracle[A, B any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqB func(B, B) bool,
	impl func(A) B,
	reference func(A) B,
) {
	t.Helper()
	coreprop.OracleEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqB, impl, reference)
}

// Idempotent verifies f(f(x)) == f(x) using hegel generators.
func Idempotent[A any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqA func(A, A) bool,
	f func(A) A,
) {
	t.Helper()
	coreprop.IdempotentEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqA, f)
}

// Commutative verifies f(a, b) == f(b, a) using hegel generators.
func Commutative[A, B any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	eqB func(B, B) bool,
	f func(A, A) B,
) {
	t.Helper()
	coreprop.CommutativeEngine(t, hegelRunner, name, fpthegel.Wrap(genA), eqB, f)
}

// Invariant verifies a predicate holds for all inputs using hegel generators.
func Invariant[A any](
	t *testing.T,
	name string,
	genA hegellib.Generator[A],
	predicate func(A) bool,
) {
	t.Helper()
	coreprop.InvariantEngine(t, hegelRunner, name, fpthegel.Wrap(genA), predicate)
}
```

- [ ] **Step 3: Run tests**

Run: `go test ./hegel/prop/ -v`
Expected: PASS (or SKIP)

- [ ] **Step 4: Commit**

```bash
git add hegel/prop/prop.go hegel/prop/prop_test.go
git commit -m "feat(hegel/prop): add hegel-native property testing convenience wrappers"
```

---

### Task 17: Cross-Engine Parity Tests

**Files:**
- Create: `hegel/parity_test.go`

- [ ] **Step 1: Write cross-engine parity test for Monoid**

This test runs the same monoid law verification under both rapid and hegel and confirms both engines detect the same pass/fail behavior.

```go
package hegel_test

import (
	"math"
	"testing"

	"github.com/franchb/fptest/laws"
	hegellaws "github.com/franchb/fptest/hegel/laws"
	hegellib "hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

func TestMonoidLaws_CrossEngine(t *testing.T) {
	eqInt := func(a, b int) bool { return a == b }
	add := func(a, b int) int { return a + b }

	t.Run("rapid", func(t *testing.T) {
		laws.MonoidLaws(t, rapid.Int(), eqInt, add, 0)
	})

	t.Run("hegel", func(t *testing.T) {
		hegellaws.MonoidLaws(t, hegellib.Integers[int](math.MinInt, math.MaxInt), eqInt, add, 0)
	})
}

func TestEqLaws_CrossEngine(t *testing.T) {
	eqString := func(a, b string) bool { return a == b }

	t.Run("rapid", func(t *testing.T) {
		laws.EqLaws(t, rapid.String(), eqString)
	})

	t.Run("hegel", func(t *testing.T) {
		hegellaws.EqLaws(t, hegellib.Text(0, 50), eqString)
	})
}
```

- [ ] **Step 2: Run parity tests**

Run: `go test ./hegel/ -v -run "CrossEngine"`
Expected: both rapid and hegel subtests PASS (or hegel SKIP)

- [ ] **Step 3: Commit**

```bash
git add hegel/parity_test.go
git commit -m "test: add cross-engine parity tests for rapid vs hegel"
```

---

### Task 18: Hegel Example

**Files:**
- Create: `examples/hegel_test.go`

- [ ] **Step 1: Create example test demonstrating hegel-backed law verification**

```go
package examples_test

import (
	"math"
	"testing"

	hegellaws "github.com/franchb/fptest/hegel/laws"
	hegelprop "github.com/franchb/fptest/hegel/prop"
	"github.com/franchb/fptest/hegel/hegelgen"
	hegellib "hegel.dev/go/hegel"
)

// This example demonstrates using hegel (Hypothesis-backed) as the PBT engine
// for algebraic law verification and property testing.
//
// To run: ensure hegel-core is installed (uv tool install hegel-core),
// then: go test ./examples/ -v -run TestHegel

func TestHegelMonoidLaws(t *testing.T) {
	hegellaws.MonoidLaws(t,
		hegellib.Text(0, 100),
		func(a, b string) bool { return a == b },
		func(a, b string) string { return a + b },
		"",
	)
}

func TestHegelOrdLaws(t *testing.T) {
	hegellaws.OrdLaws(t,
		hegellib.Integers[int](math.MinInt, math.MaxInt),
		func(a, b int) bool { return a == b },
		func(a, b int) int {
			if a < b { return -1 }
			if a > b { return 1 }
			return 0
		},
	)
}

func TestHegelRoundTrip(t *testing.T) {
	// Using hegel's domain generator for realistic email round-trip
	hegelprop.Invariant(t, "emails",
		hegelgen.Emails(),
		func(email string) bool {
			return len(email) > 0
		},
	)
}
```

**Note:** This file will only compile if `examples/` has access to the hegel module. If examples is in the core module, it cannot import `hegel/`. In that case, move this file to `hegel/examples_test.go` or add it as a separate test file under `hegel/`.

- [ ] **Step 2: Verify example compiles and runs**

Run: `go test ./hegel/ -v -run TestHegel` (if placed in hegel/)
Expected: PASS (or SKIP)

- [ ] **Step 3: Commit**

```bash
git add hegel/examples_test.go
git commit -m "docs: add hegel-backed law verification examples"
```

---

### Task 19: Full Test Suite + Final Verification

**Files:**
- No new files

- [ ] **Step 1: Run entire core module test suite**

Run: `go test ./engine/... ./gen/... ./laws/... ./prop/... ./assert/... ./mock/... -v -count=1`
Expected: all PASS

- [ ] **Step 2: Run entire hegel module test suite**

Run: `go test ./hegel/... -v -count=1`
Expected: all PASS (or SKIP if hegel-core not available)

- [ ] **Step 3: Run examples**

Run: `go test ./examples/... -v -count=1`
Expected: all PASS

- [ ] **Step 4: Verify go vet passes**

Run: `go vet ./... && cd hegel && go vet ./...`
Expected: no issues

- [ ] **Step 5: Commit any fixes from verification**

Only if needed. Otherwise skip.

---

### Task 20: CI Configuration for Hegel Module

**Files:**
- Modify: `.github/workflows/` (existing CI, or create new)

- [ ] **Step 1: Check existing CI configuration**

Run: `ls .github/workflows/`
Examine what's there.

- [ ] **Step 2: Add hegel test job**

Add a separate job that:
1. Installs Python 3.12+ and `uv`
2. Installs hegel-core: `uv tool install hegel-core`
3. Runs: `cd hegel && go test ./... -v`

This job should be `continue-on-error: true` or not required for merge, since it depends on external infrastructure.

- [ ] **Step 3: Verify core CI still works unchanged**

The core module job should not reference hegel or require Python.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/
git commit -m "ci: add hegel module test job with hegel-core setup"
```
