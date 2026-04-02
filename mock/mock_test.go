package mock_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/franchb/fptest/assert"
	"github.com/franchb/fptest/mock"
)

func TestIORef(t *testing.T) {
	ref := mock.NewIORef(0)()

	// Read initial value
	val := ref.Read()()
	if val != 0 {
		t.Fatalf("initial value = %d, want 0", val)
	}

	// Write and read
	ref.Write(42)()
	val = ref.Read()()
	if val != 42 {
		t.Fatalf("after Write(42), Read = %d", val)
	}

	// Modify
	ref.Modify(func(x int) int { return x + 8 })()
	val = ref.Read()()
	if val != 50 {
		t.Fatalf("after Modify(+8), Read = %d", val)
	}

	// ReadUnsafe
	val = ref.ReadUnsafe()
	if val != 50 {
		t.Fatalf("ReadUnsafe = %d, want 50", val)
	}
}

func TestCallTracker(t *testing.T) {
	ct := mock.NewCallTracker()()

	ct.RecordSync("FindByID", 1)
	ct.RecordSync("FindByID", 2)
	ct.RecordSync("Save", "user")

	if ct.CallCount("FindByID") != 2 {
		t.Fatalf("FindByID call count = %d, want 2", ct.CallCount("FindByID"))
	}
	if ct.CallCount("Save") != 1 {
		t.Fatalf("Save call count = %d, want 1", ct.CallCount("Save"))
	}
	if ct.CallCount("Delete") != 0 {
		t.Fatalf("Delete call count = %d, want 0", ct.CallCount("Delete"))
	}

	finds := ct.CallsFor("FindByID")
	if len(finds) != 2 {
		t.Fatalf("CallsFor(FindByID) len = %d, want 2", len(finds))
	}
	if finds[0].Args[0] != 1 || finds[1].Args[0] != 2 {
		t.Fatalf("unexpected FindByID args: %v", finds)
	}
}

func TestCallTrackerIO(t *testing.T) {
	ct := mock.NewCallTracker()()

	// Record via IO action
	ct.Record("Method", "arg1", "arg2")()

	calls := ct.Calls()()
	if len(calls) != 1 {
		t.Fatalf("calls len = %d, want 1", len(calls))
	}
	if calls[0].Method != "Method" {
		t.Fatalf("call method = %q, want %q", calls[0].Method, "Method")
	}
}

func TestStub(t *testing.T) {
	stub := mock.Stub[int, string, bool](true)
	result := stub(42)()
	val := assert.AssertRight[string](t, result)
	if !val {
		t.Fatal("stub returned false, want true")
	}
}

func TestStubError(t *testing.T) {
	stub := mock.StubError[int, string, bool]("oops")
	result := stub(42)()
	err := assert.AssertLeft[string, bool](t, result)
	if err != "oops" {
		t.Fatalf("stub error = %q, want %q", err, "oops")
	}
}

func TestTrackedStub(t *testing.T) {
	ct := mock.NewCallTracker()()
	stub := mock.TrackedStub[int, string, bool](ct, "DoSomething", true)

	result := stub(42)()
	assert.AssertRight[string](t, result)

	if ct.CallCount("DoSomething") != 1 {
		t.Fatalf("call count = %d, want 1", ct.CallCount("DoSomething"))
	}
	calls := ct.CallsFor("DoSomething")
	if calls[0].Args[0] != 42 {
		t.Fatalf("call arg = %v, want 42", calls[0].Args[0])
	}
}

func TestTrackedFunc(t *testing.T) {
	ct := mock.NewCallTracker()()
	impl := func(id int) ioeither.IOEither[string, string] {
		if id > 0 {
			return ioeither.Right[string]("found")
		}
		return ioeither.Left[string, string]("not found")
	}
	tracked := mock.TrackedFunc(ct, "FindByID", impl)

	// Positive case
	result := tracked(1)()
	assert.AssertRightEq[string](t, result, "found")

	// Negative case
	result = tracked(-1)()
	err := assert.AssertLeft[string, string](t, result)
	if err != "not found" {
		t.Fatalf("err = %q, want %q", err, "not found")
	}

	if ct.CallCount("FindByID") != 2 {
		t.Fatalf("call count = %d, want 2", ct.CallCount("FindByID"))
	}
}

// Test Reader-based dependency injection pattern
type UserRepo interface {
	FindByID(id int) ioeither.IOEither[string, string]
}

type Deps struct {
	UserRepo UserRepo
}

type mockUserRepo struct {
	tracker *mock.CallTracker
	findFn  func(int) ioeither.IOEither[string, string]
}

func (m *mockUserRepo) FindByID(id int) ioeither.IOEither[string, string] {
	return func() either.Either[string, string] {
		m.tracker.RecordSync("FindByID", id)
		return m.findFn(id)()
	}
}

func TestReaderDIPattern(t *testing.T) {
	ct := mock.NewCallTracker()()
	repo := &mockUserRepo{
		tracker: ct,
		findFn: func(id int) ioeither.IOEither[string, string] {
			return ioeither.Right[string]("Alice")
		},
	}
	deps := Deps{UserRepo: repo}

	// Simulate a Reader-based workflow: func(Deps) IOEither[string, string]
	getUserName := func(d Deps) ioeither.IOEither[string, string] {
		return d.UserRepo.FindByID(1)
	}

	result := getUserName(deps)()
	name := assert.AssertRight[string](t, result)
	if name != "Alice" {
		t.Fatalf("name = %q, want %q", name, "Alice")
	}
	if ct.CallCount("FindByID") != 1 {
		t.Fatalf("call count = %d, want 1", ct.CallCount("FindByID"))
	}
}
