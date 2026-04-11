package hegel_test

import (
	"math"
	"testing"

	"github.com/franchb/fptest-go/laws"
	hegellaws "github.com/franchb/fptest-go/hegel/laws"
	hegellib "hegel.dev/go/hegel"
	"pgregory.net/rapid"
)

// TestMonoidLaws_CrossEngine runs the same Monoid law check under both engines.
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

// TestEqLaws_CrossEngine runs the same Eq law check under both engines.
func TestEqLaws_CrossEngine(t *testing.T) {
	eqString := func(a, b string) bool { return a == b }

	t.Run("rapid", func(t *testing.T) {
		laws.EqLaws(t, rapid.String(), eqString)
	})

	t.Run("hegel", func(t *testing.T) {
		hegellaws.EqLaws(t, hegellib.Text(0, 50), eqString)
	})
}

// TestSemigroupLaws_CrossEngine runs the same Semigroup law check under both engines.
func TestSemigroupLaws_CrossEngine(t *testing.T) {
	eqString := func(a, b string) bool { return a == b }
	concat := func(a, b string) string { return a + b }

	t.Run("rapid", func(t *testing.T) {
		laws.SemigroupLaws(t, rapid.String(), eqString, concat)
	})

	t.Run("hegel", func(t *testing.T) {
		hegellaws.SemigroupLaws(t, hegellib.Text(0, 50), eqString, concat)
	})
}
