package detest

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Passing test", mapPassingTest},
		{"Failing test", mapFailingTest},
		{"Mix of tests", mapMixPassingAndFailingTests},
		{"Passed non-map to Map", mapPassedNonMap},
		{"Key called that does not exist in the map", mapKeyCalledThatDoesNotExist},
		{"AllValues pass", mapPassWithAllValues},
		{"AllValues fail", mapFailWithAllValues},
		{"AllValues fail with description", mapFailWithAllValuesAndDescription},
		{"AllValues not given a func", mapPassNonFuncToAllValues},
		{"AllValues func does not take the right number of args", mapFuncToAllValuesHasWrongInputSignature},
		{"AllValues func does not return the right number of args", mapFuncToAllValuesHasWrongOutputSignature},
		{"AllValues func does not return a bool", mapFuncToAllValuesDoesNotReturnBool},
		{"No call to Etc or End", mapNoCallToEtcOrEnd},
		{"Calls End but does not check all values", mapCallsEndButDoesNotCheckAllValues},
		{"Calls End but does not check all all values with nested maps", mapNestedEndChecks},
		{"Calls Etc and does not check all values", mapCallsEtcAndDoesNotCheckAllValues},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func mapPassingTest(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	d.Is(
		map[int]int{1: 2},
		d.Map(func(mt *MapTester) {
			mt.End()
			mt.Key(1, 2)
		}),
		"map[1] == 2",
	)
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Assertion ok: map[1] == 2\n")
}

func mapFailingTest(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.Key(1, 3)
		}),
		"map[1] == 3",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: 2, desc: "int"},
			expect: &value{value: 3, desc: "int"},
			op:     "==",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "[1]",
					callee: "detest.(*MapTester).Key",
					caller: "detest.mapFailingTest.func1",
				},
				{
					data:   "int",
					callee: "detest.(*D).Equal",
					caller: "detest.mapFailingTest.func1",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func mapMixPassingAndFailingTests(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2, 2: 3, 3: 4},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.Key(1, 2)
			mt.Key(2, 4)
			mt.Key(3, 4)
		}),
		"map mix",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with three output items")
	assert.Equal(
		t,
		true,
		r.record[0].output[0].result.pass,
		"first result was a pass",
	)
	assert.Equal(
		t,
		&result{
			actual: &value{value: 3, desc: "int"},
			expect: &value{value: 4, desc: "int"},
			op:     "==",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "[2]",
					callee: "detest.(*MapTester).Key",
					caller: "detest.mapMixPassingAndFailingTests.func1",
				},
				{
					data:   "int",
					callee: "detest.(*D).Equal",
					caller: "detest.mapMixPassingAndFailingTests.func1",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[1].result,
		"got the expected second result",
	)
	assert.Equal(
		t,
		true,
		r.record[0].output[2].result.pass,
		"third result was a pass",
	)
}

func mapPassedNonMap(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		42,
		r.Map(func(mt *MapTester) {
			mt.Key(0, 1)
		}),
		"non-map",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: 42, desc: "int"},
			expect: nil,
			op:     "[]",
			pass:   false,
			path: []Path{{
				data:   "int",
				callee: "detest.(*D).Map",
				caller: "detest.(*DetestRecorder).Is",
			}},
			where:       inDataStructure,
			description: "Called detest.Map() but the value being tested isn't a map, it's an int",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func mapKeyCalledThatDoesNotExist(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.Etc()
			mt.Key(42, 1)
		}),
		"does not exist in map",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: map[int]int{1: 2}, desc: "map[int]int"},
			expect: nil,
			op:     "[42]",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "[42]",
					callee: "detest.(*MapTester).Key",
					caller: "detest.mapKeyCalledThatDoesNotExist.func1",
				},
			},
			where:       inDataStructure,
			description: "Attempted to get a map key that does not exist",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func mapPassWithAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	d.Is(
		map[int]int{1: 2, 2: 3, 3: 4},
		d.Map(func(mt *MapTester) {
			mt.End()
			mt.AllValues(func(v int) bool {
				return v < 5
			})
		}),
		"AllValues < 5",
	)
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Assertion ok: AllValues < 5\n")
}

func mapFailWithAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2, 2: 6, 3: 4},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.AllValues(func(v int) bool {
				return v < 5
			})
		}),
		"AllValues < 5",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with three output items")
	// Map iteration order is not predictable so if we don't sort the
	// results it's a lot trickier to test that we got what we expected.
	sort.SliceStable(r.record[0].output, func(i, j int) bool {
		return r.record[0].output[i].result.path[2].data < r.record[0].output[j].result.path[2].data
	})
	AssertResultsAre(
		t,
		r.record[0].output,
		[]resultExpect{
			{
				pass:     true,
				dataPath: []string{"map[int]int", "range", "[1]", "int"},
			},
			{
				pass:     false,
				dataPath: []string{"map[int]int", "range", "[2]", "int"},
			},
			{
				pass:     true,
				dataPath: []string{"map[int]int", "range", "[3]", "int"},
			},
		},
		"got expected results",
	)
}

func mapFailWithAllValuesAndDescription(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2, 2: 6, 3: 4},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.AllValues(func(v int) (bool, string) {
				return v < 5, fmt.Sprintf("expected a value less than 5 but got %d", v)
			})
		}),
		"AllValues < 5",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with three output items")
	sort.SliceStable(r.record[0].output, func(i, j int) bool {
		return r.record[0].output[i].result.path[2].data < r.record[0].output[j].result.path[2].data
	})
	AssertResultsAre(
		t,
		r.record[0].output,
		[]resultExpect{
			{
				pass:     true,
				dataPath: []string{"map[int]int", "range", "[1]", "int"},
			},
			{
				pass:     false,
				dataPath: []string{"map[int]int", "range", "[2]", "int"},
			},
			{
				pass:     true,
				dataPath: []string{"map[int]int", "range", "[3]", "int"},
			},
		},
		"got expected results",
	)
	assert.Equal(
		t,
		r.record[0].output[1].result.description,
		"expected a value less than 5 but got 6",
		"AllValues func returns a string description",
	)
}

func mapPassNonFuncToAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.AllValues(42)
		}),
		"AllValues not given a func",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: map[int]int{1: 2}, desc: "map[int]int"},
			expect: nil,
			op:     "",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "range",
					callee: "detest.(*MapTester).AllValues",
					caller: "detest.mapPassNonFuncToAllValues.func1",
				},
			},
			where:       inUsage,
			description: "you passed an int to AllValues but it needs a function",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func mapFuncToAllValuesHasWrongInputSignature(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.AllValues(func(x, y int) bool { return true })
		}),
		"AllValues func takes 2 values",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: map[int]int{1: 2}, desc: "map[int]int"},
			expect: nil,
			op:     "",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "range",
					callee: "detest.(*MapTester).AllValues",
					caller: "detest.mapFuncToAllValuesHasWrongInputSignature.func1",
				},
			},
			where:       inUsage,
			description: "the function passed to AllValues must take 1 value, but yours takes 2",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func mapFuncToAllValuesHasWrongOutputSignature(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.AllValues(func(x int) (bool, error) { return true, nil })
		}),
		"AllValues func returns 2 values",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: map[int]int{1: 2}, desc: "map[int]int"},
			expect: nil,
			op:     "",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "range",
					callee: "detest.(*MapTester).AllValues",
					caller: "detest.mapFuncToAllValuesHasWrongOutputSignature.func1",
				},
			},
			where: inUsage,
			description: "the function passed to AllValues must return a string as its" +
				" second argument but yours returns an error",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func mapFuncToAllValuesDoesNotReturnBool(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[int]int{1: 2},
		r.Map(func(mt *MapTester) {
			mt.AllValues(func(x int) int { return 42 })
		}),
		"AllValues func returns int",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: map[int]int{1: 2}, desc: "map[int]int"},
			expect: nil,
			op:     "",
			pass:   false,
			path: []Path{
				{
					data:   "map[int]int",
					callee: "detest.(*D).Map",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   "range",
					callee: "detest.(*MapTester).AllValues",
					caller: "detest.mapFuncToAllValuesDoesNotReturnBool.func1",
				},
			},
			where: inUsage,
			description: "the function passed to AllValues must return a bool as its" +
				" first argument but yours returns an int",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func mapNoCallToEtcOrEnd(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[string]int{"foo": 1},
		r.Map(func(mt *MapTester) {
			mt.Key("foo", 1)
		}),
		"no call to Etc or End",
	)
	mockT.AssertNotCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 2, "record has state with two output items")
	assert.Equal(
		t,
		"The function passed to Map() did not call Etc() or End()",
		r.record[0].output[1].warning,
		"got the expected result",
	)
}

func mapCallsEndButDoesNotCheckAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[string]int{"foo": 1, "bar": 2, "baz": 3},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.Key("foo", 1)
		}),
		"called End but did not check all values",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with two output items")
	assert.False(
		t,
		r.record[0].output[1].result.pass,
		"got a failure for the second result",
	)
	assert.Equal(
		t,
		"Your map test did not check the key bar",
		r.record[0].output[1].result.description,
		"got a failure for the second result",
	)
	assert.False(
		t,
		r.record[0].output[2].result.pass,
		"got a failure for the third result",
	)
	assert.Equal(
		t,
		"Your map test did not check the key baz",
		r.record[0].output[2].result.description,
		"got a failure for the third result",
	)
}

func mapNestedEndChecks(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[string]map[string]int{
			"foo": {
				"bar": 1,
				"baz": 2,
			},
			"quux": {
				"x": 1,
				"y": 2,
			},
		},
		r.Map(func(mt *MapTester) {
			mt.End()
			mt.Key("foo", d.Map(func(mt2 *MapTester) {
				mt2.End()
				mt2.Key("bar", 1)
				mt2.Key("baz", 2)
			}))
		}),
		"called End but did not check all values of outer map",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with four output items")
	assert.Equal(
		t,
		"Your map test did not check the key quux",
		r.record[0].output[2].result.description,
		"got a failure for the second result",
	)
}

func mapCallsEtcAndDoesNotCheckAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		map[string]int{"foo": 1, "bar": 2, "baz": 3},
		r.Map(func(mt *MapTester) {
			mt.Etc()
			mt.Key("foo", 1)
		}),
		"called Etc and did not check all values",
	)
	mockT.AssertNotCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.True(
		t,
		r.record[0].output[0].result.pass,
		"got a pass for the first result",
	)
}
