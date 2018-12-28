package detest

import (
	"fmt"
	"reflect"
)

// FuncComparer implements comparison using a user-defined function.
//
// XXX - it also needs to be designed so the user can modify the description
// of the failure in the result.
type FuncComparer struct {
	comparer reflect.Value
}

// Compare calls the user-provided function with the value currently in
// `d.Actual()`. The function is expected to return a boolean indicating
// success or failure.
func (fc FuncComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())
	inType := fc.comparer.Type().In(0)
	if !v.Type().ConvertibleTo(inType) {
		d.AddResult(result{
			pass:  false,
			where: inUsage,
			description: fmt.Sprintf(
				"Cannot convert %s to %s",
				articleize(describeType(v.Type())),
				articleize(describeType(inType)),
			),
		})
		return
	}

	ret := fc.comparer.Call([]reflect.Value{v})
	r := result{
		actual: newValue(d.Actual()),
		pass:   ret[0].Bool(),
		op:     "func()",
	}
	if !r.pass {
		r.where = inValue
	}

	d.AddResult(r)
}
