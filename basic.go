package detest

import (
	"bytes"
	"reflect"
)

type ExactEqualityComparer struct {
	expect interface{}
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

func (d *D) Equal(expect interface{}) ExactEqualityComparer {
	return ExactEqualityComparer{expect}
}

func (eec ExactEqualityComparer) Compare(d *D) {
	d.PushPath(d.makePath(describeType(reflect.TypeOf(d.Actual())), 1, "detest.(*D).Equal"))
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

type ValueEqualityComparer struct {
	expect interface{}
}

func (d *D) ValueIs(actual, expect interface{}, name string) bool {
	d.ResetState(actual)
	defer d.PopActual()

	d.ValueEqual(expect).Compare(d)
	return d.ok(name)
}

func (d *D) ValueEqual(expect interface{}) ValueEqualityComparer {
	return ValueEqualityComparer{expect}
}

func (vec ValueEqualityComparer) Compare(d *D) {
	d.PushPath(d.makePath(describeType(reflect.TypeOf(d.Actual())), 1, "detest.(*D).ValueEqual"))
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
