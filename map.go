package detest

import (
	"fmt"
	"reflect"
)

// MapComparer implements comparison of map values.
type MapComparer struct {
	with func(*MapTester)
}

// Map takes a function which will be called to do further comparisons of the
// map's contents.
func (d *D) Map(with func(*MapTester)) MapComparer {
	return MapComparer{with}
}

// MapTester is the struct that will be passed to the test function passed to
// detest.Map. This struct implements the map-specific testing methods such as
// Idx() and AllValues().
type MapTester struct {
	d *D
}

// Compare compares the map value in d.Actual() by calling the function passed
// to `Map()`, which is in turn expected to further tests of the map's
// content.
func (mc MapComparer) Compare(d *D) {
	v := reflect.ValueOf(d.Actual())

	d.PushPath(d.NewPath(describeType(v.Type()), 1, "detest.(*D).Map"))
	defer d.PopPath()

	if v.Kind() != reflect.Map {
		d.AddResult(result{
			actual: newValue(d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     "[]",
			description: fmt.Sprintf(
				"Called detest.Map() but the value being tested isn't a map, it's %s",
				articleize(describeType(v.Type())),
			),
		})
		return
	}

	mc.with(&MapTester{d})
}

// Key takes a key and an expected value for that key. If the key does not
// exist, this is considered a failure.
func (mt *MapTester) Key(key interface{}, expect interface{}) {
	v := reflect.ValueOf(mt.d.Actual())

	mt.d.PushPath(mt.d.NewPath(fmt.Sprintf("[%v]", key), 0, ""))
	defer mt.d.PopPath()

	kv := reflect.ValueOf(key)
	if kv.Type() != v.Type().Key() {
		mt.d.AddResult(result{
			actual: newValue(mt.d.Actual()),
			pass:   false,
			where:  inDataStructure,
			op:     fmt.Sprintf("[%v]", key),
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
		mt.d.AddResult(result{
			actual:      newValue(mt.d.Actual()),
			pass:        false,
			where:       inDataStructure,
			op:          fmt.Sprintf("[%v]", key),
			description: "Attempted to get a map key that does not exist",
		})
		return
	}

	mt.d.PushActual(found.Interface())
	defer mt.d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(mt.d)
	} else {
		mt.d.Equal(expect).Compare(mt.d)
	}
}

// AllValues takes a function and turns it into a `FuncComparer`. It then
// passes every map value to that comparer in turn. The function must take
// exactly one value matching the map key's type and return a single boolean
// value.
func (mt *MapTester) AllValues(check interface{}) {
	mt.d.PushPath(mt.d.NewPath("range", 0, ""))
	defer mt.d.PopPath()

	comparer, err := mt.d.NamedFunc(check, "AllValues")
	if err != nil {
		mt.d.AddResult(result{
			actual:      newValue(mt.d.Actual()),
			pass:        false,
			where:       inUsage,
			description: err.Error(),
		})
		return
	}

	mapVal := reflect.ValueOf(mt.d.Actual())
	for _, k := range mapVal.MapKeys() {
		mt.Key(k.Interface(), comparer)
	}
}
