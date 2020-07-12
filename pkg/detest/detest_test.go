package detest

import (
	"testing"

	"github.com/houseabsolute/detest/pkg/detest/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestNewPath(t *testing.T) {
	t.Run("simple NewPath call with 3 calls on stack", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.stack3",
				caller: "detest.stack2",
			},
			stack1(0),
		)
	})
	t.Run("simple NewPath call with 2 calls on stack", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.stack3",
				caller: "detest.stack2",
			},
			stack2(0),
		)
	})
	t.Run("simple NewPath call with 1 call on stack", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.stack3",
				caller: "detest.TestNewPath.func3",
			},
			stack3(0),
		)
	})
	t.Run("NewPath call with 3 calls on stack and 1 skip", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.stack2",
				caller: "detest.stack1",
			},
			stack1(1),
		)
	})
	t.Run("NewPath call with external package on stack", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.TestNewPath.func5.1",
				caller: "internal/testhelper/testhelper.go@9",
			},
			testhelper.Callback(func() interface{} {
				return New(new(mockT)).NewPath("foo", 0, "")
			}),
		)
	})
}

func TestRegisterPackage(t *testing.T) {
	ourPackages[testhelper.PackageName()] = true
	defer func() {
		delete(ourPackages, testhelper.PackageName())
	}()
	t.Run("NewPath call with registered package on stack", func(t *testing.T) {
		assert.Equal(
			t,
			Path{
				data:   "foo",
				callee: "detest.TestRegisterPackage.func2.1",
				caller: "testhelper.Callback",
			},
			testhelper.Callback(func() interface{} {
				return New(new(mockT)).NewPath("foo", 0, "")
			}),
		)
	})
}

func stack1(n int) Path {
	return stack2(n)
}

func stack2(n int) Path {
	return stack3(n)
}

func stack3(n int) Path {
	return New(new(mockT)).NewPath("foo", n, "")
}
