## Why fptest-go

### The problem: algebraic contracts are invisible at compile time

When you write `option.Map(f)` or `either.Chain(g)`, you rely on an implicit contract: these operations obey algebraic laws. `Map(id) == id`. `Chain` is associative. `Of` is a left and right identity for `Chain`. These aren't style preferences — they are the axioms that make functional composition predictable. If `Map` breaks the functor identity law, then refactoring `Map(id)` into a no-op changes behavior. If `Chain` isn't associative, then regrouping a monadic pipeline produces different results depending on parenthesization. Your code compiles, your types check, and your program is subtly wrong.

Go's type system — even with generics — cannot express these contracts. There is no `Functor` constraint, no typeclass coherence check, no compiler-enforced law. In Haskell, libraries like [hedgehog-classes](https://hackage.haskell.org/package/hedgehog-classes) and [QuickCheck-classes](https://hackage.haskell.org/package/quickcheck-classes) exist precisely because even Haskell's type system cannot enforce laws — it can only enforce types. In Scala, [cats-discipline](https://typelevel.org/cats/typeclasses/lawtesting.html) tests every typeclass instance against its laws using ScalaCheck. In TypeScript, [fp-ts-laws](https://github.com/gcanti/fp-ts-laws) does the same with fast-check.

fp-go has **none of this** as a public API. It ships 23+ internal law-testing packages (`internal/functor/testing`, `internal/monad/testing`, etc.) that verify its own implementations — but these are invisible to you behind Go's package boundary. When you build your own `Monoid`, your own `Lens`, your own `ReaderIOEither` pipeline, you are on your own. fptest-go fills this gap.

### What algebraic laws actually guarantee

Each set of laws eliminates a specific class of bugs:

| Laws | What they prevent |
|---|---|
| **Functor** (identity, composition) | `Map` that silently mutates state, drops elements, reorders, or fails to distribute over function composition |
| **Monad** (left/right identity, associativity) | `Chain`/`FlatMap` pipelines that produce different results depending on how you group them; `Of` that doesn't act as a neutral element |
| **Applicative** (identity, homomorphism, interchange, composition) | `Ap`-based parallel composition that disagrees with sequential `Chain` or corrupts independent computations |
| **Monoid** (identity, associativity) | `Concat` that fails when one operand is `Empty`; fold/reduce operations that depend on chunking order |
| **Semigroup** (associativity) | Parallel aggregation that produces different results depending on how work is partitioned across goroutines |
| **Eq** (reflexivity, symmetry, transitivity) | Equality that isn't an equivalence relation — breaks maps, sets, deduplication, caching, and any code that assumes `a == b && b == c → a == c` |
| **Ord** (totality, antisymmetry, transitivity, Eq-consistency) | Sort instability, binary search failures, priority queue corruption, incorrect min/max |
| **Lens** (get-set, set-get, set-set) | Optics that lose data on round-trip, ignore the value you just set, or depend on how many times you set |

These are not theoretical concerns. A `Monoid` whose `Empty` isn't a true identity breaks `FoldMap` over an empty collection. A `Lens` that violates set-get means your UI state management silently discards user input. A non-associative `Semigroup` means your concurrent aggregation pipeline produces non-deterministic results depending on goroutine scheduling.

### Why property-based testing, not example-based testing

Example-based tests (`assert.Equal(t, f(3), 9)`) check specific input-output pairs. They cannot find the edge case you didn't think of. Property-based testing inverts the approach: you state what must *always* be true, and the framework searches for counterexamples.

Consider testing a `Monoid[Permissions]` for a bitflag type:

```go
// Example-based: checks three cases you thought of
func TestPermMonoid(t *testing.T) {
    assert.Equal(t, Concat(Read, Write), ReadWrite)
    assert.Equal(t, Concat(Empty, Read), Read)
    assert.Equal(t, Concat(Admin, Empty), Admin)
}

// Property-based: checks the *contract* across hundreds of random inputs
func TestPermMonoid(t *testing.T) {
    laws.MonoidLaws(t, genPerm, Empty, Concat, eqPerm)
    // Runs associativity, left identity, right identity
    // with shrinking — if Concat(a, Concat(b, c)) ≠ Concat(Concat(a, b), c)
    // for some a,b,c, you get the *minimal* counterexample
}
```

The property-based version tests the invariant, not individual cases. Rapid's automatic shrinking (the same bitstream-based approach used by Hypothesis) means when a law violation is found, the reported counterexample is minimal — not a 500-element slice, but the smallest inputs that trigger the failure.

### Why not just `pgregory.net/rapid` alone

Rapid is an excellent property-based testing engine. fptest-go does not replace it — it builds on it. What rapid gives you is generators and a test runner. What it does *not* give you:

**No fp-go type awareness.** Rapid cannot generate `Option[A]`, `Either[E, A]`, `IO[A]`, or `IOEither[E, A]` out of the box. You would write `rapid.Custom(func(t *rapid.T) Option[int] { ... })` by hand, every time, in every test file. fptest-go provides `gen.GenOption`, `gen.GenResult`, `gen.GenIO`, `gen.GenIOResult` as composable, reusable generators.

**No law-testing infrastructure.** Rapid gives you `Check(t, property)`. It does not give you "verify that this `Map` function satisfies the functor laws." You would need to manually write the identity law, the composition law, generate random functions, compare results with a custom equality — and repeat for every typeclass, every type. fptest-go encodes each law once and exposes it as a single function call.

**No monadic generator composition.** Rapid's `*Generator[V]` is an opaque type. You cannot `Chain` two generators to express dependent generation (generate `lo`, then generate `hi > lo`) without dropping into the imperative `rapid.Custom` + `Draw` style. fptest-go's `Gen[A] = func(*rapid.T) A` is a transparent function type with a full `Functor`/`Applicative`/`Monad` — you can `Map`, `Chain`, and `Ap` generators the same way you compose fp-go pipelines.

**No assertion helpers for sum types.** Rapid does not know that `Either` has a `Right` and a `Left`. You would `Fold` manually and call `t.Fatal` yourself. fptest-go provides `AssertRight`, `AssertSome`, `AssertIORight` etc. that unwrap and return the inner value for further chaining.

### Why not `testing` + `testify/assert`

The standard library's `testing` package and `testify/assert` are designed for Go's idiomatic `(value, error)` world. They have no concept of `Option`, `Either`, `IO`, or any algebraic structure.

`testify/assert.Equal(t, got, want)` uses `reflect.DeepEqual`. This fails on function types (`IO[A] = func() A` is always unequal via reflection), panics on certain recursive structures, and provides no way to plug in a custom `Eq` instance. When your value is `IOEither[error, User]` — a function returning a function returning a sum type — `reflect.DeepEqual` is meaningless.

`testify/assert` also cannot express *properties*. It asserts specific values. There is no `assert.ForAll`, no shrinking, no law verification. It is a point-check tool, not a contract-verification tool.

fptest-go integrates with `*testing.T` and `t.Run` — your law tests appear as standard Go subtests (`TestMyType/Functor/Identity`, `TestMyType/Monad/Associativity`). They run with `go test`, filter with `-run`, and report with `t.Helper()` pointing to your call site. You don't leave the Go testing ecosystem — you extend it with algebraic rigor.

### Concrete use cases

**You wrote a custom `Monoid` for merging configs.** How do you know `Concat(Empty, x) == x` for all `x`? That `Concat(Concat(a, b), c) == Concat(a, Concat(b, c))` even when configs contain nested optionals? One call to `laws.MonoidLaws` verifies both properties across hundreds of random configs with automatic shrinking.

**You built an optics `Lens` for a deeply nested struct.** Does `Set(Get(whole), whole) == whole`? Does `Get(Set(part, whole)) == part`? If your `Set` function has an off-by-one in a slice index, `laws.LensLaws` will find the minimal struct that breaks it.

**You have a `ReaderIOEither` pipeline for an HTTP handler.** You want to verify that your error-handling middleware composes correctly — that mapping over it preserves identity, that chaining is associative. `laws.FunctorLaws` and `laws.MonadLaws` verify this with generated inputs, without hitting a real HTTP server.

**You serialize domain types to JSON and back.** `prop.RoundTrip` generates random domain values, encodes them, decodes them, and verifies equality — catching encoding bugs that a hand-picked example would miss (empty strings, Unicode edge cases, zero values, nil slices vs empty slices).

**You implemented `Eq` for a case-insensitive string type.** Is it actually an equivalence relation? `laws.EqLaws` checks reflexivity, symmetry, and transitivity — catching the classic bug where `Equals("ß", "SS")` is true but `Equals("SS", "ß")` is false due to asymmetric Unicode case folding.

**You have a concurrent aggregation pipeline using `Semigroup.Concat`.** If `Concat` isn't associative, the result depends on how goroutines are scheduled. `laws.SemigroupLaws` catches this before it becomes a production Heisenbug.
