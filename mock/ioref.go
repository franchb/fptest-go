// Package mock provides FP-style mock construction utilities for testing
// Reader-based dependency injection patterns.
package mock

import (
	"sync"

	"github.com/IBM/fp-go/v2/io"
)

// IORef is a mutable reference cell within the IO monad.
// It provides thread-safe read/write access to a value, suitable for
// tracking side effects in tests without breaking referential transparency
// at the type level.
type IORef[A any] struct {
	mu  sync.Mutex
	val A
}

// NewIORef creates an IO action that allocates a new IORef with the given initial value.
func NewIORef[A any](initial A) io.IO[*IORef[A]] {
	return func() *IORef[A] {
		return &IORef[A]{val: initial}
	}
}

// Read returns an IO action that reads the current value of the ref.
func (r *IORef[A]) Read() io.IO[A] {
	return func() A {
		r.mu.Lock()
		defer r.mu.Unlock()
		return r.val
	}
}

// Write returns an IO action that sets the value of the ref.
func (r *IORef[A]) Write(a A) io.IO[struct{}] {
	return func() struct{} {
		r.mu.Lock()
		r.val = a
		r.mu.Unlock()
		return struct{}{}
	}
}

// Modify returns an IO action that applies a function to the current value.
func (r *IORef[A]) Modify(f func(A) A) io.IO[struct{}] {
	return func() struct{} {
		r.mu.Lock()
		r.val = f(r.val)
		r.mu.Unlock()
		return struct{}{}
	}
}

// ReadUnsafe reads the current value directly (not wrapped in IO).
// Use only in test assertions where IO composition is not needed.
func (r *IORef[A]) ReadUnsafe() A {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.val
}
