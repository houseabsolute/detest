package detest

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(1, 1, "1 == 1")
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: 1 == 1\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(1, 2, "1 == 2")
		mockT.AssertCalled(t, "Fail")
	})

	t.Run("Equivalent values do not compare as true", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(int32(1), int64(2), "int32(1) == int64(2)")
		mockT.AssertCalled(t, "Fail")
	})

	t.Run("second argument is Comparer - pass", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(42, GTComparer(41), "42 > 41")
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: 42 > 41\n")
	})

	t.Run("second argument is Comparer - fail", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(42, GTComparer(43), "42 > 43")
		mockT.AssertCalled(t, "Fail")
	})
}

func TestValueIs(t *testing.T) {
	t.Run("Numeric comparisons", testNumericComparisons)
	t.Run("Numeric comparison overflow failures", testNumericComparisonOverflowFailures)
	t.Run("Numeric comparison cannot conversion failures", testNumericComparisonCannotConvertFailures)
	t.Run("Complex comparisons", testComplexComparisons)
	t.Run("String comparisons", testStringComparisons)
	t.Run("Struct comparisons", testStructComparisons)
}

type intish int
type uintish uint

func testNumericComparisons(t *testing.T) {
	vals := []interface{}{
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		intish(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintish(1),
		float32(1),
		float64(1),
		byte(1),
		rune(1),
	}

	for _, actual := range vals {
		for _, expect := range vals {
			actual := actual
			expect := expect
			actualType := reflect.TypeOf(actual)
			expectType := reflect.TypeOf(expect)
			t.Run(fmt.Sprintf("Passing test - %s and %s", actualType.Name(), expectType.Name()), func(t *testing.T) {
				mockT := new(mockT)
				d := NewWithOutput(mockT, mockT)
				d.ValueIs(actual, expect, "1 == 1")
				mockT.AssertNotCalled(t, "Fail")
				mockT.AssertCalled(t, "WriteString", "Passed test: 1 == 1\n")
			})
		}
	}
}

func testComplexComparisons(t *testing.T) {
	vals := []interface{}{
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
	}

	for _, actual := range vals {
		for _, expect := range vals {
			actual := actual
			expect := expect
			actualType := reflect.TypeOf(actual)
			expectType := reflect.TypeOf(expect)
			t.Run(fmt.Sprintf("Passing test - %s and %s", actualType.Name(), expectType.Name()), func(t *testing.T) {
				mockT := new(mockT)
				d := NewWithOutput(mockT, mockT)
				d.ValueIs(actual, expect, "1,1 == 1,1")
				mockT.AssertNotCalled(t, "Fail")
				mockT.AssertCalled(t, "WriteString", "Passed test: 1,1 == 1,1\n")
			})
		}
	}
}

func testNumericComparisonOverflowFailures(t *testing.T) {
	overflows := [][2]interface{}{
		{uint8(math.MaxUint8), int8(0)},
		{uint16(math.MaxUint16), int16(0)},
		{uint32(math.MaxUint32), int32(0)},
		{uint64(math.MaxUint64), int64(0)},
	}
	for _, pair := range overflows {
		testOverflowHandling(t, pair[0], pair[1], pair[0])
		testOverflowHandling(t, pair[1], pair[0], pair[0])
	}
}

func testOverflowHandling(t *testing.T, actual, expect, larger interface{}) {
	actualType := reflect.TypeOf(actual)
	expectType := reflect.TypeOf(expect)
	t.Run(
		fmt.Sprintf("Failing test - %s and %s with overflow", actualType.Name(), expectType.Name()),
		func(t *testing.T) {
			mockT := new(mockT)
			d := NewWithOutput(mockT, mockT)
			d.ValueIs(actual, expect, "overflow")
			mockT.AssertCalled(t, "Fail")
			call := mockT.FindCall("WriteString")
			assert.NotNil(t, call, "WriteString was called")
			assert.Len(t, call.Args, 1, "WriteString was called with one argument")
			assert.Regexp(
				t,
				regexp.MustCompile(fmt.Sprintf(`Cannot convert .+\(%d\).+ without overflow`, larger)),
				call.Args[0],
				"WriteString was called with expected error message",
			)
		},
	)
}

func testNumericComparisonCannotConvertFailures(t *testing.T) {
	complexes := []interface{}{
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
	}

	reals := []interface{}{
		int(1),
		int8(1),
		int32(1),
		int64(1),
		intish(1),
		uint(1),
		uint8(1),
		uint32(1),
		uint64(1),
		uintish(1),
		float32(1),
		float64(1),
	}

	for _, complex := range complexes {
		for _, real := range reals {
			testCannotConvert(t, complex, real)
			testCannotConvert(t, real, complex)
		}
	}
}

func testCannotConvert(t *testing.T, actual, expect interface{}) {
	actualType := reflect.TypeOf(actual)
	expectType := reflect.TypeOf(expect)
	t.Run(
		fmt.Sprintf("Failing test - cannot convert between %s and %s", actualType.Name(), expectType.Name()),
		func(t *testing.T) {
			mockT := new(mockT)
			d := NewWithOutput(mockT, mockT)
			d.ValueIs(actual, expect, "overflow")
			mockT.AssertCalled(t, "Fail")
			call := mockT.FindCall("WriteString")
			assert.NotNil(t, call, "WriteString was called")
			assert.Len(t, call.Args, 1, "WriteString was called with one argument")
			assert.Regexp(
				t,
				regexp.MustCompile(fmt.Sprintf(`Cannot convert between .+ %s and .+ %s`, actualType.Name(), expectType.Name())),
				call.Args[0],
				"WriteString was called with expected error message",
			)
		},
	)
}

func testStringComparisons(t *testing.T) {
	type stringish string

	vals := []interface{}{
		"foo",
		stringish("foo"),
	}

	for _, actual := range vals {
		for _, expect := range vals {
			actual := actual
			expect := expect
			actualType := reflect.TypeOf(actual)
			expectType := reflect.TypeOf(expect)
			t.Run(fmt.Sprintf("Passing test - %s and %s", actualType.Name(), expectType.Name()), func(t *testing.T) {
				mockT := new(mockT)
				d := NewWithOutput(mockT, mockT)
				d.ValueIs(actual, expect, "\"foo\" == \"foo\"")
				mockT.AssertNotCalled(t, "Fail")
				mockT.AssertCalled(t, "WriteString", "Passed test: \"foo\" == \"foo\"\n")
			})
		}
	}
}

func testStructComparisons(t *testing.T) {
	type struct1 struct {
		val int
	}
	type struct2 struct {
		val int
	}

	vals := []interface{}{
		struct1{1},
		struct2{1},
	}

	for _, actual := range vals {
		for _, expect := range vals {
			actual := actual
			expect := expect
			actualType := reflect.TypeOf(actual)
			expectType := reflect.TypeOf(expect)
			t.Run(fmt.Sprintf("Passing test - %s and %s", actualType.Name(), expectType.Name()), func(t *testing.T) {
				mockT := new(mockT)
				d := NewWithOutput(mockT, mockT)
				d.ValueIs(actual, expect, "{val: 1} == {val: 1}")
				mockT.AssertNotCalled(t, "Fail")
				mockT.AssertCalled(t, "WriteString", "Passed test: {val: 1} == {val: 1}\n")
			})
		}
	}
}
