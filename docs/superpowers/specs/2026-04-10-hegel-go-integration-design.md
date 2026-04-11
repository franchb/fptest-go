# hegel-go Integration Design

**Date:** 2026-04-10
**Status:** Draft
**Module:** `github.com/franchb/fptest-go`

## Problem

fptest-go is a property-based testing and algebraic law verification library built on rapid (pure Go PBT engine). Users who prefer hegel-go (Hypothesis-backed PBT with superior shrinking and rich domain generators) cannot use fptest-go's law verification and property testing facilities.

## Goals

1. **Multi-engine support** — users choose rapid or hegel as their PBT backend
2. **Better shrinking** — leverage Hypothesis's shrinking for law verification
3. **Domain generators** — expose hegel's built-in generators (emails, URLs, dates, regex) to fptest-go users
4. **Ecosystem reach** — fptest-go is usable by teams that chose hegel over rapid

## Non-Goals

- Replacing rapid as the default engine
- Supporting other PBT engines (gopter, etc.) in this iteration — the abstraction enables it but we only implement rapid + hegel
- Making `gen.Gen[A]` engine-agnostic (it stays `func(*rapid.T) A`)

## Architecture: Monorepo with Go Workspace

Single repository, two Go modules connected by a Go workspace.

```
fptest-go/
├── go.work                      # workspace: . + ./hegel
├── go.mod                       # github.com/franchb/fptest-go (rapid + fp-go)
├── engine/
│   ├── engine.go                # T, Generator[A], Runner interfaces
│   └── rapid/
│       └── rapid.go             # RapidRunner, RapidGen[A] adapter
├── gen/
│   └── adapt.go                 # Gen[A] -> engine.Generator[A] adapter
├── laws/                        # Updated: optional ...Option for Runner
├── prop/                        # Updated: optional ...Option for Runner
├── assert/                      # Unchanged
├── mock/                        # Unchanged
├── hegel/
│   ├── go.mod                   # github.com/franchb/fptest-go/hegel
│   ├── runner.go                # HegelRunner implements engine.Runner
│   ├── gen.go                   # hegel.Generator[T] -> engine.Generator[T]
│   ├── hegelgen/
│   │   ├── domain.go            # Emails(), URLs(), Dates(), FromRegex()
│   │   └── fp.go                # hegel generators for Option, Either, IOEither
│   └── laws/                    # Convenience wrappers with hegel-native signatures
└── examples/
    └── hegel_test.go            # Example: law verification with hegel backend
```

**Dependency graph:**
- `engine/` — zero external deps (only `testing`)
- `engine/rapid/` — depends on `pgregory.net/rapid`
- `laws/`, `prop/` — depend on `engine/` + `engine/rapid/` (for default runner)
- `hegel/` — depends on `engine/` + `hegel.dev/go/hegel`

## Engine Abstraction (`engine/`)

Three interfaces define the PBT engine contract:

```go
package engine

import "testing"

// T is the test context provided by a PBT engine during property execution.
// Both rapid.T and hegel.T satisfy this via adapters.
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

### Rapid Adapter (`engine/rapid/`)

```go
package rapid

import (
    "testing"

    "github.com/franchb/fptest-go/engine"
    rapidlib "pgregory.net/rapid"
)

// RapidRunner implements engine.Runner using rapid.MakeCheck.
type RapidRunner struct{}

func (RapidRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
    t.Run(name, rapidlib.MakeCheck(func(rt *rapidlib.T) {
        prop(rt) // rapid.T satisfies engine.T (it embeds testing.TB)
    }))
}

// RapidGen wraps *rapid.Generator[A] as engine.Generator[A].
type RapidGen[A any] struct {
    G *rapidlib.Generator[A]
}

func (rg RapidGen[A]) Draw(t engine.T, label string) A {
    return rg.G.Draw(t.(*rapidlib.T), label)
}

func Wrap[A any](g *rapidlib.Generator[A]) engine.Generator[A] {
    return RapidGen[A]{G: g}
}
```

## Generator Interop

### Conversion directions

| From | To | How | Location |
|------|----|-----|----------|
| `*rapid.Generator[A]` | `engine.Generator[A]` | `rapid.Wrap(g)` | `engine/rapid/` |
| `gen.Gen[A]` | `engine.Generator[A]` | `gen.ToEngine(g)` (via `gen.ToRapid` then `rapid.Wrap`) | `gen/adapt.go` |
| `hegel.Generator[T]` | `engine.Generator[T]` | `hegelmod.Wrap(g)` | `hegel/gen.go` |
| `gen.Gen[A]` -> `hegel.Generator` | **Not supported** | `Gen[A]` is inherently rapid-native | — |

### Design decision

`gen.Gen[A]` stays as `func(*rapid.T) A`. It is the rapid-native generator monad. Hegel users work with `engine.Generator[A]` directly. This avoids a leaky abstraction and keeps each engine's idioms natural.

### Domain generators from hegel

```go
// hegel/hegelgen/domain.go
package hegelgen

func Emails() engine.Generator[string]           { return Wrap(hegel.Emails()) }
func URLs() engine.Generator[string]             { return Wrap(hegel.URLs()) }
func Dates() engine.Generator[time.Time]         { return Wrap(hegel.Dates()) }
func Datetimes() engine.Generator[time.Time]     { return Wrap(hegel.Datetimes()) }
func FromRegex(pat string) engine.Generator[string] { return Wrap(hegel.FromRegex(pat)) }
func IPAddresses() engine.Generator[string]      { return Wrap(hegel.IPAddresses()) }
```

### FP-type generators from hegel

```go
// hegel/hegelgen/fp.go
package hegelgen

func Option[A any](genA hegel.Generator[A]) engine.Generator[option.Option[A]]
func Some[A any](genA hegel.Generator[A]) engine.Generator[option.Option[A]]
func None[A any]() engine.Generator[option.Option[A]]
func Either[E, A any](genE hegel.Generator[E], genA hegel.Generator[A]) engine.Generator[either.Either[E, A]]
func Right[E, A any](genA hegel.Generator[A]) engine.Generator[either.Either[E, A]]
func Left[E, A any](genE hegel.Generator[E]) engine.Generator[either.Either[E, A]]
```

## Law & Property Function Migration

### Backwards-compatible API evolution

Law functions gain an optional variadic `...Option` parameter:

```go
// laws/option.go
type Option func(*config)

type config struct {
    runner engine.Runner
}

func WithRunner(r engine.Runner) Option {
    return func(c *config) { c.runner = r }
}

func resolveConfig(opts []Option) config {
    cfg := config{runner: rapid.RapidRunner{}} // default: rapid
    for _, o := range opts {
        o(&cfg)
    }
    return cfg
}
```

**Example — FunctorLaws before:**
```go
func FunctorLaws[FA, FB, FC, A, B, C any](
    t *testing.T,
    genFA *rapid.Generator[FA],
    // ...
)
```

**After:**
```go
func FunctorLaws[FA, FB, FC, A, B, C any](
    t *testing.T,
    genFA *rapid.Generator[FA],
    // ...same params...
    opts ...Option,
)
```

Internally, `FunctorLaws` wraps the `*rapid.Generator` params via `rapid.Wrap()` and delegates to an internal `functorLawsEngine()` that uses `engine.Generator[A]` + `engine.Runner`.

Existing callers pass no options and get rapid (unchanged behavior).

### Engine-generic API (advanced)

For users who want to mix engines or use the abstraction directly, each law function also has an engine-generic internal variant exposed as a public function:

```go
// laws/functor_engine.go
func FunctorLawsEngine[FA, FB, FC, A, B, C any](
    t *testing.T,
    runner engine.Runner,
    genFA engine.Generator[FA],
    genAB engine.Generator[func(A) B],
    genBC engine.Generator[func(B) C],
    eqFA func(FA, FA) bool,
    eqFC func(FC, FC) bool,
    // ...operations...
)
```

This is the foundation that both the rapid-native API and hegel wrappers delegate to.

### Hegel convenience wrappers

```go
// hegel/laws/functor.go
package laws

import (
    "testing"

    hegelmod "github.com/franchb/fptest-go/hegel"
    corelaws "github.com/franchb/fptest-go/laws"
    "hegel.dev/go/hegel"
)

func FunctorLaws[FA, FB, FC, A, B, C any](
    t *testing.T,
    genFA hegel.Generator[FA],
    genAB hegel.Generator[func(A) B],
    genBC hegel.Generator[func(B) C],
    eqFA func(FA, FA) bool,
    eqFC func(FC, FC) bool,
    // ...operations...
) {
    // Convert hegel generators to engine.Generator, call core with HegelRunner
    // Implementation delegates to the engine-parametric internal functions
}
```

This pattern repeats for all law functions: MonadLaws, ApplicativeLaws, MonoidLaws, SemigroupLaws, EqLaws, OrdLaws, LensLaws, ChainAssociativity, ApplyAssociativeComposition.

And for all prop functions: RoundTrip, RoundTripPartial, RoundTripError, Oracle, Idempotent, Commutative, Invariant.

## Packages Unchanged

- **`assert/`** — operates on `testing.TB` and fp-go types directly. No PBT engine involvement.
- **`mock/`** — `IORef`, `CallTracker`, stubs are engine-independent. No changes needed.

## Testing Strategy

### Unit tests

- `engine/rapid/` — verify RapidRunner and RapidGen correctly delegate to rapid
- `hegel/` — verify HegelRunner and HegelGen correctly delegate to hegel

### Cross-engine parity tests

Run identical law verification under both engines and confirm:
- Both pass for correct implementations (e.g., string monoid with `+` and `""`)
- Both fail for broken implementations (e.g., monoid with wrong identity)

```go
func TestMonoidLaws_CrossEngine(t *testing.T) {
    eqString := func(a, b string) bool { return a == b }
    concat := func(a, b string) string { return a + b }
    empty := ""

    // rapid: uses *rapid.Generator directly (existing API)
    t.Run("rapid", func(t *testing.T) {
        genString := rapid.String()
        laws.MonoidLaws(t, genString, eqString, concat, empty)
    })

    // hegel: uses hegel convenience wrappers (hegel-native generators)
    t.Run("hegel", func(t *testing.T) {
        genString := hegel.Text(0, 100)
        hegellaws.MonoidLaws(t, genString, eqString, concat, empty)
    })
}
```

Note: the two engine paths use different entry points. The core `laws.MonoidLaws` accepts `*rapid.Generator` (backwards-compatible). The `hegel/laws.MonoidLaws` wrapper accepts `hegel.Generator` and internally delegates to the engine-parametric implementation with `HegelRunner`.

### CI configuration

- **Core module CI:** runs as today (rapid only, fast, no external deps)
- **Hegel module CI:** separate job that installs Python + `uv` + hegel-core, then runs `hegel/` tests
- The hegel CI job is non-blocking for core module releases

### Graceful degradation

If hegel-core is not available at test time:
```go
func (HegelRunner) MakeCheck(t *testing.T, name string, prop func(engine.T)) {
    if !hegelAvailable() {
        t.Skip("hegel-core not available: install with 'uv tool install hegel-core'")
        return
    }
    // ...
}
```

## Benchmarks (Stretch Goal)

Benchmark suite in `hegel/bench/` comparing:
- Wall-clock time for identical law suites (MonoidLaws, FunctorLaws, etc.)
- Iterations per second
- Memory per iteration
- Shrinking time on intentionally broken laws

Output as Go benchmark results for `benchstat` comparison.

## Fuzz Interop (Stretch Goal, Exploratory)

Potential bridge: generate N interesting values from either engine and seed Go's native `testing.F` corpus. Would live in a `fuzz/` package in the core module. Not blocking initial release.

## Migration Path

1. Add `engine/` package with interfaces (no existing code changes)
2. Add `engine/rapid/` adapter (no existing code changes)
3. Add `gen/adapt.go` for `Gen[A]` -> `engine.Generator[A]`
4. Refactor `laws/` and `prop/` internals to use engine abstraction, keeping public API backwards-compatible via `...Option`
5. Create `hegel/` sub-module with runner, generator adapters, and convenience wrappers
6. Add cross-engine parity tests
7. Add examples and documentation
8. (Stretch) Add benchmarks and fuzz interop

## Open Questions

- **hegel-go stability:** hegel-go is beta (v0.1.x). Breaking changes in hegel-go's API may require updates to the adapter layer. Pin to a specific version and document compatibility.
- **Go 1.25 requirement:** Both fptest-go and hegel-go require Go 1.25+. This is aligned — no issue.
- **hegel.T type assertion:** The engine.T interface requires type-asserting to the concrete engine T (rapid.T or hegel.T) inside adapters. This is safe because the Runner controls which T it passes to the property function, but it's a runtime invariant not enforced by the type system.
