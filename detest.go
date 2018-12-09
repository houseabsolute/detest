// Package detest implements a DSL-ish interface for testing complicated Go
// data structure, as well as structured output on test failures.
package detest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

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
	callee string
	caller string
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

// TestingT is an interface wrapper around `*testing.T` for the portion of its
// API that we care about.
type TestingT interface {
	Fail()
}

// StringWriter is an interface used for writing strings.
type StringWriter interface {
	WriteString(string) (int, error)
}

// D contains state for the current set of tests. You should create a new `D`
// in every `Test*` function or subtest.
type D struct {
	t             TestingT
	callerPackage string
	state         *state
	output        StringWriter
}

var ourPackages = map[string]bool{}

func init() {
	ourPackages[packageFromFrame(findFrame(0))] = true
}

// RegisterPackage adds the caller's package to the list of "internal"
// packages for the purposes of presenting paths in test failure
// output. Specifically, when a function in a registered package is found as
// the caller for a path, detest will use the function name as the caller
// rather than showing the file and line where the call occurred.
func RegisterPackage() {
	ourPackages[packageFromFrame(findFrame(1))] = true
}

func findFrame(s int) runtime.Frame {
	pc := make([]uintptr, 1)
	n := runtime.Callers(s+1, pc)
	if n == 0 {
		panic("Cannot get New() from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return frame
}

var packageRE = regexp.MustCompile(`((?:[^/]+/)*[^\.]+)\.`)

func packageFromFrame(frame runtime.Frame) string {
	m := packageRE.FindStringSubmatch(frame.Function)
	if len(m) == 1 {
		return ""
	}
	return m[1]
}

// New takes any implementer of the `TestingT` interface and returns a new
// `*detest.D`. A `*D` created this way will send its output to `os.Stdout`.
func New(t TestingT) *D {
	return &D{
		t:             t,
		callerPackage: packageFromFrame(findFrame(1)),
		output:        os.Stdout,
	}
}

// NewWithOutput takes any implementer of the `TestingT` interface and a
// `StringWriter` implementer and returns a new `*detest.D`. This is provided
// primarily for the benefit of testing code that wants to capture the output
// from detest.
func NewWithOutput(t TestingT, o StringWriter) *D {
	return &D{t: t, output: o}
}

// ResetState resets the internal state of the `*detest.D` struct. This is
// public for the benefit of test packages that want to provide their own
// comparers or test functions like `detest.Is`.
func (d *D) ResetState() {
	d.state = &state{}
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

var funcNameRE = regexp.MustCompile(`^.+/`)

// NewPath takes a data path element, the number of frames to skip, and an
// optional function name. It returns a new `Path` struct. If the function
// name is given, then this is used as the called function rather than looking
// at the call frames .
//
// When the desired frame is from a package marked as internal to detest, then
// the caller's line and file is replaced with a function name so that we
// don't show (unhelpful) information about the detest internals when
// displaying the path.
func (d *D) NewPath(data string, skip int, function string) Path {
	pc := make([]uintptr, 2)
	// The hard-coded "2" is here because we want to skip this frame and the
	// frame of the caller. We're interested in the frames before that.
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return Path{data: data}
	}

	frames := runtime.CallersFrames(pc)
	frame, more := frames.Next()

	var callee = calleeFromFrame(frame, function)

	if !more {
		return Path{
			data:   data,
			callee: funcNameRE.ReplaceAllLiteralString(callee, ""),
		}
	}

	frame, _ = frames.Next()

	return Path{
		data:   data,
		callee: funcNameRE.ReplaceAllLiteralString(callee, ""),
		caller: d.callerFromFrame(frame),
	}
}

func calleeFromFrame(frame runtime.Frame, function string) string {
	if function != "" {
		return function
	}

	callee := frame.Function
	if callee == "" {
		callee = "<unknown>"
	}

	return callee
}

func (d *D) callerFromFrame(frame runtime.Frame) string {
	if ourPackages[packageFromFrame(frame)] {
		return funcNameRE.ReplaceAllLiteralString(frame.Function, "")
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	file := frame.File
	for _, p := range strings.Split(gopath, ":") {
		src := filepath.Join(p, "src")
		if strings.HasPrefix(file, src) {
			// We want to remove the remaining separator character that's left
			// after the trim.
			file = strings.TrimPrefix(frame.File, src)[1:]
			break
		}
	}
	// If the caller is in the package that created our *D then we can strip
	// that from the caller path and just show a path relative to the package
	// root.
	if strings.HasPrefix(file, d.callerPackage) {
		file = strings.TrimPrefix(file, d.callerPackage)[1:]
	}

	return fmt.Sprintf("%s@%d", file, frame.Line)
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
			_, err = d.output.WriteString(fmt.Sprintf("Passed test: %s\n", name))
		} else {
			pass = false
			d.t.Fail()
			_, err = d.output.WriteString(r.describe(name, ansi.DefaultScheme))
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
	return fmt.Sprintf("%s called %s", p.caller, p.callee)
}
