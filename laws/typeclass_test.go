package laws_test

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/gen"
	"github.com/franchb/fptest/laws"
	"pgregory.net/rapid"
)

func TestPointedOptionCompiles(t *testing.T) {
	ptd := laws.MakePointed(option.Of[int])
	got := ptd.Of(42)
	if !option.IsSome(got) {
		t.Fatal("expected Some(42), got None")
	}
	v := option.GetOrElse(func() int { return 0 })(got)
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestFunctorOptionCompiles(t *testing.T) {
	f := laws.MakeFunctor(option.Map[int, string])
	got := f.Map(func(i int) string { return fmt.Sprintf("%d", i) })(option.Of(42))
	if !option.IsSome(got) {
		t.Fatal("expected Some(\"42\"), got None")
	}
	v := option.GetOrElse(func() string { return "" })(got)
	if v != "42" {
		t.Fatalf("expected \"42\", got %q", v)
	}
}

func TestApplyOptionCompiles(t *testing.T) {
	ap := laws.MakeApply[int, string, option.Option[int], option.Option[string], option.Option[func(int) string]](
		option.Map[int, string],
		option.Ap[string, int],
	)

	// Test Map
	mapped := ap.Map(func(i int) string { return fmt.Sprintf("%d", i) })(option.Of(10))
	if !option.IsSome(mapped) {
		t.Fatal("Map: expected Some, got None")
	}
	v := option.GetOrElse(func() string { return "" })(mapped)
	if v != "10" {
		t.Fatalf("Map: expected \"10\", got %q", v)
	}

	// Test Ap
	fn := option.Of(func(i int) string { return fmt.Sprintf("v=%d", i) })
	applied := ap.Ap(option.Of(7))(fn)
	if !option.IsSome(applied) {
		t.Fatal("Ap: expected Some, got None")
	}
	v = option.GetOrElse(func() string { return "" })(applied)
	if v != "v=7" {
		t.Fatalf("Ap: expected \"v=7\", got %q", v)
	}
}

func TestApplicativeOptionCompiles(t *testing.T) {
	ap := laws.MakeApplicative[int, string, option.Option[int], option.Option[string], option.Option[func(int) string]](
		option.Of[int],
		option.Map[int, string],
		option.Ap[string, int],
	)

	// Test Of
	fa := ap.Of(99)
	if !option.IsSome(fa) {
		t.Fatal("Of: expected Some(99), got None")
	}
	vi := option.GetOrElse(func() int { return 0 })(fa)
	if vi != 99 {
		t.Fatalf("Of: expected 99, got %d", vi)
	}

	// Test Map
	mapped := ap.Map(func(i int) string { return fmt.Sprintf("%d!", i) })(fa)
	vs := option.GetOrElse(func() string { return "" })(mapped)
	if vs != "99!" {
		t.Fatalf("Map: expected \"99!\", got %q", vs)
	}

	// Test Ap
	fn := option.Of(func(i int) string { return fmt.Sprintf("n=%d", i) })
	applied := ap.Ap(option.Of(5))(fn)
	vs = option.GetOrElse(func() string { return "" })(applied)
	if vs != "n=5" {
		t.Fatalf("Ap: expected \"n=5\", got %q", vs)
	}
}

func TestChainableOptionCompiles(t *testing.T) {
	ch := laws.MakeChainable[int, string, option.Option[int], option.Option[string]](
		option.Chain[int, string],
	)
	result := ch.Chain(func(i int) option.Option[string] {
		return option.Of(fmt.Sprintf("%d", i))
	})(option.Of(42))
	if !option.IsSome(result) {
		t.Fatal("expected Some(\"42\"), got None")
	}
	v := option.GetOrElse(func() string { return "" })(result)
	if v != "42" {
		t.Fatalf("expected \"42\", got %q", v)
	}

	// Also verify None propagation
	noneResult := ch.Chain(func(i int) option.Option[string] {
		return option.Of(fmt.Sprintf("%d", i))
	})(option.None[int]())
	if option.IsSome(noneResult) {
		t.Fatal("expected None for chain on None input")
	}
}

func TestMonadOptionCompiles(t *testing.T) {
	m := laws.MakeMonad[int, string, option.Option[int], option.Option[string], option.Option[func(int) string]](
		option.Of[int],
		option.Map[int, string],
		option.Ap[string, int],
		option.Chain[int, string],
	)

	// Test Of
	fa := m.Of(7)
	if !option.IsSome(fa) {
		t.Fatal("Of: expected Some(7), got None")
	}
	vi := option.GetOrElse(func() int { return 0 })(fa)
	if vi != 7 {
		t.Fatalf("Of: expected 7, got %d", vi)
	}

	// Test Map
	mapped := m.Map(func(i int) string { return fmt.Sprintf("%d", i) })(fa)
	vs := option.GetOrElse(func() string { return "" })(mapped)
	if vs != "7" {
		t.Fatalf("Map: expected \"7\", got %q", vs)
	}

	// Test Ap
	fn := option.Of(func(i int) string { return fmt.Sprintf("x=%d", i) })
	applied := m.Ap(option.Of(3))(fn)
	vs = option.GetOrElse(func() string { return "" })(applied)
	if vs != "x=3" {
		t.Fatalf("Ap: expected \"x=3\", got %q", vs)
	}

	// Test Chain
	chained := m.Chain(func(i int) option.Option[string] {
		return option.Of(fmt.Sprintf("chained:%d", i))
	})(option.Of(5))
	vs = option.GetOrElse(func() string { return "" })(chained)
	if vs != "chained:5" {
		t.Fatalf("Chain: expected \"chained:5\", got %q", vs)
	}
}

func TestApplyCompositionOption(t *testing.T) {
	laws.ApplyAssociativeComposition[
		option.Option[int],
		option.Option[string],
		option.Option[bool],
		option.Option[func(int) string],
		option.Option[func(string) bool],
		option.Option[func(int) bool],
		option.Option[func(func(int) string) func(int) bool],
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
