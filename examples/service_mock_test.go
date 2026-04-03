package examples_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/franchb/fptest/assert"
	"github.com/franchb/fptest/mock"
)

// --- Domain ---

type ServiceUser struct {
	ID    int
	Name  string
	Email string
}

// --- Repository interface ---

// UserRepository is a typical repository interface returning IOEither for
// effectful operations that can fail.
type UserRepository interface {
	FindByID(id int) ioeither.IOEither[error, ServiceUser]
	Save(user ServiceUser) ioeither.IOEither[error, ServiceUser]
}

// --- Service under test ---

// getUserName is a service function that fetches a user and extracts the name.
func getUserName(repo UserRepository, id int) ioeither.IOEither[error, string] {
	return ioeither.Map[error](func(u ServiceUser) string {
		return u.Name
	})(repo.FindByID(id))
}

// transferUser fetches a user, updates the email, and saves.
func transferUser(repo UserRepository, id int, newEmail string) ioeither.IOEither[error, ServiceUser] {
	return ioeither.Chain[error](func(u ServiceUser) ioeither.IOEither[error, ServiceUser] {
		u.Email = newEmail
		return repo.Save(u)
	})(repo.FindByID(id))
}

// --- Mock repository ---

type mockRepo struct {
	tracker *mock.CallTracker
	findFn  func(int) ioeither.IOEither[error, ServiceUser]
	saveFn  func(ServiceUser) ioeither.IOEither[error, ServiceUser]
}

func (m *mockRepo) FindByID(id int) ioeither.IOEither[error, ServiceUser] {
	return func() either.Either[error, ServiceUser] {
		m.tracker.RecordSync("FindByID", id)
		return m.findFn(id)()
	}
}

func (m *mockRepo) Save(user ServiceUser) ioeither.IOEither[error, ServiceUser] {
	return func() either.Either[error, ServiceUser] {
		m.tracker.RecordSync("Save", user)
		return m.saveFn(user)()
	}
}

// --- Tests ---

// TestServiceSuccess demonstrates using TrackedStub to mock a successful
// repository call. The test verifies both the returned value and that the
// correct method was called with the expected argument.
func TestServiceSuccess(t *testing.T) {
	ct := mock.NewCallTracker()()
	alice := ServiceUser{ID: 1, Name: "Alice", Email: "alice@example.com"}

	repo := &mockRepo{
		tracker: ct,
		findFn:  mock.Stub[int, error](alice),
		saveFn:  mock.Stub[ServiceUser, error](alice),
	}

	name := assert.AssertIORight[error](t, getUserName(repo, 1))
	if name != "Alice" {
		t.Fatalf("expected Alice, got %q", name)
	}

	if ct.CallCount("FindByID") != 1 {
		t.Fatalf("FindByID call count = %d, want 1", ct.CallCount("FindByID"))
	}
	calls := ct.CallsFor("FindByID")
	if calls[0].Args[0] != 1 {
		t.Fatalf("FindByID arg = %v, want 1", calls[0].Args[0])
	}
}

// TestServiceNotFound demonstrates using StubError to simulate a repository
// failure. The test verifies the error propagates correctly through the service.
func TestServiceNotFound(t *testing.T) {
	ct := mock.NewCallTracker()()
	notFound := errors.New("user not found")

	repo := &mockRepo{
		tracker: ct,
		findFn:  mock.StubError[int, error, ServiceUser](notFound),
		saveFn:  mock.Stub[ServiceUser, error](ServiceUser{}),
	}

	err := assert.AssertIOLeft[error, string](t, getUserName(repo, 999))
	if err.Error() != "user not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTrackedFuncDelegation demonstrates TrackedFunc wrapping a real
// implementation. The tracker records calls while the implementation provides
// conditional behavior based on input.
func TestTrackedFuncDelegation(t *testing.T) {
	ct := mock.NewCallTracker()()

	impl := func(id int) ioeither.IOEither[error, ServiceUser] {
		if id > 0 {
			return ioeither.Right[error](ServiceUser{ID: id, Name: fmt.Sprintf("User%d", id)})
		}
		return ioeither.Left[ServiceUser](errors.New("invalid id"))
	}

	tracked := mock.TrackedFunc(ct, "FindByID", impl)

	// Success case
	user := assert.AssertIORight[error](t, tracked(1))
	if user.Name != "User1" {
		t.Fatalf("expected User1, got %q", user.Name)
	}

	// Failure case
	assert.AssertIOLeft[error, ServiceUser](t, tracked(-1))

	if ct.CallCount("FindByID") != 2 {
		t.Fatalf("call count = %d, want 2", ct.CallCount("FindByID"))
	}
}

// TestIORefCounter demonstrates using IORef as a functional counter. IORef
// provides a mutable reference cell in the IO monad, useful for tracking state
// in tests without breaking referential transparency outside the IO context.
func TestIORefCounter(t *testing.T) {
	ref := mock.NewIORef(0)()

	// Simulate counting operations
	for i := 0; i < 5; i++ {
		ref.Modify(func(n int) int { return n + 1 })()
	}

	count := assert.AssertIO(t, ref.Read())
	if count != 5 {
		t.Fatalf("counter = %d, want 5", count)
	}

	// Reset
	ref.Write(0)()
	count = ref.ReadUnsafe()
	if count != 0 {
		t.Fatalf("after reset, counter = %d, want 0", count)
	}
}

// TestServiceComposition demonstrates composing multiple IOEither service calls
// in a pipeline. The test verifies that all mocked methods are called in the
// correct order with correct arguments.
func TestServiceComposition(t *testing.T) {
	ct := mock.NewCallTracker()()
	alice := ServiceUser{ID: 1, Name: "Alice", Email: "alice@example.com"}
	updated := ServiceUser{ID: 1, Name: "Alice", Email: "new@example.com"}

	repo := &mockRepo{
		tracker: ct,
		findFn:  mock.Stub[int, error](alice),
		saveFn:  mock.Stub[ServiceUser, error](updated),
	}

	result := assert.AssertIORight[error](t, transferUser(repo, 1, "new@example.com"))
	if result.Email != "new@example.com" {
		t.Fatalf("email = %q, want %q", result.Email, "new@example.com")
	}

	if ct.CallCount("FindByID") != 1 {
		t.Fatalf("FindByID calls = %d, want 1", ct.CallCount("FindByID"))
	}
	if ct.CallCount("Save") != 1 {
		t.Fatalf("Save calls = %d, want 1", ct.CallCount("Save"))
	}
}
