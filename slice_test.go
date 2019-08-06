package detest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This looks like a dupe of the code in map_test because it's so similar in
// many spots.
//
// nolint: dupl
func TestSlice(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(
			[]int{1},
			d.Slice(func(st *SliceTester) {
				st.Idx(0, 1)
				st.End()
			}),
			"slice[0] == 1",
		)
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: slice[0] == 1\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.Idx(0, 2)
				st.End()
			}),
			"slice[0] == 2",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: 1, desc: "int"},
				expect: &value{value: 2, desc: "int"},
				op:     "==",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "[0]",
						callee: "detest.(*SliceTester).Idx",
						caller: "detest.TestSlice.func2.1",
					},
					{
						data:   "int",
						callee: "detest.(*D).Equal",
						caller: "detest.TestSlice.func2.1",
					},
				},
				where:       inValue,
				description: "",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Mix of tests", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1, 2, 3},
			r.Slice(func(st *SliceTester) {
				st.Idx(0, 1)
				st.Idx(1, 3)
				st.Idx(2, 3)
				st.End()
			}),
			"slice mix",
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
				actual: &value{value: 2, desc: "int"},
				expect: &value{value: 3, desc: "int"},
				op:     "==",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "[1]",
						callee: "detest.(*SliceTester).Idx",
						caller: "detest.TestSlice.func3.1",
					},
					{
						data:   "int",
						callee: "detest.(*D).Equal",
						caller: "detest.TestSlice.func3.1",
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
	})

	t.Run("Passed non-slice to Slice", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			42,
			r.Slice(func(st *SliceTester) {
				st.Idx(0, 1)
				st.End()
			}),
			"non-slice",
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
					callee: "detest.(*D).Slice",
					caller: "detest.(*DetestRecorder).Is",
				}},
				where:       inDataStructure,
				description: "Called detest.Slice() but the value being tested isn't a slice, it's an int",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Idx called past end of slice", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.Idx(1, 1)
			}),
			"past end of slice",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1}, desc: "[]int"},
				expect: nil,
				op:     "[1]",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "[1]",
						callee: "detest.(*SliceTester).Idx",
						caller: "detest.TestSlice.func5.1",
					},
				},
				where:       inDataStructure,
				description: "Attempted to get an index (1) past the end of a 1-element slice",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("AllValues pass", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(
			[]int{1, 2, 3},
			d.Slice(func(st *SliceTester) {
				st.AllValues(func(v int) bool {
					return v < 5
				})
				st.End()
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
			[]int{1, 2, 6, 3},
			r.Slice(func(st *SliceTester) {
				st.AllValues(func(v int) bool {
					return v < 5
				})
				st.End()
			}),
			"AllValues < 5",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 4, "record has state with four output items")
		AssertResultsAre(
			t,
			r.record[0].output,
			[]resultExpect{
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[0]", "int"},
				},
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[1]", "int"},
				},
				{
					pass:     false,
					dataPath: []string{"[]int", "range", "[2]", "int"},
				},
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[3]", "int"},
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
			[]int{1, 2, 6, 3},
			r.Slice(func(st *SliceTester) {
				st.AllValues(func(v int) (bool, string) {
					return v < 5, fmt.Sprintf("expected a value less than 5 but got %d", v)
				})
				st.End()
			}),
			"AllValues < 5",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 4, "record has state with four output items")
		AssertResultsAre(
			t,
			r.record[0].output,
			[]resultExpect{
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[0]", "int"},
				},
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[1]", "int"},
				},
				{
					pass:     false,
					dataPath: []string{"[]int", "range", "[2]", "int"},
				},
				{
					pass:     true,
					dataPath: []string{"[]int", "range", "[3]", "int"},
				},
			},
			"got expected results",
		)
		assert.Equal(
			t,
			r.record[0].output[2].result.description,
			"expected a value less than 5 but got 6",
			"AllValues func returns a string description",
		)
	})

	t.Run("AllValues not given a func", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.AllValues(42)
			}),
			"AllValues not given a func",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1}, desc: "[]int"},
				expect: nil,
				op:     "",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "range",
						callee: "detest.(*SliceTester).AllValues",
						caller: "detest.TestSlice.func9.1",
					},
				},
				where:       inUsage,
				description: "You passed an int to AllValues but it needs a function",
			},
			r.record[0].output[0].result,
			"got expected results",
		)
	})

	t.Run("AllValues func does not take the right number of args", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.AllValues(func(x, y int) bool { return true })
				st.End()
			}),
			"AllValues func takes 2 values",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1}, desc: "[]int"},
				expect: nil,
				op:     "",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "range",
						callee: "detest.(*SliceTester).AllValues",
						caller: "detest.TestSlice.func10.1",
					},
				},
				where:       inUsage,
				description: "The function passed to AllValues must take 1 value, but yours takes 2",
			},
			r.record[0].output[0].result,
			"got expected results",
		)
	})

	t.Run("AllValues func does not return the right number of args", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.AllValues(func(x int) (bool, error) { return true, nil })
				st.End()
			}),
			"AllValues func returns 2 values",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1}, desc: "[]int"},
				expect: nil,
				op:     "",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "range",
						callee: "detest.(*SliceTester).AllValues",
						caller: "detest.TestSlice.func11.1",
					},
				},
				where: inUsage,
				description: "The function passed to AllValues must return a string as its" +
					" second argument but yours returns an error",
			},
			r.record[0].output[0].result,
			"got expected results",
		)
	})

	t.Run("AllValues func does not return a bool", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		r := NewRecorder(d)
		r.Is(
			[]int{1},
			r.Slice(func(st *SliceTester) {
				st.AllValues(func(x int) int { return 42 })
			}),
			"AllValues func returns int",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1}, desc: "[]int"},
				expect: nil,
				op:     "",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "detest.(*D).Slice",
						caller: "detest.(*DetestRecorder).Is",
					},
					{
						data:   "range",
						callee: "detest.(*SliceTester).AllValues",
						caller: "detest.TestSlice.func12.1",
					},
				},
				where: inUsage,
				description: "The function passed to AllValues must return a bool as its" +
					" first argument but yours returns an int",
			},
			r.record[0].output[0].result,
			"got expected results",
		)
	})
}
