package gen_test

import (
	"testing"

	"github.com/franchb/fptest/gen"
	"pgregory.net/rapid"
)

func TestOf(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.Int().Draw(t, "n")
		g := gen.Of(n)
		got := g(t)
		if got != n {
			t.Fatalf("Of(%d) produced %d", n, got)
		}
	})
}

func TestMap(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.Int().Draw(t, "n")
		g := gen.Map(gen.Of(n), func(x int) int { return x * 2 })
		got := g(t)
		if got != n*2 {
			t.Fatalf("Map(Of(%d), *2) = %d, want %d", n, got, n*2)
		}
	})
}

func TestChain(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(0, 100).Draw(t, "n")
		g := gen.Chain(gen.Of(n), func(x int) gen.Gen[int] {
			return gen.Of(x + 1)
		})
		got := g(t)
		if got != n+1 {
			t.Fatalf("Chain(Of(%d), +1) = %d, want %d", n, got, n+1)
		}
	})
}

func TestAp(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.Int().Draw(t, "n")
		gf := gen.Of(func(x int) string {
			if x > 0 {
				return "positive"
			}
			return "non-positive"
		})
		ga := gen.Of(n)
		g := gen.Ap(gf, ga)
		got := g(t)
		want := "non-positive"
		if n > 0 {
			want = "positive"
		}
		if got != want {
			t.Fatalf("Ap produced %q, want %q for n=%d", got, want, n)
		}
	})
}

func TestToRapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// ToRapid requires the Gen to use at least one rapid generator
		g := gen.Map(gen.FromRapid(rapid.IntRange(1, 10)), func(x int) int { return x * 10 })
		rg := gen.ToRapid(g)
		got := rg.Draw(t, "val")
		if got < 10 || got > 100 {
			t.Fatalf("ToRapid drew %d, want 10..100", got)
		}
	})
}

func TestFromRapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		rg := rapid.IntRange(1, 10)
		g := gen.FromRapid(rg)
		got := g(t)
		if got < 1 || got > 10 {
			t.Fatalf("FromRapid(IntRange(1,10)) produced %d", got)
		}
	})
}

func TestMap2(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := rapid.Int().Draw(t, "a")
		b := rapid.Int().Draw(t, "b")
		g := gen.Map2(gen.Of(a), gen.Of(b), func(x, y int) int { return x + y })
		got := g(t)
		if got != a+b {
			t.Fatalf("Map2(Of(%d), Of(%d), +) = %d, want %d", a, b, got, a+b)
		}
	})
}

func TestSlice(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		g := gen.Slice(gen.Of(1), 2, 5)
		got := g(t)
		if len(got) < 2 || len(got) > 5 {
			t.Fatalf("Slice len = %d, want 2..5", len(got))
		}
		for i, v := range got {
			if v != 1 {
				t.Fatalf("Slice[%d] = %d, want 1", i, v)
			}
		}
	})
}

func TestFilter(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		g := gen.Filter(gen.FromRapid(rapid.IntRange(-100, 100)), func(x int) bool {
			return x > 0
		})
		got := g(t)
		if got <= 0 {
			t.Fatalf("Filter(>0) produced %d", got)
		}
	})
}
