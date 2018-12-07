package detest

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/houseabsolute/detest/internal/table"
	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/style"
)

type failure int

const (
	inType failure = iota
	inValue
	inDataStructure
	inUsage
)

type value struct {
	value interface{}
}

type result struct {
	actual      *value
	expect      *value
	op          string
	pass        bool
	path        []path
	where       failure
	description string
}

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
	t           *testing.T
	state       *state
	ourPackages map[string]bool
}

func New(t *testing.T) *D {
	pc := make([]uintptr, 1)
	n := runtime.Callers(1, pc)
	if n == 0 {
		panic("Cannot get New() from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	pkg := packageFromFrame(frame)

	return &D{
		t:           t,
		ourPackages: map[string]bool{pkg: true},
	}
}

func packageFromFrame(frame runtime.Frame) string {
	s := strings.Split(frame.Function, ".")
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func (d *D) Is(actual, expect interface{}, name string) bool {
	d.ResetState(actual)
	defer d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(d)
	} else {
		d.Equal(expect).Compare(d)
	}
	return d.ok(name)
}

func (d *D) ValueIs(actual, expect interface{}, name string) bool {
	d.ResetState(actual)
	defer d.PopActual()

	d.ValueEqual(expect).Compare(d)
	return d.ok(name)
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
	if d.ourPackages[packageFromFrame(frame)] {
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
			os.Stdout.WriteString(r.describe(name))
		}
	}

	return pass
}

func (r result) describe(name string) string {
	var aType, eType string
	var showActual, showExpect bool

	if r.actual != nil {
		aType = typeOf(r.actual.value)
		if r.actual.value == nil {
			aType += " <nil>"
		}
		showActual = true
	}

	if r.expect != nil {
		eType = typeOf(r.expect.value)
		if r.expect.value == nil {
			eType += " <nil>"
		}
		showExpect = true
	}

	scheme := ansi.DefaultScheme

	t := table.NewWithTitle(scheme.Strong(fmt.Sprintf("Failed test: %s", name)))

	addHeaders(t, r)

	var actual, expect, op string
	if showActual {
		actual = fmt.Sprintf("%v", r.actual.value)
	}
	if showExpect {
		expect = fmt.Sprintf("%v", r.expect.value)
	}
	op = r.op

	if r.where == inType {
		aType = scheme.Incorrect(aType)
		eType = scheme.Correct(eType)
	} else if r.where == inValue {
		actual = scheme.Incorrect(actual)
		expect = scheme.Correct(expect)
	} else if r.where == inDataStructure {
		op = scheme.Incorrect(op)
	}

	lastBodyRow := []interface{}{}
	if len(r.path) != 0 {
		lastBodyRow = append(lastBodyRow, "")
	}
	if showActual {
		lastBodyRow = append(lastBodyRow, aType, actual)
	}
	if op != "" {
		lastBodyRow = append(lastBodyRow, op)
	}
	if showExpect {
		lastBodyRow = append(lastBodyRow, eType, expect)
	}
	if len(r.path) != 0 {
		lastBodyRow = append(lastBodyRow, "")
	}

	body := [][]interface{}{}
	for _, p := range r.path {
		body = append(
			body,
			[]interface{}{
				p.data,
				cell.NewWithParams("", len(lastBodyRow)-2, cell.AlignLeft),
				pathSummary(p),
			},
		)
	}
	body = append(body, lastBodyRow)

	for _, b := range body {
		t.AddRow(b...)
	}

	if r.description != "" {
		span := 0
		if len(r.path) != 0 {
			span += 2
		}
		if showActual {
			span += 2
		}
		if op != "" {
			span += 1
		}
		if showExpect {
			span += 2
		}
		t.AddFooterRow(
			cell.NewWithParams(scheme.Strong(scheme.Incorrect(r.description)), span, cell.AlignLeft),
		)
	}

	rendered, err := t.Render(style.Default)
	if err != nil {
		panic(err)
	}
	return rendered
}

func addHeaders(t *table.Table, r result) {
	first := []interface{}{}
	if len(r.path) != 0 {
		first = append(first, "")
	}
	if r.actual != nil {
		first = append(
			first,
			cell.NewWithParams("ACTUAL", 2, cell.AlignCenter),
		)
		if r.op != "" {
			first = append(first, cell.NewWithParams("", 1, cell.AlignCenter))
		}
	}

	if r.expect != nil {
		first = append(first, cell.NewWithParams("EXPECT", 2, cell.AlignCenter))
	}
	if len(r.path) != 0 {
		first = append(first, "")
	}

	t.AddHeaderRow(first...)

	second := []interface{}{}
	if len(r.path) != 0 {
		second = append(second, cell.NewWithParams("PATH", 1, cell.AlignCenter))
	}

	if r.actual != nil {
		second = append(
			second,
			cell.NewWithParams("TYPE", 1, cell.AlignCenter),
			cell.NewWithParams("VALUE", 1, cell.AlignCenter),
		)
	}
	if r.op != "" {
		second = append(second, cell.NewWithParams("OP", 1, cell.AlignCenter))
	}

	if r.expect != nil {
		second = append(
			second,
			cell.NewWithParams("TYPE", 1, cell.AlignCenter),
			cell.NewWithParams("VALUE", 1, cell.AlignCenter),
		)
	}
	if len(r.path) != 0 {
		second = append(second, cell.NewWithParams("CALLER", 1, cell.AlignCenter))
	}

	t.AddHeaderRow(second...)

}

func pathSummary(p path) string {
	return fmt.Sprintf("%s called %s", p.at, p.caller)
}

func typeOf(val interface{}) string {
	return describeType(reflect.ValueOf(val).Type())
}

func describeType(ty reflect.Type) string {
	k := ty.Kind().String()
	// This is only true for simple types like string, float64, etc. If it's
	// not composite or it's not a built-in then the name doesn't match the
	// kind.
	if k == ty.Name() {
		return k
	}

	switch ty.Kind() {
	case reflect.Array:
		return fmt.Sprintf("[%d]", ty.Len()) + describeType(ty.Elem())
	case reflect.Chan:
		return fmt.Sprintf("chan(%s)", describeType(ty.Elem()))
	case reflect.Func:
		return describeFunc(ty)
	case reflect.Interface:
		// Can this happen?
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", describeType(ty.Key()), describeType(ty.Elem()))
	case reflect.Ptr:
		return "*" + describeType(ty.Elem())
	case reflect.Slice:
		return "[]" + describeType(ty.Elem())
	case reflect.Struct:
		return describeStruct(ty)
	case reflect.UnsafePointer:
		return "*<unsafe>"
	}

	// wtf - should not get here
	return ""
}

func describeFunc(ty reflect.Type) string {
	desc := "func "
	if name := ty.Name(); name != "" {
		desc = desc + name + " "
	}

	desc = desc + "("
	for i := 0; i < ty.NumIn(); i++ {
		desc = desc + describeType(ty.In(i))
	}
	if ty.IsVariadic() {
		desc = desc + "..."
	}
	desc = desc + ") "

	if ty.NumOut() > 1 {
		desc = desc + "("
	}
	for i := 0; i < ty.NumOut(); i++ {
		desc = desc + describeType(ty.Out(i))
	}
	if ty.NumOut() > 1 {
		desc = desc + ")"
	}

	return desc
}

func describeStruct(ty reflect.Type) string {
	if ty.Name() != "" {
		return ty.Name()
	}

	return "<anon struct>"
}

var vowelRE = regexp.MustCompile(`^[aeiou]`)

func articleize(noun string) string {
	if vowelRE.MatchString(noun) {
		return "an " + noun
	}
	return "a " + noun
}
