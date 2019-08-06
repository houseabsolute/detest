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
	d      *D
	ending CollectionEnding
	seen   map[int]bool
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

	st := &SliceTester{d: d, seen: map[int]bool{}}
	defer st.enforceEnding()
	sc.with(st)
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

	st.seen[idx] = true

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

	comparer, err := st.d.FuncFor(check, "AllValues")
	if err != nil {
		st.d.AddResult(result{
			actual:      newValue(st.d.Actual()),
			pass:        false,
			where:       inUsage,
			description: err.Error(),
		})
		return
	}

	array := reflect.ValueOf(st.d.Actual())
	for i := 0; i < array.Len(); i++ {
		st.Idx(i, comparer)
	}
}

// Etc means that not all elements of the slice will be tested.
func (st *SliceTester) Etc() {
	st.ending = Etc
}

// End means that all elements of the slice must be tested or else the test
// will fail.
func (st *SliceTester) End() {
	st.ending = End
}

func (st *SliceTester) enforceEnding() {
	// If we got an error in anything but a value check that means the test
	// aborted. This could mean attempting to get an index past the end of the
	// slice, passing an incorrect type to AllValues, etc.
	if !st.d.lastResultIsValueError() {
		return
	}

	if st.ending == Etc {
		return
	}

	if st.ending == Unset {
		st.d.AddWarning("The function passed to Slice() did not call Etc() or End()")
		return
	}

	for i := 0; i < reflect.ValueOf(st.d.Actual()).Len(); i++ {
		if !st.seen[i] {
			st.d.AddResult(result{
				pass:        false,
				where:       inUsage,
				description: fmt.Sprintf("Your slice test did not check index %d", i),
			})
		}
	}
}
