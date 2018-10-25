package detest

import (
	"fmt"
	"reflect"
)

type FuncComparer struct {
	comparer reflect.Value
}

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
	if ret[0].Bool() {
		return
	}

	d.AddResult(result{
		actual: &value{d.Actual()},
		pass:   false,
		where:  inValue,
		op:     "func()",
	})
}
