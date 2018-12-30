package detest

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This looks like a dupe of the code in slice_test because it's so similar in
// many spots.
//
// nolint: dupl
func TestMap(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(
			map[int]int{1: 2},
			d.Map(func(mt *MapTester) {
				mt.Key(1, 2)
			}),
			"map[1] == 2",
		)
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: map[1] == 2\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(mt *MapTester) {
				mt.Key(1, 3)
			}),
			"map[1] == 3",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func2.1",
					},
					{
						data:   "int",
						callee: "detest.(*D).Equal",
						caller: "detest.TestMap.func2.1",
					},
				},
				where:       inValue,
				description: "",
			},
			r.record[0].results[0],
			"got the expected result",
		)
	})

	t.Run("Mix of tests", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2, 2: 3, 3: 4},
			r.Map(func(st *MapTester) {
				st.Key(1, 2)
				st.Key(2, 4)
				st.Key(3, 4)
			}),
			"map mix",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 3, "record has state with three results")
		assert.Equal(
			t,
			true,
			r.record[0].results[0].pass,
			"first result was a pass",
		)
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func3.1",
					},
					{
						data:   "int",
						callee: "detest.(*D).Equal",
						caller: "detest.TestMap.func3.1",
					},
				},
				where:       inValue,
				description: "",
			},
			r.record[0].results[1],
			"got the expected second result",
		)
		assert.Equal(
			t,
			true,
			r.record[0].results[2].pass,
			"third result was a pass",
		)
	})

	t.Run("Passed non-map to Map", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			42,
			r.Map(func(st *MapTester) {
				st.Key(0, 1)
			}),
			"non-map",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
			r.record[0].results[0],
			"got the expected result",
		)
	})

	t.Run("Key called that does not exist in the map", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(st *MapTester) {
				st.Key(42, 1)
			}),
			"does not exist in map",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func5.1",
					},
				},
				where:       inDataStructure,
				description: "Attempted to get a map key that does not exist",
			},
			r.record[0].results[0],
			"got the expected result",
		)
	})

	t.Run("AllValues pass", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(
			map[int]int{1: 2, 2: 3, 3: 4},
			d.Map(func(st *MapTester) {
				st.AllValues(func(v int) bool {
					return v < 5
				})
			}),
			"AllValues < 5",
		)
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: AllValues < 5\n")
	})

	t.Run("AllValues fail", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2, 2: 6, 3: 4},
			r.Map(func(st *MapTester) {
				st.AllValues(func(v int) bool {
					return v < 5
				})
			}),
			"AllValues < 5",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 3, "record has state with three results")
		// Map iteration order is not predictable so if we don't sort the
		// results it's a lot trickier to test that we got what we expected.
		sort.SliceStable(r.record[0].results, func(i, j int) bool {
			return r.record[0].results[i].path[2].data < r.record[0].results[j].path[2].data
		})
		AssertResultsAre(
			t,
			r.record[0].results,
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
	})

	t.Run("AllValues fail with description", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2, 2: 6, 3: 4},
			r.Map(func(st *MapTester) {
				st.AllValues(func(v int) (bool, string) {
					return v < 5, fmt.Sprintf("expected a value less than 5 but got %d", v)
				})
			}),
			"AllValues < 5",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 3, "record has state with three results")
		sort.SliceStable(r.record[0].results, func(i, j int) bool {
			return r.record[0].results[i].path[2].data < r.record[0].results[j].path[2].data
		})
		AssertResultsAre(
			t,
			r.record[0].results,
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
			r.record[0].results[1].description,
			"expected a value less than 5 but got 6",
			"AllValues func returns a string description",
		)
	})

	t.Run("AllValues not given a func", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(st *MapTester) {
				st.AllValues(42)
			}),
			"AllValues not given a func",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func9.1",
					},
				},
				where:       inUsage,
				description: "You passed an int to AllValues but it needs a function",
			},
			r.record[0].results[0],
			"got expected results",
		)
	})

	t.Run("AllValues func does not take the right number of args", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(st *MapTester) {
				st.AllValues(func(x, y int) bool { return true })
			}),
			"AllValues func takes 2 values",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func10.1",
					},
				},
				where:       inUsage,
				description: "The function passed to AllValues must take 1 value, but yours takes 2",
			},
			r.record[0].results[0],
			"got expected results",
		)
	})

	t.Run("AllValues func does not return the right number of args", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(st *MapTester) {
				st.AllValues(func(x int) (bool, error) { return true, nil })
			}),
			"AllValues func returns 2 values",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func11.1",
					},
				},
				where: inUsage,
				description: "The function passed to AllValues must return a string as its" +
					" second argument but yours returns an error",
			},
			r.record[0].results[0],
			"got expected results",
		)
	})

	t.Run("AllValues func does not return a bool", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			map[int]int{1: 2},
			r.Map(func(st *MapTester) {
				st.AllValues(func(x int) int { return 42 })
			}),
			"AllValues func returns int",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].results, 1, "record has state with one result")
		assert.Equal(
			t,
			result{
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
						caller: "detest.TestMap.func12.1",
					},
				},
				where: inUsage,
				description: "The function passed to AllValues must return a bool as its" +
					" first argument but yours returns an int",
			},
			r.record[0].results[0],
			"got expected results",
		)
	})
}
