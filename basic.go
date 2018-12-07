package detest

import (
	"bytes"
	"reflect"
)

// ExactEqualityComparer implements exact comparison of two values.
type ExactEqualityComparer struct {
	expect interface{}
}

// Is tests that two variables are exactly equal. The first variable is the
// actual variable and the second is what is expected. The `expect` argument
// can be either a literal value or anything that implements the
// detest.Comparer interface. The final argument is the test name.
//
// Under the hood this is implemented with the ExactEqualityComparer.
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

// Equal takes an expected literal value and returns and ExactEqualityComparer
// for later use.
func (d *D) Equal(expect interface{}) ExactEqualityComparer {
	return ExactEqualityComparer{expect}
}

// Compare compares the value in d.Actual() to the expected value passed to
// Equal().
func (eec ExactEqualityComparer) Compare(d *D) {
	d.PushPath(d.NewPath(describeType(reflect.TypeOf(d.Actual())), 1, "detest.(*D).Equal"))
	defer d.PopPath()

	actual := d.Actual()
	expect := eec.expect
	result := result{
		actual: newValue(actual),
		expect: newValue(expect),
		op:     "==",
	}
	if actual == nil || expect == nil {
		// Two nils are only equal if they're also the same type.
		if actual == expect {
			result.pass = true
		} else {
			result.pass = false
			result.where = inType
		}
	}

	exp, ok := expect.([]byte)
	if !ok {
		// Need to replace this with something that traverses a data structure
		// recording our path as we go.
		result.pass = reflect.DeepEqual(actual, expect)
		result.where = inValue
	} else {
		act, ok := actual.([]byte)
		if !ok {
			result.pass = false
			result.where = inValue
		}
		if exp == nil || act == nil {
			result.pass = exp == nil && act == nil
			result.where = inValue
		} else {
			result.pass = bytes.Equal(act, exp)
			result.where = inValue
		}
	}

	d.AddResult(result)
}

// ValueEqualityComparer implements value-based comparison of two values.
type ValueEqualityComparer struct {
	expect interface{}
}

// ValueIs tests that two variables contain the same value. The first variable
// is the actual variable and the second is what is expected. The `expect`
// argument can be either a literal value or anything that implements the
// detest.Comparer interface. The final argument is the test name.
//
// If the two variables are of different types this is fine as long as one
// type can be converted to the other (for example `int32` and `int64`).
//
// Under the hood this is implemented with the ValueEqualityComparer.
func (d *D) ValueIs(actual, expect interface{}, name string) bool {
	d.ResetState(actual)
	defer d.PopActual()

	d.ValueEqual(expect).Compare(d)
	return d.ok(name)
}

// ValueEqual takes an expected literal value and returns and
// ValueEqualityComparer for later use.
func (d *D) ValueEqual(expect interface{}) ValueEqualityComparer {
	return ValueEqualityComparer{expect}
}

// Compare compares the value in d.Actual() to the expected value passed to
// Equal().
//
// XXX - this is just a compare of the ExactEqualityComparer.Compare method
// for now.
func (vec ValueEqualityComparer) Compare(d *D) {
	d.PushPath(d.NewPath(describeType(reflect.TypeOf(d.Actual())), 1, "detest.(*D).ValueEqual"))
	defer d.PopPath()

	actual := d.Actual()
	expect := vec.expect
	result := result{
		actual: newValue(actual),
		expect: newValue(expect),
		op:     "== (value)",
	}
	if actual == nil || expect == nil {
		// Two nils are only equal if they're also the same type.
		if actual == expect {
			result.pass = true
		} else {
			result.pass = false
			result.where = inType
		}
	}

	exp, ok := expect.([]byte)
	if !ok {
		// Need to replace this with something that traverses a data structure
		// recording our path as we go.
		result.pass = reflect.DeepEqual(actual, expect)
		result.where = inValue
	} else {
		act, ok := actual.([]byte)
		if !ok {
			result.pass = false
			result.where = inValue
		}
		if exp == nil || act == nil {
			result.pass = exp == nil && act == nil
			result.where = inValue
		} else {
			result.pass = bytes.Equal(act, exp)
			result.where = inValue
		}
	}

	d.AddResult(result)
}
