package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/option"
)

// AssertIO executes an IO action and returns its result.
func AssertIO[A any](t testing.TB, action io.IO[A]) A {
	t.Helper()
	return action()
}

// AssertIORight executes an IOEither and extracts the Right value, failing if Left.
func AssertIORight[E, A any](t testing.TB, action ioeither.IOEither[E, A]) A {
	t.Helper()
	return AssertRight[E, A](t, action())
}

// AssertIOLeft executes an IOEither and extracts the Left value, failing if Right.
func AssertIOLeft[E, A any](t testing.TB, action ioeither.IOEither[E, A]) E {
	t.Helper()
	return AssertLeft[E, A](t, action())
}

// AssertIOSome executes an IO[Option[A]] and extracts the Some value, failing if None.
func AssertIOSome[A any](t testing.TB, action io.IO[option.Option[A]]) A {
	t.Helper()
	return AssertSome(t, action())
}

// AssertIONone executes an IO[Option[A]] and asserts the result is None.
func AssertIONone[A any](t testing.TB, action io.IO[option.Option[A]]) {
	t.Helper()
	AssertNone(t, action())
}

// AssertIOEq executes an IO action and checks the result equals the expected value.
func AssertIOEq[A comparable](t testing.TB, action io.IO[A], want A) {
	t.Helper()
	got := action()
	if got != want {
		t.Fatalf("IO value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}

// AssertIOEitherEq executes an IOEither and asserts the result equals the expected Either.
func AssertIOEitherEq[E, A comparable](t testing.TB, action ioeither.IOEither[E, A], want either.Either[E, A], eqEither func(either.Either[E, A], either.Either[E, A]) bool) {
	t.Helper()
	got := action()
	if !eqEither(got, want) {
		t.Fatalf("IOEither value mismatch:\n  got  = %v\n  want = %v", got, want)
	}
}
