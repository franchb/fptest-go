package mock

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/ioeither"
)

// Stub creates a stub IOEither function that always returns Right(value).
func Stub[A, E, R any](value R) func(A) ioeither.IOEither[E, R] {
	return func(_ A) ioeither.IOEither[E, R] {
		return ioeither.Right[E](value)
	}
}

// StubError creates a stub IOEither function that always returns Left(err).
func StubError[A, E, R any](err E) func(A) ioeither.IOEither[E, R] {
	return func(_ A) ioeither.IOEither[E, R] {
		return ioeither.Left[R](err)
	}
}

// TrackedStub creates a stub IOEither that records calls to the tracker.
func TrackedStub[A, E, R any](ct *CallTracker, method string, value R) func(A) ioeither.IOEither[E, R] {
	return func(a A) ioeither.IOEither[E, R] {
		return func() either.Either[E, R] {
			ct.RecordSync(method, a)
			return either.Right[E](value)
		}
	}
}

// TrackedStubError creates a failing stub IOEither that records calls to the tracker.
func TrackedStubError[A, E, R any](ct *CallTracker, method string, err E) func(A) ioeither.IOEither[E, R] {
	return func(a A) ioeither.IOEither[E, R] {
		return func() either.Either[E, R] {
			ct.RecordSync(method, a)
			return either.Left[R](err)
		}
	}
}

// TrackedFunc creates a tracked IOEither function that delegates to the provided implementation.
func TrackedFunc[A, E, R any](ct *CallTracker, method string, impl func(A) ioeither.IOEither[E, R]) func(A) ioeither.IOEither[E, R] {
	return func(a A) ioeither.IOEither[E, R] {
		return func() either.Either[E, R] {
			ct.RecordSync(method, a)
			return impl(a)()
		}
	}
}
