package detest

import (
	"fmt"
	"reflect"
)

// SliceComparer implements comparison of slice values.
type SliceComparer struct {
	with func(*D)
}

// Slice takes a function which will be called to do further comparisons of
// the slice's contents.
func (d *D) Slice(with func(*D)) SliceComparer {
	return SliceComparer{with}
}

// Compare compares the slice value in d.Actual() by calling the function
// passed to `Slice()`, which is in turn expected to further tests of the
// slice's content.
func (sc SliceComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())
	if v.Kind() != reflect.Slice {
		d.AddResult(result{
			actual: newValue(d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     "[]",
			description: fmt.Sprintf(
				"Called detest.Slice() but the value being tested isn't a slice, it's %s",
				articleize(describeType(v.Type())),
			),
		})
		return
	}

	d.PushPath(d.NewPath(describeType(v.Type()), 1, "detest.(*D).Slice"))
	defer d.PopPath()

	sc.with(d)
}

// Idx takes a slice index and an expected value for that index. If the index
// is past the end of the array, this is considered a failure.
func (d *D) Idx(idx int, expect interface{}) {
	v := reflect.ValueOf(d.Actual())

	d.PushPath(d.NewPath(fmt.Sprintf("[%d]", idx), 0, ""))
	defer d.PopPath()

	if idx >= v.Len() {
		d.AddResult(result{
			actual:      newValue(d.Actual()),
			pass:        false,
			where:       inDataStructure,
			op:          fmt.Sprintf("[%d]", idx),
			description: "Attempted to get an element past the end of the slice",
		})
		return
	}

	d.PushActual(v.Index(idx).Interface())
	defer d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(d)
	} else {
		d.Equal(expect).Compare(d)
	}
}

// AllSliceValues takes a function and turns it into a `FuncComparer`. It then
// passes every slice value to that comparer in turn. The function must take
// exactly one value matching the slice values' type and return a single boolean
// value.
func (d *D) AllSliceValues(check interface{}) {
	d.PushPath(d.NewPath("{...}", 0, ""))
	defer d.PopPath()

	v := reflect.ValueOf(check)
	t := v.Type()
	if v.Kind() != reflect.Func {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("You passed a %s to AllValues but it needs a function", describeType(t)),
		})
		return
	}

	if t.NumIn() != 1 {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must take one value, but yours takes %d", t.NumIn()),
		})
		return
	}

	if t.NumOut() != 1 {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must return one value, but yours returns %d", t.NumOut()),
		})
		return
	}

	if t.Out(0).Name() != "bool" {
		d.AddResult(result{
			pass:  false,
			where: inUsage,
			description: fmt.Sprintf(
				"The function passed to AllValues must return a boolean, but yours returns %s",
				articleize(describeType(t.Out(0))),
			),
		})
		return
	}

	comparer := FuncComparer{comparer: v}
	array := reflect.ValueOf(d.Actual())
	for i := 0; i < array.Len(); i++ {
		d.Idx(i, comparer)
	}
}
