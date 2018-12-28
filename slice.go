package detest

import (
	"fmt"
	"reflect"
)

// SliceComparer implements comparison of slice values.
type SliceComparer struct {
	with func(*SliceTester)
}

// Slice takes a function which will be called to do further comparisons of
// the slice's contents.
func (d *D) Slice(with func(*SliceTester)) SliceComparer {
	return SliceComparer{with}
}

// SliceTester is the struct that will be passed to the test function passed
// to detest.Slice. This struct implements the slice-specific testing methods
// such as Idx() and AllValues().
type SliceTester struct {
	d *D
}

// Compare compares the slice value in d.Actual() by calling the function
// passed to `Slice()`, which is in turn expected to further tests of the
// slice's content.
func (sc SliceComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())
	d.PushPath(d.NewPath(describeType(v.Type()), 1, "detest.(*D).Slice"))
	defer d.PopPath()

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

	sc.with(&SliceTester{d})
}

// Idx takes a slice index and an expected value for that index. If the index
// is past the end of the array, this is considered a failure.
func (st *SliceTester) Idx(idx int, expect interface{}) {
	v := reflect.ValueOf(st.d.Actual())

	st.d.PushPath(st.d.NewPath(fmt.Sprintf("[%d]", idx), 0, ""))
	defer st.d.PopPath()

	if idx >= v.Len() {
		st.d.AddResult(result{
			actual: newValue(st.d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     fmt.Sprintf("[%d]", idx),
			description: fmt.Sprintf(
				"Attempted to get an index (%d) past the end of a %d-element slice", idx, v.Len()),
		})
		return
	}

	st.d.PushActual(v.Index(idx).Interface())
	defer st.d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(st.d)
	} else {
		st.d.Equal(expect).Compare(st.d)
	}
}

// AllValues takes a function and turns it into a `FuncComparer`. It then
// passes every slice value to that comparer in turn. The function must take
// exactly one value matching the slice values' type and return a single
// boolean value.
func (st *SliceTester) AllValues(check interface{}) {
	st.d.PushPath(st.d.NewPath("range", 0, ""))
	defer st.d.PopPath()

	v := reflect.ValueOf(check)
	t := v.Type()
	if v.Kind() != reflect.Func {
		st.d.AddResult(result{
			actual:      newValue(st.d.Actual()),
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("You passed %s to AllValues but it needs a function", articleize(describeType(t))),
		})
		return
	}

	if t.NumIn() != 1 {
		st.d.AddResult(result{
			actual:      newValue(st.d.Actual()),
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must take 1 value, but yours takes %d", t.NumIn()),
		})
		return
	}

	if t.NumOut() != 1 {
		st.d.AddResult(result{
			actual:      newValue(st.d.Actual()),
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must return 1 value, but yours returns %d", t.NumOut()),
		})
		return
	}

	if t.Out(0).Name() != "bool" {
		st.d.AddResult(result{
			actual: newValue(st.d.Actual()),
			pass:   false,
			where:  inUsage,
			description: fmt.Sprintf(
				"The function passed to AllValues must return a bool, but yours returns %s",
				articleize(describeType(t.Out(0))),
			),
		})
		return
	}

	comparer := FuncComparer{comparer: v}
	array := reflect.ValueOf(st.d.Actual())
	for i := 0; i < array.Len(); i++ {
		st.Idx(i, comparer)
	}
}
