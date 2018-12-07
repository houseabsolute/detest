package detest

import (
	"fmt"
	"reflect"

	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/houseabsolute/detest/internal/table"
	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/style"
)

// We wrap values in a struct so that we can use a nil *value to indicate that
// the value wasn't present, as opposed to a nil *value.value, which is a nil
// _value_.
type value struct {
	value interface{}
	desc  string
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

func newValue(val interface{}) *value {
	return &value{value: val}
}

func (r result) hasPath() bool {
	return len(r.path) != 0
}

func (r result) showActual() bool {
	return r.actual != nil
}

func (r result) showExpect() bool {
	return r.expect != nil
}

type describer struct {
	r result
	t *table.Table
	s ansi.Scheme
}

func (r result) describe(name string, s ansi.Scheme) string {
	t := table.NewWithTitle(s.Strong(fmt.Sprintf("Failed test: %s", name)))
	return describer{r, t, s}.table()
}

func (d describer) table() string {
	d.addHeaders()

	lastBodyRow := d.lastBodyRow()

	body := [][]interface{}{}
	for _, p := range d.r.path {
		body = append(
			body,
			[]interface{}{
				p.data,
				cell.NewWithParams("", len(lastBodyRow)-2, cell.AlignLeft),
				d.pathSummary(p),
			},
		)
	}
	body = append(body, lastBodyRow)
	for _, b := range body {
		d.t.AddRow(b...)
	}

	if d.r.description != "" {
		span := 0
		if d.r.hasPath() {
			span += 2
		}
		if d.r.showActual() {
			span += 2
		}
		if d.r.op != "" {
			span += 1
		}
		if d.r.showExpect() {
			span += 2
		}
		d.t.AddFooterRow(
			cell.NewWithParams(d.s.Strong(d.s.Incorrect(d.r.description)), span, cell.AlignLeft),
		)
	}

	rendered, err := d.t.Render(style.Default)
	if err != nil {
		panic(err)
	}
	return rendered
}

func (d describer) addHeaders() {
	first := []interface{}{}
	if d.r.hasPath() {
		first = append(first, "")
	}
	if d.r.showActual() {
		first = append(
			first,
			cell.NewWithParams("ACTUAL", 2, cell.AlignCenter),
		)
		if d.r.op != "" {
			first = append(first, cell.NewWithParams("", 1, cell.AlignCenter))
		}
	}

	if d.r.showExpect() {
		first = append(first, cell.NewWithParams("EXPECT", 2, cell.AlignCenter))
	}
	if d.r.hasPath() {
		first = append(first, "")
	}

	d.t.AddHeaderRow(first...)

	second := []interface{}{}
	if d.r.hasPath() {
		second = append(second, cell.NewWithParams("PATH", 1, cell.AlignCenter))
	}

	if d.r.showActual() {
		second = append(
			second,
			cell.NewWithParams("TYPE", 1, cell.AlignCenter),
			cell.NewWithParams("VALUE", 1, cell.AlignCenter),
		)
	}
	if d.r.op != "" {
		second = append(second, cell.NewWithParams("OP", 1, cell.AlignCenter))
	}

	if d.r.showExpect() {
		second = append(
			second,
			cell.NewWithParams("TYPE", 1, cell.AlignCenter),
			cell.NewWithParams("VALUE", 1, cell.AlignCenter),
		)
	}
	if d.r.hasPath() {
		second = append(second, cell.NewWithParams("CALLER", 1, cell.AlignCenter))
	}

	d.t.AddHeaderRow(second...)
}

func (d describer) lastBodyRow() []interface{} {
	var actual, expect, op string
	if d.r.showActual() {
		actual = fmt.Sprintf("%v", d.r.actual.value)
	}
	if d.r.showExpect() {
		expect = fmt.Sprintf("%v", d.r.expect.value)
	}
	op = d.r.op

	var aType, eType string
	if d.r.showActual() {
		aType = d.r.actual.description()
	}
	if d.r.showExpect() {
		eType = d.r.expect.description()
	}
	if d.r.where == inType {
		aType = d.s.Incorrect(aType)
		eType = d.s.Correct(eType)
	} else if d.r.where == inValue {
		actual = d.s.Incorrect(actual)
		expect = d.s.Correct(expect)
	} else if d.r.where == inDataStructure {
		op = d.s.Incorrect(op)
	}

	lastBodyRow := []interface{}{}
	if d.r.hasPath() {
		lastBodyRow = append(lastBodyRow, "")
	}
	if d.r.showActual() {
		aType := d.r.actual.description()
		lastBodyRow = append(lastBodyRow, aType, actual)
	}
	if op != "" {
		lastBodyRow = append(lastBodyRow, op)
	}
	if d.r.showExpect() {
		lastBodyRow = append(lastBodyRow, eType, expect)
	}
	if d.r.hasPath() {
		lastBodyRow = append(lastBodyRow, "")
	}

	return lastBodyRow
}

func (d describer) pathSummary(p path) string {
	return fmt.Sprintf("%s called %s", p.at, p.caller)
}

func (v *value) description() string {
	if v.desc != "" {
		return v.desc
	}

	v.desc = describeTypeOfValue(v.value)
	if v.value == nil {
		v.desc += " <nil>"
	}
	return v.desc
}

func describeTypeOfValue(val interface{}) string {
	return describeType(reflect.TypeOf(val))
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