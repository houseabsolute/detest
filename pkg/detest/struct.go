package detest

import (
	"fmt"
	"reflect"
	"unsafe"
)

// StructComparer implements comparison of struct values.
type StructComparer struct {
	with func(*StructTester)
}

// Struct takes a function which will be called to do further comparisons of
// the struct's contents. Note that you must pass a struct _pointer_ to this
// method if you want to access private fields with the StructComparer.Field()
// method.
func (d *D) Struct(with func(*StructTester)) StructComparer {
	return StructComparer{with}
}

// StructTester is the struct that will be passed to the test function passed
// to detest.Struct. This struct implements the struct-specific testing methods
// such as Idx() and AllValues().
type StructTester struct {
	d *D
}

// Compare compares the struct value in d.Actual() by calling the function
// passed to `Struct()`, which is in turn expected to further tests of the
// struct's content.
func (sc StructComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())

	d.PushPath(d.NewPath(describeTypeOfReflectValue(v), 1, "detest.(*D).Struct"))
	defer d.PopPath()

	if !v.IsValid() ||
		(v.Kind() != reflect.Struct &&
			(v.Kind() != reflect.Ptr &&
				v.Elem().Kind() != reflect.Struct)) {
		d.AddResult(result{
			actual: newValue(d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     ".",
			description: fmt.Sprintf(
				"Called detest.Struct() but the value being tested isn't a struct, it's %s",
				articleize(describeTypeOfReflectValue(v)),
			),
		})
		return
	}

	st := &StructTester{d: d}
	sc.with(st)
}

// Field takes a field name and an expected value for that field. If the field
// does not exist, this is considered a failure.
func (st *StructTester) Field(field string, expect interface{}) {
	v := reflect.ValueOf(st.d.Actual())

	st.d.PushPath(st.d.NewPath(fmt.Sprintf(".%v", field), 0, ""))
	defer st.d.PopPath()

	// This is hack to be able to get private fields from structs (as opposed
	// to struct pointers, where this is a little simpler). We need to copy
	// the original Value into an addressable Value.
	v2 := v
	if v.Kind() == reflect.Struct {
		v2 = reflect.New(v.Type()).Elem()
		v2.Set(v)
	}

	f := v2.FieldByName(field)
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	if !f.IsValid() {
		st.d.AddResult(result{
			actual:      newValue(st.d.Actual()),
			pass:        false,
			where:       inDataStructure,
			op:          fmt.Sprintf(".%s", field),
			description: "Attempted to get a struct field that does not exist",
		})
		return
	}

	st.d.PushActual(f.Interface())
	defer st.d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(st.d)
	} else {
		st.d.Equal(expect).Compare(st.d)
	}
}
