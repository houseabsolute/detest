package detest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Passing test", slicePassingTest},
		{"Failing test", sliceFailingTest},
		{"Mix of tests", sliceMixPassingAndFailingTests},
		{"Passed non-slice to Slice", slicePassedNonSlice},
		{"Idx called past end of slice", sliceIdxCalledPastEndOfSlice},
		{"AllValues pass", slicePassWithAllValues},
		{"AllValues fail", sliceFailWithAllValues},
		{"AllValues fail with description", sliceFailWithAllValuesAndDescription},
		{"AllValues not given a func", slicePassNonFuncToAllValues},
		{"AllValues func does not take the right number of args", sliceFuncToAllValuesHasWrongInputSignature},
		{"AllValues func does not return the right number of args", sliceFuncToAllValuesHasWrongOutputSignature},
		{"AllValues func does not return a bool", sliceFuncToAllValuesDoesNotReturnBool},
		{"No call to Etc or End", sliceNoCallToEtcOrEnd},
		{"Calls End but does not check all values", sliceCallsEndButDoesNotCheckAllValues},
		{"Calls End but does not check all all values with nested slices", sliceNestedEndChecks},
		{"Calls Etc and does not check all values", sliceCallsEtcAndDoesNotCheckAllValues},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func slicePassingTest(t *testing.T) {
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
}

func sliceFailingTest(t *testing.T) {
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
					caller: "detest.sliceFailingTest.func1",
				},
				{
					data:   "int",
					callee: "detest.(*D).Equal",
					caller: "detest.sliceFailingTest.func1",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func sliceMixPassingAndFailingTests(t *testing.T) {
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
					caller: "detest.sliceMixPassingAndFailingTests.func1",
				},
				{
					data:   "int",
					callee: "detest.(*D).Equal",
					caller: "detest.sliceMixPassingAndFailingTests.func1",
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

func slicePassedNonSlice(t *testing.T) {
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
}

func sliceIdxCalledPastEndOfSlice(t *testing.T) {
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
					caller: "detest.sliceIdxCalledPastEndOfSlice.func1",
				},
			},
			where:       inDataStructure,
			description: "Attempted to get an index (1) past the end of a 1-element slice",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func slicePassWithAllValues(t *testing.T) {
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
}

func sliceFailWithAllValues(t *testing.T) {
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
}

func sliceFailWithAllValuesAndDescription(t *testing.T) {
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
}

func slicePassNonFuncToAllValues(t *testing.T) {
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
					caller: "detest.slicePassNonFuncToAllValues.func1",
				},
			},
			where:       inUsage,
			description: "you passed an int to AllValues but it needs a function",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func sliceFuncToAllValuesHasWrongInputSignature(t *testing.T) {
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
					caller: "detest.sliceFuncToAllValuesHasWrongInputSignature.func1",
				},
			},
			where:       inUsage,
			description: "the function passed to AllValues must take 1 value, but yours takes 2",
		},
		r.record[0].output[0].result,
		"got expected results",
	)
}

func sliceFuncToAllValuesHasWrongOutputSignature(t *testing.T) {
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
					caller: "detest.sliceFuncToAllValuesHasWrongOutputSignature.func1",
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

func sliceFuncToAllValuesDoesNotReturnBool(t *testing.T) {
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
					caller: "detest.sliceFuncToAllValuesDoesNotReturnBool.func1",
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

func sliceNoCallToEtcOrEnd(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		[]int{1},
		r.Slice(func(st *SliceTester) {
			st.Idx(0, 1)
		}),
		"no call to Etc or End",
	)
	mockT.AssertNotCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 2, "record has state with two output items")
	assert.Equal(
		t,
		"The function passed to Slice() did not call Etc() or End()",
		r.record[0].output[1].warning,
		"got the expected result",
	)
}

func sliceCallsEndButDoesNotCheckAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		[]int{1, 2, 3},
		r.Slice(func(st *SliceTester) {
			st.End()
			st.Idx(0, 1)
		}),
		"called End but did not check all values",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 3, "record has state with three output items")
	assert.False(
		t,
		r.record[0].output[1].result.pass,
		"got a failure for the second result",
	)
	assert.Equal(
		t,
		"Your slice test did not check index 1",
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
		"Your slice test did not check index 2",
		r.record[0].output[2].result.description,
		"got a failure for the third result",
	)
}

func sliceNestedEndChecks(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		[][]int{{1, 2, 3}, {4, 5, 6}},
		r.Slice(func(st *SliceTester) {
			st.End()
			st.Idx(0, d.Slice(func(st2 *SliceTester) {
				st2.End()
				st2.Idx(0, 1)
				st2.Idx(1, 2)
				st2.Idx(2, 3)
			}))
		}),
		"called End but did not check all values of outer slice",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 4, "record has state with four output items")
	assert.Equal(
		t,
		"Your slice test did not check index 1",
		r.record[0].output[3].result.description,
		"got a failure for the second result",
	)
}

func sliceCallsEtcAndDoesNotCheckAllValues(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		[]int{1, 2, 3},
		r.Slice(func(st *SliceTester) {
			st.Etc()
			st.Idx(0, 1)
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
