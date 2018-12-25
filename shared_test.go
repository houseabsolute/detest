package detest

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

type Call struct {
	Method string
	Args   []interface{}
}

type mockT struct {
	calls []Call
}

func (mt *mockT) called(args ...interface{}) {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	if n == 0 {
		panic("Cannot get caller from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	mt.calls = append(mt.calls, Call{Method: methodName(frame.Function), Args: args})
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
		if c.Method == method {
			t.Errorf("The %s method was called when it should not have been", method)
			return
		}
	}
}

func (mt *mockT) AssertCalled(t *testing.T, method string, args ...interface{}) {
	for _, c := range mt.calls {
		if c.Method == method {
			_, differences := mock.Arguments(args).Diff(c.Args)
			if differences == 0 {
				return
			}

			t.Errorf("Expected the %s method to be called with:\n%v\nbut it was called with:\n%v\n", method, args, c.Args)
			return
		}
	}
	t.Errorf("Expected the %s method to be called with:\n%v\nbut it was never called", method, args)
}

func (mt *mockT) FindCall(method string) *Call {
	for _, c := range mt.calls {
		if c.Method == method {
			c := c
			return &c
		}
	}
	return nil
}

func (mt *mockT) Fail() {
	mt.called()
}

func (mt *mockT) WriteString(s string) (int, error) {
	mt.called(s)
	return len([]byte(s)), nil
}

type GTComparer int

func (sc GTComparer) Compare(d *D) {
	d.AddResult(result{
		pass:   d.Actual().(int) > int(sc),
		actual: newValue(d.Actual()),
		expect: newValue(sc),
		op:     ">",
		where:  inValue,
	})
}
