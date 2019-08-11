package detest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	t.Run("Func with no name passes", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(s []int) bool {
			return len(s) < 4
		})
		assert.NoError(t, err, "no error calling Func()")
		d.Is(
			[]int{1, 2, 3},
			f,
			"len(s) < 4",
		)
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: len(s) < 4\n")
	})

	t.Run("Func with name passes", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.NamedFunc(func(s []int) bool {
			return len(s) < 4
		}, "Has a name")
		assert.NoError(t, err, "no error calling Func()")
		d.Is(
			[]int{1, 2, 3},
			f,
			"len(s) < 4",
		)
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: len(s) < 4\n")
	})

	t.Run("Func with no name fails", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(s []int) bool {
			return len(s) < 4
		})
		assert.NoError(t, err, "no error calling Func()")
		r := NewRecorder(d)
		r.Is(
			[]int{1, 2, 3, 4},
			f,
			"Func checks that array is less than 4 elements",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1, 2, 3, 4}, desc: "[]int"},
				expect: nil,
				op:     "func()",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "Func()",
						caller: "detest.(*DetestRecorder).Is",
					},
				},
				where:       inValue,
				description: "",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Func with a name fails", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.NamedFunc(func(s []int) bool {
			return len(s) < 4
		}, "Has a name")
		assert.NoError(t, err, "no error calling Func()")
		r := NewRecorder(d)
		r.Is(
			[]int{1, 2, 3, 4},
			f,
			"Func checks that array is less than 4 elements",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1, 2, 3, 4}, desc: "[]int"},
				expect: nil,
				op:     "func()",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "Has a name",
						caller: "detest.(*DetestRecorder).Is",
					},
				},
				where:       inValue,
				description: "",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Func with no name fails with description", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(s []int) (bool, string) {
			return len(s) < 4, fmt.Sprintf("Slice is %d elements long but cannot be more than 3", len(s))
		})
		assert.NoError(t, err, "no error calling Func()")
		r := NewRecorder(d)
		r.Is(
			[]int{1, 2, 3, 4},
			f,
			"Func checks that array is less than 4 elements",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1, 2, 3, 4}, desc: "[]int"},
				expect: nil,
				op:     "func()",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "Func()",
						caller: "detest.(*DetestRecorder).Is",
					},
				},
				where:       inValue,
				description: "Slice is 4 elements long but cannot be more than 3",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Func with name fails with description", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.NamedFunc(func(s []int) (bool, string) {
			return len(s) < 4, fmt.Sprintf("Slice is %d elements long but cannot be more than 3", len(s))
		}, "Has a name")
		assert.NoError(t, err, "no error calling Func()")
		r := NewRecorder(d)
		r.Is(
			[]int{1, 2, 3, 4},
			f,
			"Func checks that array is less than 4 elements",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: []int{1, 2, 3, 4}, desc: "[]int"},
				expect: nil,
				op:     "func()",
				pass:   false,
				path: []Path{
					{
						data:   "[]int",
						callee: "Has a name",
						caller: "detest.(*DetestRecorder).Is",
					},
				},
				where:       inValue,
				description: "Slice is 4 elements long but cannot be more than 3",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Func cannot accept the given argument", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.NamedFunc(func(s []int) bool {
			return len(s) < 4
		}, "Has a name")
		assert.NoError(t, err, "no error calling Func()")
		r := NewRecorder(d)
		r.Is(
			42,
			f,
			"Func checks that array is less than 4 elements",
		)
		mockT.AssertCalled(t, "Fail")
		assert.Len(t, r.record, 1, "one state was recorded")
		assert.Len(t, r.record[0].output, 1, "record has state with one output item")
		assert.Equal(
			t,
			&result{
				actual: &value{value: 42, desc: "int"},
				expect: nil,
				op:     "func()",
				pass:   false,
				path: []Path{
					{
						data:   "int",
						callee: "Has a name",
						caller: "detest.(*DetestRecorder).Is",
					},
				},
				where:       inUsage,
				description: "Called a function as a comparison that takes a []int but it was passed an int",
			},
			r.record[0].output[0].result,
			"got the expected result",
		)
	})

	t.Run("Func creation errors", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)

		_, err := d.Func(42)
		assert.EqualError(
			t,
			err,
			"you passed an int to detest.Func() but it needs a function",
			"passed an int to Func()",
		)

		_, err = d.NamedFunc(42, "foo")
		assert.EqualError(
			t,
			err,
			"you passed an int to detest.NamedFunc() but it needs a function",
			"passed an int to NamedFunc()",
		)

		_, err = d.Func(func(_, _ int) bool { return false })
		assert.EqualError(
			t,
			err,
			"the function passed to detest.Func() must take 1 value, but yours takes 2",
			"function takes 2 arguments",
		)

		_, err = d.Func(func(_ int) {})
		assert.EqualError(
			t,
			err,
			"the function passed to detest.Func() must return 1 or 2 values, but yours returns 0",
			"function returns 0 arguments",
		)

		_, err = d.Func(func(_ int) (bool, string, bool) { return false, "", false })
		assert.EqualError(
			t,
			err,
			"the function passed to detest.Func() must return 1 or 2 values, but yours returns 3",
			"function returns 3 arguments",
		)

		_, err = d.Func(func(_ int) int { return 42 })
		assert.EqualError(
			t,
			err,
			"the function passed to detest.Func() must return a bool as its first argument but yours returns an int",
			"function returns an int instead of a bool",
		)

		_, err = d.Func(func(_ int) (bool, int) { return false, 42 })
		assert.EqualError(
			t,
			err,
			"the function passed to detest.Func() must return a string as its second argument but yours returns an int",
			"function returns an int instead of a string",
		)
	})
}
