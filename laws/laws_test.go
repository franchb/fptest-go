package laws_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	"github.com/IBM/fp-go/v2/semigroup"
	"github.com/franchb/fptest/gen"
	"github.com/franchb/fptest/laws"
	"pgregory.net/rapid"
)

// --- Helpers ---

func eqOption[A any](eqA eq.Eq[A]) func(option.Option[A], option.Option[A]) bool {
	oEq := option.Eq(eqA)
	return oEq.Equals
}

func eqEither[E, A any](eqE eq.Eq[E], eqA eq.Eq[A]) func(either.Either[E, A], either.Either[E, A]) bool {
	eEq := either.Eq(eqE, eqA)
	return eEq.Equals
}

var (
	eqInt    = eq.FromStrictEquals[int]()
	eqString = eq.FromStrictEquals[string]()
	eqBool   = eq.FromStrictEquals[bool]()
)

// --- Functor Laws ---

func TestOptionFunctorLaws(t *testing.T) {
	laws.FunctorLaws[
		option.Option[int],
		option.Option[string],
		option.Option[bool],
		int, string, bool,
	](
		t,
		gen.GenOption(rapid.Int()),
		gen.GenFunc[int](rapid.String()),
		gen.GenFunc[string](rapid.Bool()),
		eqOption[int](eqInt),
		eqOption[bool](eqBool),
		option.Map[int, int],
		option.Map[int, string],
		option.Map[string, bool],
		option.Map[int, bool],
		function.Identity[int],
		func(f func(int) string, g func(string) bool) func(int) bool {
			return function.Flow2(f, g)
		},
	)
}

func TestEitherFunctorLaws(t *testing.T) {
	laws.FunctorLaws[
		either.Either[string, int],
		either.Either[string, string],
		either.Either[string, bool],
		int, string, bool,
	](
		t,
		gen.GenEither(rapid.String(), rapid.Int()),
		gen.GenFunc[int](rapid.String()),
		gen.GenFunc[string](rapid.Bool()),
		eqEither[string, int](eqString, eqInt),
		eqEither[string, bool](eqString, eqBool),
		either.Map[string, int, int],
		either.Map[string, int, string],
		either.Map[string, string, bool],
		either.Map[string, int, bool],
		function.Identity[int],
		func(f func(int) string, g func(string) bool) func(int) bool {
			return function.Flow2(f, g)
		},
	)
}

// --- Monad Laws ---

func TestOptionMonadLaws(t *testing.T) {
	laws.MonadLawsFull[
		option.Option[int],
		option.Option[string],
		option.Option[bool],
		int, string, bool,
	](
		t,
		rapid.Int(),
		gen.GenOption(rapid.Int()),
		gen.GenFunc[int](gen.GenOption(rapid.String())),
		gen.GenFunc[string](gen.GenOption(rapid.Bool())),
		eqOption[string](eqString),
		eqOption[int](eqInt),
		eqOption[bool](eqBool),
		option.Of[int],
		option.Chain[int, int],
		option.Chain[int, string],
		option.Chain[string, bool],
		option.Chain[int, bool],
	)
}

// --- Semigroup Laws ---

func TestIntSumSemigroupLaws(t *testing.T) {
	sg := semigroup.MakeSemigroup(func(a, b int) int { return a + b })
	laws.SemigroupInterfaceLaws(t, rapid.Int(), eqInt.Equals, sg)
}

func TestStringSemigroupLaws(t *testing.T) {
	sg := semigroup.MakeSemigroup(func(a, b string) string { return a + b })
	laws.SemigroupInterfaceLaws(t, rapid.String(), eqString.Equals, sg)
}

// --- Monoid Laws ---

func TestIntSumMonoidLaws(t *testing.T) {
	m := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
	laws.MonoidInterfaceLaws(t, rapid.Int(), eqInt.Equals, m)
}

func TestStringMonoidLaws(t *testing.T) {
	m := monoid.MakeMonoid(func(a, b string) string { return a + b }, "")
	laws.MonoidInterfaceLaws(t, rapid.String(), eqString.Equals, m)
}

// --- Eq Laws ---

func TestIntEqLaws(t *testing.T) {
	laws.EqInterfaceLaws(t, rapid.Int(), eqInt)
}

func TestStringEqLaws(t *testing.T) {
	laws.EqInterfaceLaws(t, rapid.String(), eqString)
}

func TestOptionEqLaws(t *testing.T) {
	laws.EqLaws(t, gen.GenOption(rapid.Int()), eqOption[int](eqInt))
}

// --- Ord Laws ---

func TestIntOrdLaws(t *testing.T) {
	ordInt := ord.FromCompare(func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	laws.OrdInterfaceLaws(t, rapid.Int(), ordInt)
}

// --- Lens Laws ---

type Person struct {
	Name string
	Age  int
}

func TestPersonNameLensLaws(t *testing.T) {
	genPerson := rapid.Custom(func(t *rapid.T) Person {
		return Person{
			Name: rapid.String().Draw(t, "name"),
			Age:  rapid.IntRange(0, 150).Draw(t, "age"),
		}
	})

	laws.LensLaws(t, genPerson, rapid.String(),
		func(a, b Person) bool { return a == b },
		func(a, b string) bool { return a == b },
		func(p Person) string { return p.Name },
		func(name string) func(Person) Person {
			return func(p Person) Person {
				p.Name = name
				return p
			}
		},
	)
}

// --- Applicative Laws ---

func TestOptionApplicativeLaws(t *testing.T) {
	laws.ApplicativeLaws[
		option.Option[int],
		option.Option[string],
		option.Option[func(int) string],
		int, string,
	](
		t,
		rapid.Int(),
		gen.GenOption(rapid.Int()),
		gen.GenFunc[int](rapid.String()),
		eqOption[int](eqInt),
		eqOption[string](eqString),
		option.Of[int],
		option.Of[string],
		option.Of[func(int) string],
		option.Map[int, int],
		option.Ap[string, int],
		function.Identity[int],
	)
}

