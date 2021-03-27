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
	actual := d.Actual()
	v := reflect.ValueOf(actual)

	d.PushPath(d.NewPath(describeTypeOfReflectValue(v), 1, fc.name))
	defer d.PopPath()

	inType := fc.comparer.Type().In(0)
	okInput := false
	if v.IsValid() {
		// Either the types are the same or the input implements the input
		// type is an interface and the value implements it.
		if v.Type() == inType ||
			(inType.Kind() == reflect.Interface && v.Type().Implements(inType)) {
			okInput = true
		}
	} else {
		if isNilable(inType.Kind()) {
			okInput = true
		}
	}

	if !okInput {
		d.AddResult(result{
			actual: newValue(actual),
			pass:   false,
			op:     "func()",
			where:  inUsage,
			description: fmt.Sprintf(
				"Called a function as a comparison that takes %s but it was passed %s",
				articleize(describeType(inType)),
				articleize(describeTypeOfReflectValue(v)),
			),
		})
		return
	}

	// If it's a bare nil we need to make a zero value of whatever type the
	// func is expecting. If we try to just pass the (invalid) bare nil, then
	// the `.Call(...)` will panic with "Call using zero Value argument".
	//
	// This seems really wonky but AFAICT this is actually what the
	// interpreter is doing to! Run this code to see it in action:
	//
	// f := func(s []int) {
	//     log.Printf("%v", reflect.TypeOf(s))
	// }
	// f(s)
	// f(nil)

	if !v.IsValid() {
		v = reflect.Zero(inType)
	}
	ret := fc.comparer.Call([]reflect.Value{v})
	r := result{
		actual: newValue(actual),
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
