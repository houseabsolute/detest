package detest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStruct(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Passing test", structPassingTest},
		{"Failing test", structFailingTest},
		{"Mix of tests", structMixPassingAndFailingTests},
		{"Passed non-struct to Struct", structPassedNonStruct},
		{"Passed nil to Struct", structPassedNil},
		{"Field called that does not exist in the struct", structFieldCalledThatDoesNotExist},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

type s struct {
	foo string
	bar []int
}

func structPassingTest(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	d.Is(
		s{foo: "x"},
		d.Struct(func(st *StructTester) {
			st.Field("foo", "x")
			st.Field("bar", nil)
		}),
		"s.foo == x && s.bar == nil",
	)
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Assertion ok: s.foo == x && s.bar == nil\n")
}

func structFailingTest(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		s{foo: "x"},
		r.Struct(func(st *StructTester) {
			st.Field("foo", "y")
		}),
		"s.foo == y",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: "x", desc: "string"},
			expect: &value{value: "y", desc: "string"},
			op:     "==",
			pass:   false,
			path: []Path{
				{
					data:   "s",
					callee: "detest.(*D).Struct",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   ".foo",
					callee: "detest.(*StructTester).Field",
					caller: "detest.structFailingTest.func1",
				},
				{
					data:   "string",
					callee: "detest.(*D).Equal",
					caller: "detest.structFailingTest.func1",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func structMixPassingAndFailingTests(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		s{foo: "x"},
		r.Struct(func(st *StructTester) {
			st.Field("foo", "y")
			st.Field("bar", nil)
		}),
		"s.foo == y && s.bar == nil",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 2, "record has state with two output items")
	assert.Equal(
		t,
		&result{
			actual: &value{value: "x", desc: "string"},
			expect: &value{value: "y", desc: "string"},
			op:     "==",
			pass:   false,
			path: []Path{
				{
					data:   "s",
					callee: "detest.(*D).Struct",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   ".foo",
					callee: "detest.(*StructTester).Field",
					caller: "detest.structMixPassingAndFailingTests.func1",
				},
				{
					data:   "string",
					callee: "detest.(*D).Equal",
					caller: "detest.structMixPassingAndFailingTests.func1",
				},
			},
			where:       inValue,
			description: "",
		},
		r.record[0].output[0].result,
		"got the expected first result",
	)
	assert.Equal(
		t,
		true,
		r.record[0].output[1].result.pass,
		"second result was a pass",
	)
}

func structPassedNonStruct(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		42,
		r.Struct(func(st *StructTester) {
			st.Field("foo", "x")
		}),
		"non-struct",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: 42, desc: "int"},
			expect: nil,
			op:     ".",
			pass:   false,
			path: []Path{{
				data:   "int",
				callee: "detest.(*D).Struct",
				caller: "detest.(*DetestRecorder).Is",
			}},
			where:       inDataStructure,
			description: "Called detest.Struct() but the value being tested isn't a struct, it's an int",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func structPassedNil(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		nil,
		r.Struct(func(st *StructTester) {
			st.Field("foo", "x")
		}),
		"non-struct",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: nil, desc: "nil <nil>"},
			expect: nil,
			op:     ".",
			pass:   false,
			path: []Path{{
				data:   "nil",
				callee: "detest.(*D).Struct",
				caller: "detest.(*DetestRecorder).Is",
			}},
			where:       inDataStructure,
			description: "Called detest.Struct() but the value being tested isn't a struct, it's a nil",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}

func structFieldCalledThatDoesNotExist(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	r := NewRecorder(d)
	r.Is(
		s{foo: "x", bar: []int{1}},
		r.Struct(func(st *StructTester) {
			st.Field("what", 42)
		}),
		"does not exist in struct",
	)
	mockT.AssertCalled(t, "Fail")
	assert.Len(t, r.record, 1, "one state was recorded")
	assert.Len(t, r.record[0].output, 1, "record has state with one output item")
	assert.Equal(
		t,
		&result{
			actual: &value{value: s{foo: "x", bar: []int{1}}, desc: "s"},
			expect: nil,
			op:     ".what",
			pass:   false,
			path: []Path{
				{
					data:   "s",
					callee: "detest.(*D).Struct",
					caller: "detest.(*DetestRecorder).Is",
				},
				{
					data:   ".what",
					callee: "detest.(*StructTester).Field",
					caller: "detest.structFieldCalledThatDoesNotExist.func1",
				},
			},
			where:       inDataStructure,
			description: "Attempted to get a struct field that does not exist",
		},
		r.record[0].output[0].result,
		"got the expected result",
	)
}
