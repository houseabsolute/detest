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

type path struct {
	data   string
	caller string
	at     string
}

type state struct {
	results []result
	actual  []interface{}
	path    []path
}

type Comparer interface {
	Compare(*D)
}

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

func New(t *testing.T) *D {
	return &D{t: t}
}

func (d *D) ResetState(actual interface{}) {
	d.state = &state{}
	d.PushActual(actual)
}

func (d *D) PushActual(actual interface{}) {
	d.state.actual = append(d.state.actual, actual)
}

func (d *D) PopActual() {
	if len(d.state.actual) > 0 {
		d.state.actual = d.state.actual[:len(d.state.actual)-1]
	}
}

var callerRE = regexp.MustCompile(`^.+/`)

func (d *D) makePath(data string, skip int, function string) path {
	pc := make([]uintptr, 2)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return path{data: data}
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
		return path{
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

	return path{
		data:   data,
		caller: callerRE.ReplaceAllLiteralString(caller, ""),
		at:     at,
	}
}

func (d *D) PushPath(path path) {
	d.state.path = append(d.state.path, path)
}

func (d *D) PopPath() {
	if len(d.state.path) > 0 {
		d.state.path = d.state.path[:len(d.state.path)-1]
	}
}

func (d *D) Actual() interface{} {
	if len(d.state.actual) == 0 {
		panic("Actual() called before any actual values are stored in the state")
	}
	return d.state.actual[len(d.state.actual)-1]
}

func (d *D) AddResult(r result) {
	for _, p := range d.state.path {
		r.path = append(r.path, p)
	}
	d.state.results = append(d.state.results, r)
}

func (d *D) ok(name string) bool {
	pass := true
	for _, r := range d.state.results {
		if r.pass {
			os.Stdout.WriteString(fmt.Sprintf("Passed test: %s\n", name))
		} else {
			pass = false
			d.t.Fail()
			os.Stdout.WriteString(r.describe(name, ansi.DefaultScheme))
		}
	}

	return pass
}

var vowelRE = regexp.MustCompile(`^[aeiou]`)

func articleize(noun string) string {
	if vowelRE.MatchString(noun) {
		return "an " + noun
	}
	return "a " + noun
}
