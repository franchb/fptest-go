package mock

import (
	"github.com/IBM/fp-go/v2/io"
)

// Call represents a recorded method invocation.
type Call struct {
	Method string
	Args   []any
}

// CallTracker records method calls using an IORef, enabling functional call verification.
type CallTracker struct {
	ref *IORef[[]Call]
}

// NewCallTracker creates an IO action that allocates a new CallTracker.
func NewCallTracker() io.IO[*CallTracker] {
	return func() *CallTracker {
		ref := NewIORef([]Call{})(/* execute */)
		return &CallTracker{ref: ref}
	}
}

// Record returns an IO action that appends a call record.
func (ct *CallTracker) Record(method string, args ...any) io.IO[struct{}] {
	return ct.ref.Modify(func(calls []Call) []Call {
		return append(calls, Call{Method: method, Args: args})
	})
}

// RecordSync records a call synchronously (not wrapped in IO).
// Use in mock method implementations where IO composition is impractical.
func (ct *CallTracker) RecordSync(method string, args ...any) {
	ct.ref.Modify(func(calls []Call) []Call {
		return append(calls, Call{Method: method, Args: args})
	})()
}

// Calls returns an IO action that reads all recorded calls.
func (ct *CallTracker) Calls() io.IO[[]Call] {
	return ct.ref.Read()
}

// CallsUnsafe reads all recorded calls directly (not wrapped in IO).
func (ct *CallTracker) CallsUnsafe() []Call {
	return ct.ref.ReadUnsafe()
}

// CallCount returns the number of calls recorded for the given method.
func (ct *CallTracker) CallCount(method string) int {
	calls := ct.CallsUnsafe()
	count := 0
	for _, c := range calls {
		if c.Method == method {
			count++
		}
	}
	return count
}

// CallsFor returns all calls recorded for the given method.
func (ct *CallTracker) CallsFor(method string) []Call {
	calls := ct.CallsUnsafe()
	var result []Call
	for _, c := range calls {
		if c.Method == method {
			result = append(result, c)
		}
	}
	return result
}
