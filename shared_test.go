package detest

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type call struct {
	method string
	args   []interface{}
}

type mockT struct {
	calls []call
}

func (mt *mockT) called(args ...interface{}) {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	if n == 0 {
		panic("Cannot get caller from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	mt.calls = append(mt.calls, call{method: methodName(frame.Function), args: args})
}

func methodName(f string) string {
	s := strings.Split(f, ".")
	if len(s) == 0 {
		return ""
	}
	return s[len(s)-1]
}

func (mt *mockT) AssertNotCalled(t *testing.T, method string) {
	for _, c := range mt.calls {
		if c.method == method {
			t.Errorf("The %s method was called when it should not have been", method)
			return
		}
	}
}

func (mt *mockT) AssertCalled(t *testing.T, method string, args ...interface{}) {
	for _, c := range mt.calls {
		if c.method == method {
			if reflect.DeepEqual(c.args, args) {
				return
			}
		}
	}
	spew.Dump(mt.calls)
	t.Errorf("Expected the %s method to be called with %d args but it was not", method, len(args))
}

func (mt *mockT) Fail() {
	mt.called()
}

func (mt *mockT) WriteString(s string) (int, error) {
	mt.called(s)
	return len([]byte(s)), nil
}
