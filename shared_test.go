package detest

import (
	"runtime"
	"strings"
	"testing"

	"github.com/autarch/testify/mock"
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
			_, differences := mock.Arguments(args).Diff(c.args)
			if differences == 0 {
				return
			}
		}
	}
	t.Errorf("Expected the %s method to be called with:\n%v\nbut it was not", method, args)
}

func (mt *mockT) Fail() {
	mt.called()
}

func (mt *mockT) WriteString(s string) (int, error) {
	mt.called(s)
	return len([]byte(s)), nil
}
