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

// Func takes a function and returns a new FuncComparer using that
// function. The function provided must accept one value and return one or two
// values. The first return value must be a bool. If the function returns two
// values, the second must be a string describing a failure.
func (d *D) Func(with interface{}) (FuncComparer, error) {
	return d.newFunc(with, "Func")
}

// NamedFunc works like Func but uses the given creator name in error messages
// about the function (accepts the wrong number of values, etc.). This is
// useful for things like the slice AllValues() operation, which wants to be
// able to return an error like "The function passed to AllValues must take 1
// value, yours takes 2".
func (d *D) NamedFunc(with interface{}, creator string) (FuncComparer, error) {
	return d.newFunc(with, creator)
}

func (d *D) newFunc(with interface{}, creator string) (FuncComparer, error) {
	v := reflect.ValueOf(with)
	t := v.Type()
	if v.Kind() != reflect.Func {
		return FuncComparer{},
			fmt.Errorf("You passed %s to %s but it needs a function", articleize(describeType(t)), creator)
	}

	if t.NumIn() != 1 {
		return FuncComparer{},
			fmt.Errorf("The function passed to %s must take 1 value, but yours takes %d", creator, t.NumIn())
	}
	if !(t.NumOut() == 1 || t.NumOut() == 2) {
		return FuncComparer{},
			fmt.Errorf("The function passed to %s must return 1 or 2 values, but yours returns %d", creator, t.NumOut())
	}
	if t.Out(0).Kind() != reflect.Bool {
		return FuncComparer{},
			fmt.Errorf("The function passed to %s must return a bool as its first argument but it returns %s",
				creator, articleize(describeType(t.Out(0))))
	}
	if t.NumOut() == 2 && t.Out(1).Kind() != reflect.String {
		return FuncComparer{},
			fmt.Errorf("The function passed to %s must return a string as its second argument but it returns %s",
				creator, articleize(describeType(t.Out(1))))
	}

	return FuncComparer{v}, nil
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
	if len(ret) == 2 {
		r.description = ret[1].String()
	}
	if !r.pass {
		r.where = inValue
	}

	d.AddResult(r)
}
