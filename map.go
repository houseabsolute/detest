package detest

import (
	"fmt"
	"reflect"
)

type MapComparer struct {
	with func(*D)
}

func (d *D) Map(with func(*D)) MapComparer {
	return MapComparer{with}
}

func (sc MapComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())
	if v.Kind() != reflect.Map {
		d.AddResult(result{
			actual:      newValue(d.Actual()),
			pass:        false,
			where:       inDataStructure,
			op:          "[]",
			description: fmt.Sprintf("Called detest.Map() but the value being tested isn't a map, it's %s", articleize(describeType(v.Type()))),
		})
		return
	}

	d.PushPath(d.makePath(describeType(v.Type()), 1, "detest.(*D).Map"))
	defer d.PopPath()

	sc.with(d)
}

func (d *D) Key(key interface{}, expect interface{}) {
	v := reflect.ValueOf(d.Actual())

	d.PushPath(d.makePath(fmt.Sprintf("[%s]", key), 0, ""))
	defer d.PopPath()

	kv := reflect.ValueOf(key)
	if kv.Type() != v.Type().Key() {
		d.AddResult(result{
			actual: newValue(d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     fmt.Sprintf("[%s]", key),
			description: fmt.Sprintf(
				"Attempted to look up a map using a key that is %s but this map uses %s as a key",
				articleize(describeType(kv.Type())),
				articleize(describeType(v.Type().Key())),
			),
		})
		return
	}

	found := v.MapIndex(kv)
	if !found.IsValid() {
		d.AddResult(result{
			actual:      newValue(d.Actual()),
			pass:        false,
			where:       inDataStructure,
			op:          fmt.Sprintf("[%s]", key),
			description: "Attempted to get a map key that does not exist",
		})
		return
	}

	d.PushActual(found.Interface())
	defer d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(d)
	} else {
		d.Equal(expect).Compare(d)
	}
}

func (d *D) AllMapValues(check interface{}) {
	d.PushPath(d.makePath("{...}", 0, ""))
	defer d.PopPath()

	v := reflect.ValueOf(check)
	t := v.Type()
	if v.Kind() != reflect.Func {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("You passed a %s to AllValues but it needs a function", describeType(t)),
		})
		return
	}

	if t.NumIn() != 1 {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must take one value, but yours takes %d", t.NumIn()),
		})
		return
	}

	if t.NumOut() != 1 {
		d.AddResult(result{
			pass:        false,
			where:       inUsage,
			description: fmt.Sprintf("The function passed to AllValues must return one value, but yours returns %d", t.NumOut()),
		})
		return
	}

	if t.Out(0).Name() != "bool" {
		d.AddResult(result{
			pass:  false,
			where: inUsage,
			description: fmt.Sprintf(
				"The function passed to AllValues must return a boolean, but yours returns %d", describeType(t.Out(0))),
		})
		return
	}

	comparer := FuncComparer{comparer: v}
	mapVal := reflect.ValueOf(d.Actual())
	for _, k := range mapVal.MapKeys() {
		d.Key(k.Interface(), comparer)
	}
}
