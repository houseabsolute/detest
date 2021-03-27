package detest

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIs(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1, "1 == 1")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: 1 == 1\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 2, "1 == 2")
		mT.AssertCalled(t, "Fail")
	})

	t.Run("Equivalent values do not compare as equal", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(int32(1), int64(2), "int32(1) == int64(2)")
		mT.AssertCalled(t, "Fail")
	})

	t.Run("Second argument is Comparer - pass", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(42, GTComparer(41), "42 > 41")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: 42 > 41\n")
	})

	t.Run("Second argument is Comparer - fail", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(42, GTComparer(43), "42 > 43")
		mT.AssertCalled(t, "Fail")
	})

	t.Run("Can handle nil", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(nil, nil, "nil == nil")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: nil == nil\n")
	})
}

func TestIsNot(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.IsNot(1, 42, "1 != 42")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: 1 != 42\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.IsNot(1, 1, "1 != 1")
		mT.AssertCalled(t, "Fail")
	})

	t.Run("Equivalent values do not compare as equal", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.IsNot(int32(1), int64(2), "int32(1) != int64(2)")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: int32(1) != int64(2)\n")
	})

	t.Run("Can handle nil", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.IsNot(nil, nil, "nil != nil")
		mT.AssertCalled(t, "Fail")
	})
}

// Go has this (so great) thing where there are different types of nils. A
// bare `nil` is not the same as a nil which comes from a types but
// uninitialized variable. But they should compare the same. Previously detest
// was panicking on these or treating them as not the same.
func TestIsNilTypeHandling(t *testing.T) {
	var sl []int
	var ma map[string]int
	var f func()
	var p *string
	var up unsafe.Pointer
	var c chan int

	type something interface {
		foo()
	}
	var intf something

	for _, v := range []interface{}{sl, ma, f, p, up, c, intf} {
		desc := describe(v)

		t.Run(fmt.Sprintf("Is(%s == nil)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)

			d.Is(v, nil, "uninit == nil")
			mT.AssertNotCalled(t, "Fail")
			mT.AssertCalled(t, "WriteString", "Assertion ok: uninit == nil\n")
		})
		t.Run(fmt.Sprintf("Is(nil == %s)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)
			d.Is(nil, v, "nil == uninit")
			mT.AssertNotCalled(t, "Fail")
			mT.AssertCalled(t, "WriteString", "Assertion ok: nil == uninit\n")
		})

		t.Run(fmt.Sprintf("IsNot(%s == nil)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)

			d.IsNot(v, nil, "uninit != nil")
			mT.AssertCalled(t, "Fail")
		})
		t.Run(fmt.Sprintf("IsNot(nil == %s)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)
			d.IsNot(nil, v, "nil != uninit")
			mT.AssertCalled(t, "Fail")
		})

		t.Run(fmt.Sprintf("ValueIs(%s == nil)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)

			d.ValueIs(v, nil, "uninit == nil")
			mT.AssertNotCalled(t, "Fail")
			mT.AssertCalled(t, "WriteString", "Assertion ok: uninit == nil\n")
		})
		t.Run(fmt.Sprintf("ValueIs(nil == %s)", desc), func(t *testing.T) {
			mT := new(mockT)
			d := NewWithOutput(mT, mT)
			d.ValueIs(nil, v, "nil == uninit")
			mT.AssertNotCalled(t, "Fail")
			mT.AssertCalled(t, "WriteString", "Assertion ok: nil == uninit\n")
		})

		// An uninitialized interface is special. Only like other nil types
		// above, it's treated as untyped and is equal to any other typed nil.
		if desc == "something interface" {
			return
		}

		for _, v2 := range []interface{}{sl, ma, f, p, up, c} {
			desc2 := describe(v2)
			if desc != desc2 {
				t.Run(fmt.Sprintf("Is(%s == %s)", desc, desc2), func(t *testing.T) {
					mT := new(mockT)
					d := NewWithOutput(mT, mT)
					d.Is(v, v2, "nil == uninit")
					mT.AssertCalled(t, "Fail")
				})
				t.Run(fmt.Sprintf("IsNot(%s == %s)", desc, desc2), func(t *testing.T) {
					mT := new(mockT)
					d := NewWithOutput(mT, mT)
					d.IsNot(v, v2, "nil != uninit")
					mT.AssertNotCalled(t, "Fail")
				})
				t.Run(fmt.Sprintf("ValueIs(%s == %s)", desc, desc2), func(t *testing.T) {
					mT := new(mockT)
					d := NewWithOutput(mT, mT)
					d.ValueIs(v, v2, "nil == uninit")
					mT.AssertCalled(t, "Fail")
				})
			}
		}
	}
}

func describe(v interface{}) string {
	d := describeTypeOfActualValue(v)
	if d == "nil" {
		return "something interface"
	}
	return d
}

func TestPasses(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(42, GTComparer(41), "42 > 41")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: 42 > 41\n")
	})

	t.Run("fail", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(42, GTComparer(43), "42 > 43")
		mT.AssertCalled(t, "Fail")
	})
}

func TestRequire(t *testing.T) {
	t.Run("d.Require passes", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Require(d.Is(1, 1, "1 == 1"))
		mT.AssertNotCalled(t, "Fatal")
		mT.AssertCalled(t, "WriteString", "Assertion ok: 1 == 1\n")
	})

	t.Run("d.Require fails", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Require(d.Is(1, 2, "1 == 1"))
		mT.AssertCalled(t, "Fatal", []interface{}{"required test failed"})
	})

	t.Run("d.Require fails and has name", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Require(d.Is(1, 2, "1 == 1"), "must have numeric equality!")
		mT.AssertCalled(t, "Fatal", []interface{}{"must have numeric equality!"})
	})
}

func TestValueIs(t *testing.T) {
	t.Run("Numeric comparisons", testNumericComparisons)
	t.Run("Numeric comparison overflow failures", testNumericComparisonOverflowFailures)
	t.Run("Numeric comparison cannot conversion failures", testNumericComparisonCannotConvertFailures)
	t.Run("Complex comparisons", testComplexComparisons)
	t.Run("String comparisons", testStringComparisons)
	t.Run("Struct comparisons", testStructComparisons)
	t.Run("Can handle nil", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(nil, nil, "nil == nil")
		mT.AssertNotCalled(t, "Fail")
		mT.AssertCalled(t, "WriteString", "Assertion ok: nil == nil\n")
	})
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
				mT := new(mockT)
				d := NewWithOutput(mT, mT)
				d.ValueIs(actual, expect, "1 == 1")
				mT.AssertNotCalled(t, "Fail")
				mT.AssertCalled(t, "WriteString", "Assertion ok: 1 == 1\n")
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
				mT := new(mockT)
				d := NewWithOutput(mT, mT)
				d.ValueIs(actual, expect, "1,1 == 1,1")
				mT.AssertNotCalled(t, "Fail")
				mT.AssertCalled(t, "WriteString", "Assertion ok: 1,1 == 1,1\n")
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
			mT := new(mockT)
			d := NewWithOutput(mT, mT)
			d.ValueIs(actual, expect, "overflow")
			mT.AssertCalled(t, "Fail")
			call := mT.FindCall("WriteString")
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
			mT := new(mockT)
			d := NewWithOutput(mT, mT)
			d.ValueIs(actual, expect, "overflow")
			mT.AssertCalled(t, "Fail")
			call := mT.FindCall("WriteString")
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
				mT := new(mockT)
				d := NewWithOutput(mT, mT)
				d.ValueIs(actual, expect, "\"foo\" == \"foo\"")
				mT.AssertNotCalled(t, "Fail")
				mT.AssertCalled(t, "WriteString", "Assertion ok: \"foo\" == \"foo\"\n")
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
				mT := new(mockT)
				d := NewWithOutput(mT, mT)
				d.ValueIs(actual, expect, "{val: 1} == {val: 1}")
				mT.AssertNotCalled(t, "Fail")
				mT.AssertCalled(t, "WriteString", "Assertion ok: {val: 1} == {val: 1}\n")
			})
		}
	}
}

func TestNameGeneration(t *testing.T) {
	t.Run("d.Is with no name", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1)
		mT.AssertCalled(t, "WriteString", "Assertion ok: unnamed d.Is call\n")
	})

	t.Run("d.Is with one string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1, "one string")
		mT.AssertCalled(t, "WriteString", "Assertion ok: one string\n")
	})

	t.Run("d.Is with multiple args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1, "got %d %s", 5, "dogs")
		mT.AssertCalled(t, "WriteString", "Assertion ok: got 5 dogs\n")
	})

	t.Run("d.Is with one non-string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1, []int{1, 2, 3})
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]\n")
	})

	t.Run("d.Is with multiple non-string args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.Is(1, 1, []int{1, 2, 3}, "foo")
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]%!(EXTRA string=foo)\n")
	})

	t.Run("d.Passes with no name", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		f, err := d.Func(func(v int) bool { return v == 1 })
		require.NoError(t, err)
		d.Passes(1, f)
		mT.AssertCalled(t, "WriteString", "Assertion ok: unnamed d.Passes call\n")
	})

	t.Run("d.Passes with one string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		f, err := d.Func(func(v int) bool { return v == 1 })
		require.NoError(t, err)
		d.Passes(1, f, "one string")
		mT.AssertCalled(t, "WriteString", "Assertion ok: one string\n")
	})

	t.Run("d.Passes with multiple args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		f, err := d.Func(func(v int) bool { return v == 1 })
		require.NoError(t, err)
		d.Passes(1, f, "got %d %s", 5, "dogs")
		mT.AssertCalled(t, "WriteString", "Assertion ok: got 5 dogs\n")
	})

	t.Run("d.Passes with one non-string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		f, err := d.Func(func(v int) bool { return v == 1 })
		require.NoError(t, err)
		d.Passes(1, f, []int{1, 2, 3})
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]\n")
	})

	t.Run("d.Passes with multiple non-string args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		f, err := d.Func(func(v int) bool { return v == 1 })
		require.NoError(t, err)
		d.Passes(1, f, []int{1, 2, 3}, "foo")
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]%!(EXTRA string=foo)\n")
	})

	t.Run("d.ValueIs with no name", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(1, 1)
		mT.AssertCalled(t, "WriteString", "Assertion ok: unnamed d.ValueIs call\n")
	})

	t.Run("d.ValueIs with one string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(1, 1, "one string")
		mT.AssertCalled(t, "WriteString", "Assertion ok: one string\n")
	})

	t.Run("d.ValueIs with multiple args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(1, 1, "got %d %s", 5, "dogs")
		mT.AssertCalled(t, "WriteString", "Assertion ok: got 5 dogs\n")
	})

	t.Run("d.ValueIs with one non-string arg", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(1, 1, []int{1, 2, 3})
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]\n")
	})

	t.Run("d.ValueIs with multiple non-string args", func(t *testing.T) {
		mT := new(mockT)
		d := NewWithOutput(mT, mT)
		d.ValueIs(1, 1, []int{1, 2, 3}, "foo")
		mT.AssertCalled(t, "WriteString", "Assertion ok: [1 2 3]%!(EXTRA string=foo)\n")
	})
}
