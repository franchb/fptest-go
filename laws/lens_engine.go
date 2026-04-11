package laws

import (
	"testing"

	"github.com/franchb/fptest-go/engine"
)

// LensLawsEngine verifies the Lens laws (get-set, set-get, set-set) using the engine abstraction.
//
// A lawful lens must satisfy:
//   - GetSet: set(get(s), s) == s (setting what you get changes nothing)
//   - SetGet: get(set(a, s)) == a (you get back what you set)
//   - SetSet: set(b, set(a, s)) == set(b, s) (setting twice is setting once)
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

		// Law: set(get(s))(s) == s
		got := set(get(s))(s)
		if !eqS(got, s) {
			et.Fatalf("Lens GetSet violated:\n  set(get(s))(s) = %v\n  s              = %v", got, s)
		}
	})

	runner.MakeCheck(t, "Lens/SetGet", func(et engine.T) {
		s := genS.Draw(et, "s")
		a := genA.Draw(et, "a")

		// Law: get(set(a)(s)) == a
		got := get(set(a)(s))
		if !eqA(got, a) {
			et.Fatalf("Lens SetGet violated:\n  get(set(a)(s)) = %v\n  a              = %v", got, a)
		}
	})

	runner.MakeCheck(t, "Lens/SetSet", func(et engine.T) {
		s := genS.Draw(et, "s")
		a := genA.Draw(et, "a")
		b := genA.Draw(et, "b")

		// Law: set(b)(set(a)(s)) == set(b)(s)
		left := set(b)(set(a)(s))
		right := set(b)(s)
		if !eqS(left, right) {
			et.Fatalf("Lens SetSet violated:\n  set(b)(set(a)(s)) = %v\n  set(b)(s)         = %v", left, right)
		}
	})
}
