package detest

import (
	"fmt"
	"reflect"
)

// FuncComparer implements comparison using a user-defined function.
type FuncComparer struct {
	comparer reflect.Value
	name     string
}

// Func takes a function and returns a new FuncComparer using that
// function. The function provided must accept one value and return one or two
// values. The first return value must be a bool. If the function returns two
// values, the second must be a string describing a failure.
func (d *D) Func(with interface{}) (FuncComparer, error) {
	return d.newFunc(with, "Func()", "detest.Func()")
}

// NamedFunc works like Func but uses the given name in the data path for test
// output.
func (d *D) NamedFunc(with interface{}, name string) (FuncComparer, error) {
	return d.newFunc(with, name, "detest.NamedFunc()")
}

// FuncFor works like Func but uses the given name in error messages about
// the function (accepts the wrong number of values, etc.). This is useful for
// things like the slice AllValues() operation, which wants to be able to
// return an error like "The function passed to AllValues must take 1 value,
// yours takes 2".
func (d *D) FuncFor(with interface{}, called string) (FuncComparer, error) {
	return d.newFunc(with, "Func()", called)
}

func (d *D) newFunc(with interface{}, name, called string) (FuncComparer, error) {
	v := reflect.ValueOf(with)
	t := v.Type()
	if v.Kind() != reflect.Func {
		return FuncComparer{},
			fmt.Errorf("you passed %s to %s but it needs a function", articleize(describeType(t)), called)
	}

	if t.NumIn() != 1 {
		return FuncComparer{},
			fmt.Errorf("the function passed to %s must take 1 value, but yours takes %d", called, t.NumIn())
	}
	if !(t.NumOut() == 1 || t.NumOut() == 2) {
		return FuncComparer{},
			fmt.Errorf("the function passed to %s must return 1 or 2 values, but yours returns %d", called, t.NumOut())
	}
	if t.Out(0).Kind() != reflect.Bool {
		return FuncComparer{},
			fmt.Errorf("the function passed to %s must return a bool as its first argument but yours returns %s",
				called, articleize(describeType(t.Out(0))))
	}
	if t.NumOut() == 2 && t.Out(1).Kind() != reflect.String {
		return FuncComparer{},
			fmt.Errorf("the function passed to %s must return a string as its second argument but yours returns %s",
				called, articleize(describeType(t.Out(1))))
	}

	return FuncComparer{v, name}, nil
}

// Compare calls the user-provided function with the value currently in
// `d.Actual()`. The function is expected to return a boolean indicating
// success or failure.
func (fc FuncComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())

	d.PushPath(d.NewPath(describeType(v.Type()), 1, fc.name))
	defer d.PopPath()

	inType := fc.comparer.Type().In(0)
	if v.Type() != inType {
		d.AddResult(result{
			actual: newValue(d.Actual()),
			pass:   false,
			op:     "func()",
			where:  inUsage,
			description: fmt.Sprintf(
				"Called a function as a comparison that takes %s but it was passed %s",
				articleize(describeType(inType)),
				articleize(describeType(v.Type())),
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
