package laws

import (
	"testing"

	"pgregory.net/rapid"
)

// LensLaws verifies the Lens laws (get-set, set-get, set-set).
//
// A lawful lens must satisfy:
//   - GetSet: set(get(s), s) == s (setting what you get changes nothing)
//   - SetGet: get(set(a, s)) == a (you get back what you set)
//   - SetSet: set(b, set(a, s)) == set(b, s) (setting twice is setting once)
func LensLaws[S, A any](
	t *testing.T,
	genS *rapid.Generator[S],
	genA *rapid.Generator[A],
	eqS func(S, S) bool,
	eqA func(A, A) bool,
	get func(S) A,
	set func(A) func(S) S,
) {
	t.Helper()

	t.Run("Lens/GetSet", rapid.MakeCheck(func(t *rapid.T) {
		s := genS.Draw(t, "s")

		// Law: set(get(s))(s) == s
		got := set(get(s))(s)
		if !eqS(got, s) {
			t.Fatalf("Lens GetSet violated:\n  set(get(s))(s) = %v\n  s              = %v", got, s)
		}
	}))

	t.Run("Lens/SetGet", rapid.MakeCheck(func(t *rapid.T) {
		s := genS.Draw(t, "s")
		a := genA.Draw(t, "a")

		// Law: get(set(a)(s)) == a
		got := get(set(a)(s))
		if !eqA(got, a) {
			t.Fatalf("Lens SetGet violated:\n  get(set(a)(s)) = %v\n  a              = %v", got, a)
		}
	}))

	t.Run("Lens/SetSet", rapid.MakeCheck(func(t *rapid.T) {
		s := genS.Draw(t, "s")
		a := genA.Draw(t, "a")
		b := genA.Draw(t, "b")

		// Law: set(b)(set(a)(s)) == set(b)(s)
		left := set(b)(set(a)(s))
		right := set(b)(s)
		if !eqS(left, right) {
			t.Fatalf("Lens SetSet violated:\n  set(b)(set(a)(s)) = %v\n  set(b)(s)         = %v", left, right)
		}
	}))
}
