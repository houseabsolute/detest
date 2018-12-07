// Package detest implements a DSL-ish interface for testing complicated Go
// data structure, as well as structured output on test failures.
package detest

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/houseabsolute/detest/internal/ansi"
)

type failure int

const (
	inType failure = iota
	inValue
	inDataStructure
	inUsage
)

// Path is used to track the data path as a test goes through a complex data
// structure. It records a place in a data structure along with information
// about the call stack at that particular point in the data path.
type Path struct {
	data   string
	caller string
	at     string
}

type state struct {
	results []result
	actual  []interface{}
	path    []Path
}

// Comparer is the interface for anything that implements the `Compare`
// method.
type Comparer interface {
	// Compare is called with one argument, the current `*detest.D`
	// object. You can call `d.Actual()` to get the variable to be tested.
	Compare(*D)
}

// D contains state for the current set of tests. You should create a new `D`
// in every `Test*` function or subtest.
type D struct {
	t     *testing.T
	state *state
}

var ourPackages = map[string]bool{}

func init() {
	pc := make([]uintptr, 1)
	n := runtime.Callers(1, pc)
	if n == 0 {
		panic("Cannot get New() from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	pkg := packageFromFrame(frame)

	ourPackages[pkg] = true
}

func packageFromFrame(frame runtime.Frame) string {
	s := strings.Split(frame.Function, ".")
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

// New takes a `*testing.T` and returns a new `*detest.D`.
func New(t *testing.T) *D {
	return &D{t: t}
}

// ResetState resets the internal state of the `*detest.D` struct. This is
// public for the benefit of test packages that want to provide their own
// comparers or test functions like `detest.Is`.
func (d *D) ResetState(actual interface{}) {
	d.state = &state{}
	d.PushActual(actual)
}

// PushActual adds an actual variable being tested to the current stack of
// variables.
func (d *D) PushActual(actual interface{}) {
	d.state.actual = append(d.state.actual, actual)
}

// PopActual removes the top element from the current stack of variables being
// tested.
func (d *D) PopActual() {
	if len(d.state.actual) > 0 {
		d.state.actual = d.state.actual[:len(d.state.actual)-1]
	}
}

var callerRE = regexp.MustCompile(`^.+/`)

// NewPath takes a data path element, the number of frames to skip, and an
// optional function name. It returns a new `Path` struct. If the function
// name is given, then this is used as the caller rather than looking at the
// call frame's function.
//
// When the desired frame is from a package marked as internal to detest, then
// the caller's line and file is replaced with a function name so that we
// don't show (unhelpful) information about the detest internals when
// displaying the path.
func (d *D) NewPath(data string, skip int, function string) Path {
	pc := make([]uintptr, 2)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return Path{data: data}
	}

	frames := runtime.CallersFrames(pc)
	frame, more := frames.Next()

	var caller string
	if function == "" {
		caller = frame.Function
		if caller == "" {
			caller = "<unknown>"
		}
	} else {
		caller = function
	}

	if !more {
		return Path{
			data:   data,
			caller: callerRE.ReplaceAllLiteralString(caller, ""),
		}
	}

	frame, _ = frames.Next()

	var at string
	if ourPackages[packageFromFrame(frame)] {
		at = callerRE.ReplaceAllLiteralString(frame.Function, "")
	} else {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fileRE := regexp.MustCompile(`^` + wd + `/`)

		at = fmt.Sprintf("%s@%d", fileRE.ReplaceAllLiteralString(frame.File, ""), frame.Line)
	}

	return Path{
		data:   data,
		caller: callerRE.ReplaceAllLiteralString(caller, ""),
		at:     at,
	}
}

// PushPath adds a path to the current path stack.
func (d *D) PushPath(path Path) {
	d.state.path = append(d.state.path, path)
}

// PopPath removes the top path from the current path stack.
func (d *D) PopPath() {
	if len(d.state.path) > 0 {
		d.state.path = d.state.path[:len(d.state.path)-1]
	}
}

// Actual returns the top actual variable from the stack of variables being
// tested.
func (d *D) Actual() interface{} {
	if len(d.state.actual) == 0 {
		panic("Actual() called before any actual values are stored in the state")
	}
	return d.state.actual[len(d.state.actual)-1]
}

// AddResult adds a test result. At the end of a test any result which is
// marked as failing is displayed as its own table.
func (d *D) AddResult(r result) {
	// We want to make a new slice since d.state.path will could get pushed
	// and popped after this result is saved.
	r.path = append(r.path, d.state.path...)
	d.state.results = append(d.state.results, r)
}

func (d *D) ok(name string) bool {
	pass := true
	for _, r := range d.state.results {
		var err error
		if r.pass {
			_, err = os.Stdout.WriteString(fmt.Sprintf("Passed test: %s\n", name))
		} else {
			pass = false
			d.t.Fail()
			_, err = os.Stdout.WriteString(r.describe(name, ansi.DefaultScheme))
		}
		if err != nil {
			panic(err)
		}
	}

	return pass
}

// CalledAt returns a string describing the function, file, and line for this
// path element.
func (p Path) CalledAt() string {
	return fmt.Sprintf("%s called %s", p.at, p.caller)
}

var vowelRE = regexp.MustCompile(`^[aeiou]`)

func articleize(noun string) string {
	if vowelRE.MatchString(noun) {
		return "an " + noun
	}
	return "a " + noun
}
