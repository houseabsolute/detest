package detest

import (
	"bytes"
	"reflect"
)

type ExactEqualityComparer struct {
	expect interface{}
}

func (d *D) Equal(expect interface{}) ExactEqualityComparer {
	return ExactEqualityComparer{expect}
}

func (eec ExactEqualityComparer) Compare(d *D) {
	actual := d.Actual()
	expect := eec.expect
	result := result{
		actual: &value{actual},
		expect: &value{expect},
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

func (d *D) ValueEqual(expect interface{}) ValueEqualityComparer {
	return ValueEqualityComparer{expect}
}

func (vec ValueEqualityComparer) Compare(d *D) {
	actual := d.Actual()
	expect := vec.expect
	result := result{
		actual: &value{actual},
		expect: &value{expect},
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
