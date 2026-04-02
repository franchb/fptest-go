package laws

// Pointed abstracts a pure/of operation that lifts a value into a type constructor.
type Pointed[A, FA any] interface {
	Of(A) FA
}

// Functor abstracts the Map operation for a type constructor.
type Functor[A, B, FA, FB any] interface {
	Map(func(A) B) func(FA) FB
}

// Apply abstracts Functor + Ap for a type constructor.
// Ap follows the fp-go value-first convention: Ap(fa) returns func(fab) fb.
type Apply[A, B, FA, FB, FAB any] interface {
	Functor[A, B, FA, FB]
	Ap(FA) func(FAB) FB
}

// Applicative combines Apply and Pointed for a type constructor.
type Applicative[A, B, FA, FB, FAB any] interface {
	Apply[A, B, FA, FB, FAB]
	Pointed[A, FA]
}

// Chainable abstracts the monadic bind (Chain) operation.
type Chainable[A, B, FA, FB any] interface {
	Chain(func(A) FB) func(FA) FB
}

// Monad combines Applicative and Chainable for a type constructor.
type Monad[A, B, FA, FB, FAB any] interface {
	Applicative[A, B, FA, FB, FAB]
	Chainable[A, B, FA, FB]
}

// --- concrete adapters ---

type pointed[A, FA any] struct {
	of func(A) FA
}

func (p pointed[A, FA]) Of(a A) FA { return p.of(a) }

// MakePointed constructs a Pointed from a pure/of function.
//
// Type parameters:
//   - A: element type
//   - FA: container type F[A]
//
// In fp-go terms: of = option.Of[int].
func MakePointed[A, FA any](of func(A) FA) Pointed[A, FA] {
	return pointed[A, FA]{of: of}
}

type functor[A, B, FA, FB any] struct {
	fmap func(func(A) B) func(FA) FB
}

func (f functor[A, B, FA, FB]) Map(ab func(A) B) func(FA) FB { return f.fmap(ab) }

// MakeFunctor constructs a Functor from a curried map function.
//
// Type parameters:
//   - A: element type
//   - B: target element type
//   - FA: container type F[A]
//   - FB: container type F[B]
//
// In fp-go terms: fmap = option.Map[int, string].
func MakeFunctor[A, B, FA, FB any](fmap func(func(A) B) func(FA) FB) Functor[A, B, FA, FB] {
	return functor[A, B, FA, FB]{fmap: fmap}
}

type apply[A, B, FA, FB, FAB any] struct {
	fmap func(func(A) B) func(FA) FB
	ap   func(FA) func(FAB) FB
}

func (a apply[A, B, FA, FB, FAB]) Map(ab func(A) B) func(FA) FB { return a.fmap(ab) }
func (a apply[A, B, FA, FB, FAB]) Ap(fa FA) func(FAB) FB       { return a.ap(fa) }

// MakeApply constructs an Apply from curried map and ap functions.
//
// Type parameters:
//   - A: element type
//   - B: target element type
//   - FA: container type F[A]
//   - FB: container type F[B]
//   - FAB: container type F[func(A) B]
//
// In fp-go terms: fmap = option.Map[int, string],
// ap = option.Ap[string, int].
func MakeApply[A, B, FA, FB, FAB any](
	fmap func(func(A) B) func(FA) FB,
	ap func(FA) func(FAB) FB,
) Apply[A, B, FA, FB, FAB] {
	return apply[A, B, FA, FB, FAB]{fmap: fmap, ap: ap}
}

type applicative[A, B, FA, FB, FAB any] struct {
	of   func(A) FA
	fmap func(func(A) B) func(FA) FB
	ap   func(FA) func(FAB) FB
}

func (a applicative[A, B, FA, FB, FAB]) Of(v A) FA                  { return a.of(v) }
func (a applicative[A, B, FA, FB, FAB]) Map(ab func(A) B) func(FA) FB { return a.fmap(ab) }
func (a applicative[A, B, FA, FB, FAB]) Ap(fa FA) func(FAB) FB      { return a.ap(fa) }

// MakeApplicative constructs an Applicative from of, map, and ap functions.
//
// Type parameters:
//   - A: element type
//   - B: target element type
//   - FA: container type F[A]
//   - FB: container type F[B]
//   - FAB: container type F[func(A) B]
//
// In fp-go terms: of = option.Of[int], fmap = option.Map[int, string],
// ap = option.Ap[string, int].
func MakeApplicative[A, B, FA, FB, FAB any](
	of func(A) FA,
	fmap func(func(A) B) func(FA) FB,
	ap func(FA) func(FAB) FB,
) Applicative[A, B, FA, FB, FAB] {
	return applicative[A, B, FA, FB, FAB]{of: of, fmap: fmap, ap: ap}
}

// --- Chainable and Monad concrete adapters ---

type chainable[A, B, FA, FB any] struct {
	chain func(func(A) FB) func(FA) FB
}

func (c chainable[A, B, FA, FB]) Chain(f func(A) FB) func(FA) FB { return c.chain(f) }

// MakeChainable constructs a Chainable from a monadic bind (chain) function.
//
// Type parameters:
//   - A: element type
//   - B: target element type
//   - FA: container type F[A]
//   - FB: container type F[B]
//
// In fp-go terms: chain = option.Chain[int, string].
func MakeChainable[A, B, FA, FB any](chain func(func(A) FB) func(FA) FB) Chainable[A, B, FA, FB] {
	return chainable[A, B, FA, FB]{chain: chain}
}

type monad[A, B, FA, FB, FAB any] struct {
	of    func(A) FA
	fmap  func(func(A) B) func(FA) FB
	ap    func(FA) func(FAB) FB
	chain func(func(A) FB) func(FA) FB
}

func (m monad[A, B, FA, FB, FAB]) Of(a A) FA                      { return m.of(a) }
func (m monad[A, B, FA, FB, FAB]) Map(f func(A) B) func(FA) FB    { return m.fmap(f) }
func (m monad[A, B, FA, FB, FAB]) Ap(fa FA) func(FAB) FB          { return m.ap(fa) }
func (m monad[A, B, FA, FB, FAB]) Chain(f func(A) FB) func(FA) FB { return m.chain(f) }

// MakeMonad constructs a Monad from of, map, ap, and chain functions.
//
// Type parameters:
//   - A: element type
//   - B: target element type
//   - FA: container type F[A]
//   - FB: container type F[B]
//   - FAB: container type F[func(A) B]
//
// In fp-go terms: of = option.Of[int], fmap = option.Map[int, string],
// ap = option.Ap[string, int], chain = option.Chain[int, string].
func MakeMonad[A, B, FA, FB, FAB any](
	of func(A) FA,
	fmap func(func(A) B) func(FA) FB,
	ap func(FA) func(FAB) FB,
	chain func(func(A) FB) func(FA) FB,
) Monad[A, B, FA, FB, FAB] {
	return monad[A, B, FA, FB, FAB]{of: of, fmap: fmap, ap: ap, chain: chain}
}
