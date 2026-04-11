package assert_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest-go/assert"
)

func TestAssertSome(t *testing.T) {
	got := assert.AssertSome(t, option.Some(42))
	if got != 42 {
		t.Fatalf("AssertSome returned %d, want 42", got)
	}
}

func TestAssertNone(t *testing.T) {
	assert.AssertNone(t, option.None[int]())
}

func TestAssertSomeEq(t *testing.T) {
	assert.AssertSomeEq(t, option.Some("hello"), "hello")
}

func TestAssertRight(t *testing.T) {
	got := assert.AssertRight[string](t, either.Right[string](42))
	if got != 42 {
		t.Fatalf("AssertRight returned %d, want 42", got)
	}
}

func TestAssertLeft(t *testing.T) {
	got := assert.AssertLeft[string, int](t, either.Left[int, string]("err"))
	if got != "err" {
		t.Fatalf("AssertLeft returned %q, want %q", got, "err")
	}
}

func TestAssertRightEq(t *testing.T) {
	assert.AssertRightEq[string](t, either.Right[string](42), 42)
}

func TestAssertLeftEq(t *testing.T) {
	assert.AssertLeftEq[string, int](t, either.Left[int, string]("err"), "err")
}

func TestAssertIO(t *testing.T) {
	got := assert.AssertIO(t, io.Of(42))
	if got != 42 {
		t.Fatalf("AssertIO returned %d, want 42", got)
	}
}

func TestAssertIORight(t *testing.T) {
	got := assert.AssertIORight[string](t, ioeither.Right[string](42))
	if got != 42 {
		t.Fatalf("AssertIORight returned %d, want 42", got)
	}
}

func TestAssertIOLeft(t *testing.T) {
	got := assert.AssertIOLeft[string, int](t, ioeither.Left[int, string]("err"))
	if got != "err" {
		t.Fatalf("AssertIOLeft returned %q, want %q", got, "err")
	}
}

func TestAssertIOSome(t *testing.T) {
	got := assert.AssertIOSome(t, io.Of(option.Some(42)))
	if got != 42 {
		t.Fatalf("AssertIOSome returned %d, want 42", got)
	}
}

func TestAssertIONone(t *testing.T) {
	assert.AssertIONone(t, io.Of(option.None[int]()))
}

func TestAssertEq(t *testing.T) {
	eqInt := eq.FromStrictEquals[int]()
	assert.AssertEq(t, eqInt, 42, 42)
}

func TestAssertEqFunc(t *testing.T) {
	assert.AssertEqFunc(t, func(a, b int) bool { return a == b }, 42, 42)
}

func TestAssertNotEq(t *testing.T) {
	eqInt := eq.FromStrictEquals[int]()
	assert.AssertNotEq(t, eqInt, 42, 43)
}

func TestChainedAssertions(t *testing.T) {
	// Demonstrate chained unwrapping: IOEither -> Either -> value
	action := ioeither.Right[string](option.Some(42))
	opt := assert.AssertIORight[string](t, action)
	val := assert.AssertSome(t, opt)
	if val != 42 {
		t.Fatalf("Chained assertion got %d, want 42", val)
	}
}
