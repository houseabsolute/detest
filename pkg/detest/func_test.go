package detest

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunc(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Func with no name passes", funcWithNoNamePasses},
		{"Func with name passes", funcWithNamePasses},
		{"Func with no name fails", funcWithNoNameFails},
		{"Func with a name fails", funcWithNameFails},
		{"Func with no name fails with description", funcWithNoNameFailsWithDescription},
		{"Func with name fails with description", funcWithNameFailsWithDescription},
		{"Func cannot accept the given argument", funcCannotAcceptArgument},
		{"Func creation errors", funcCreationErrors},
		{"Func input type is interface", funcInputTypeIsInterface},
		{"Func handles untyped nil", funcHandlesUntypedNil},
		{"Func handles typed nil", funcHandlesTypedNil},
		{"Func call additional detest methods", funcCallsAdditionalDetestMethods},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func funcWithNoNamePasses(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.Func(func(s []int) bool {
		return len(s) < 4
	})
	require.NoError(t, err, "no error calling Func()")
	d.Passes(
		[]int{1, 2, 3},
		f,
		"len(s) < 4",
	)
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Assertion ok: len(s) < 4\n")
}

func funcWithNamePasses(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.NamedFunc(func(s []int) bool {
		return len(s) < 4
	}, "Has a name")
	require.NoError(t, err, "no error calling Func()")
	d.Passes(
		[]int{1, 2, 3},
		f,
		"len(s) < 4",
	)
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Assertion ok: len(s) < 4\n")
}

func funcWithNoNameFails(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.Func(func(s []int) bool {
		return len(s) < 4
	})
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		[]int{1, 2, 3, 4},
		f,
		"Func checks that array is less than 4 elements",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
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
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcWithNameFails(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.NamedFunc(func(s []int) bool {
		return len(s) < 4
	}, "Has a name")
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		[]int{1, 2, 3, 4},
		f,
		"Func checks that array is less than 4 elements",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
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
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcWithNoNameFailsWithDescription(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.Func(func(s []int) (bool, string) {
		return len(s) < 4, fmt.Sprintf("Slice is %d elements long but cannot be more than 3", len(s))
	})
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		[]int{1, 2, 3, 4},
		f,
		"Func checks that array is less than 4 elements",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
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
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "Slice is 4 elements long but cannot be more than 3",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcWithNameFailsWithDescription(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.NamedFunc(func(s []int) (bool, string) {
		return len(s) < 4, fmt.Sprintf("Slice is %d elements long but cannot be more than 3", len(s))
	}, "Has a name")
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		[]int{1, 2, 3, 4},
		f,
		"Func checks that array is less than 4 elements",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
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
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "Slice is 4 elements long but cannot be more than 3",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcCannotAcceptArgument(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.NamedFunc(func(s []int) bool {
		return len(s) < 4
	}, "Has a name")
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		42,
		f,
		"Func checks that array is less than 4 elements",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
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
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inUsage,
			description: "Called a function as a comparison that takes a []int but it was passed an int",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcCreationErrors(t *testing.T) {
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
}

func funcInputTypeIsInterface(t *testing.T) {
	t.Run("error interface input with wrong error", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(err error) bool {
			return err != nil && errors.Is(err, os.ErrNotExist)
		})
		require.NoError(t, err, "no error calling Func()")
		d.Passes(
			errors.New("foo"),
			f,
			"can pass concrete value to func that takes interface",
		)
		mockT.AssertCalled(t, "Fail")
	})

	t.Run("error interface input with right error", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(err error) bool {
			return err != nil && errors.Is(err, os.ErrNotExist)
		})
		require.NoError(t, err, "no error calling Func()")
		d.Passes(
			os.ErrNotExist,
			f,
			"can pass concrete value to func that takes interface",
		)
		mockT.AssertNotCalled(t, "Fail")
	})

	t.Run("error interface input with wrong type", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		f, err := d.Func(func(err error) bool {
			return err != nil && errors.Is(err, os.ErrNotExist)
		})
		require.NoError(t, err, "no error calling Func()")
		d.Passes(
			42,
			f,
			"can pass wrong type to func that takes interface",
		)
		mockT.AssertCalled(t, "Fail")
	})
}

func funcHandlesUntypedNil(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.Func(func(s []int) bool {
		return s != nil && len(s) < 4
	})
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	r.Passes(
		nil,
		f,
		"len(s) < 4",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: nil, desc: "nil <nil>"},
			expect: nil,
			op:     "func()",
			pass:   false,
			path: []Path{
				{
					data:   "nil",
					callee: "Func()",
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func funcHandlesTypedNil(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	f, err := d.Func(func(s []int) bool {
		return s != nil && len(s) < 4
	})
	require.NoError(t, err, "no error calling Func()")
	r := NewRecorder(d)
	var s []int
	r.Passes(
		s,
		f,
		"len(s) < 4",
	)
	mockT.AssertCalled(t, "Fail")
	require.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: s, desc: "[]int"},
			expect: nil,
			op:     "func()",
			pass:   false,
			path: []Path{
				{
					data:   "[]int",
					callee: "Func()",
					caller: "detest.(*DetestRecorder).Passes",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

type myError struct {
	size int
}

func (m myError) Error() string {
	return "foo"
}

func funcCallsAdditionalDetestMethods(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	f, err := r.Func(func(err error) bool {
		var realErr myError
		return r.Is(errors.As(err, &realErr), true) && r.Is(realErr.size, 42)
	})
	require.NoError(t, err, "no error calling Func()")
	e := myError{size: 42}
	r.Passes(
		e,
		f,
		"myError 42",
	)
	mockT.AssertNotCalled(t, "Fail")
	require.Len(t, r.record, 3, "three states were recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: true, desc: ""},
			expect: &value{value: true, desc: ""},
			op:     "==",
			pass:   true,
			path: []Path{
				{
					data:   "bool",
					callee: "detest.(*D).Equal",
					caller: "detest.(*DetestRecorder).Is",
				},
			},
		},
		r.record[0].output[0].result,
		"got the expected result for first assertion",
	)
	assert.Equal(
		t,
		&result{
			actual: &value{value: 42, desc: ""},
			expect: &value{value: 42, desc: ""},
			op:     "==",
			pass:   true,
			path: []Path{
				{
					data:   "int",
					callee: "detest.(*D).Equal",
					caller: "detest.(*DetestRecorder).Is",
				},
			},
		},
		r.record[1].output[0].result,
		"got the expected result for second assertion",
	)
}
