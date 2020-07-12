package detest

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
)

// ExactEqualityComparer implements exact comparison of two values.
type ExactEqualityComparer struct {
	expect interface{}
}

// Is tests that two variables are exactly equal. The first variable is the
// actual variable and the second is what is expected. The `expect` argument
// can be either a literal value or anything that implements the
// detest.Comparer interface. The final argument is the assertion name.
//
// Under the hood this is implemented with the ExactEqualityComparer.
func (d *D) Is(actual, expect interface{}, name string) bool {
	d.ResetState()
	d.PushActual(actual)
	defer d.PopActual()

	if c, ok := expect.(Comparer); ok {
		c.Compare(d)
	} else {
		d.Equal(expect).Compare(d)
	}
	return d.ok(name)
}

// Equal takes an expected literal value and returns an ExactEqualityComparer
// for later use.
func (d *D) Equal(expect interface{}) ExactEqualityComparer {
	return ExactEqualityComparer{expect}
}

// Compare compares the value in d.Actual() to the expected value passed to
// Equal().
func (eec ExactEqualityComparer) Compare(d *D) {
	actual := d.Actual()
	actualType := reflect.TypeOf(actual)
	d.PushPath(d.NewPath(describeType(actualType), 1, "detest.(*D).Equal"))
	defer d.PopPath()

	expect := eec.expect
	result := result{
		actual: newValue(actual),
		expect: newValue(expect),
		op:     "==",
	}

	if reflect.TypeOf(expect) == actualType {
		result.pass = exactCompare(actual, expect)
		if !result.pass {
			result.where = inValue
		}
	} else {
		result.pass = false
		result.where = inType
	}

	d.AddResult(result)
}

func exactCompare(actual, expect interface{}) bool {
	if actual == nil || expect == nil {
		// Two nils are only equal if they're also the same type.
		return actual == expect
	}

	exp, ok := expect.([]byte)
	if !ok {
		// Need to replace this with something that traverses a data structure
		// recording our path as we go.
		return reflect.DeepEqual(actual, expect)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}

	return bytes.Equal(act, exp)
}

// ValueEqualityComparer implements value-based comparison of two values.
type ValueEqualityComparer struct {
	expect interface{}
}

// ValueIs tests that two variables contain the same value. The first variable
// is the actual variable and the second is what is expected. The `expect`
// argument can be either a literal value or anything that implements the
// detest.Comparer interface. The final argument is the assertion name.
//
// If the two variables are of different types this is fine as long as one
// type can be converted to the other (for example `int32` and `int64`).
//
// Under the hood this is implemented with the ValueEqualityComparer.
func (d *D) ValueIs(actual, expect interface{}, name string) bool {
	d.ResetState()
	d.PushActual(actual)
	defer d.PopActual()

	d.ValueEqual(expect).Compare(d)
	return d.ok(name)
}

// ValueEqual takes an expected literal value and returns a
// ValueEqualityComparer for later use.
func (d *D) ValueEqual(expect interface{}) ValueEqualityComparer {
	return ValueEqualityComparer{expect}
}

// Compare compares the value in d.Actual() to the expected value passed to
// ValueEqual().
func (vec ValueEqualityComparer) Compare(d *D) {
	actual := d.Actual()
	actualType := reflect.TypeOf(actual)
	d.PushPath(d.NewPath(describeType(actualType), 1, "detest.(*D).ValueEqual"))
	defer d.PopPath()

	expect := vec.expect
	result := result{
		actual: newValue(actual),
		expect: newValue(expect),
		op:     "== (value)",
	}

	expectType := reflect.TypeOf(expect)
	if actualType == expectType {
		result.pass = exactCompare(actual, expect)
		if !result.pass {
			result.where = inValue
		}
		d.AddResult(result)
		return
	}

	if !actualType.ConvertibleTo(expectType) {
		result.pass = false
		result.where = inType
		result.description = cannotConvertMessage(actualType, expectType)
		d.AddResult(result)
		return
	}

	if actual == nil && expect == nil {
		result.pass = true
		d.AddResult(result)
		return
	}

	actualVal, expectVal, desc := maybeConvertValues(actual, expect, actualType, expectType)
	if desc != "" {
		result.pass = false
		result.where = inType
		result.description = desc
		d.AddResult(result)
		return
	}

	result.pass = actualVal.Interface() == expectVal.Interface()
	if !result.pass {
		result.where = inValue
	}

	d.AddResult(result)
}

func cannotConvertMessage(actualType, expectType reflect.Type) string {
	return fmt.Sprintf(
		"Cannot convert between %s and %s",
		articleize(describeType(actualType)),
		articleize(describeType(expectType)),
	)
}

func maybeConvertValues(actual, expect interface{}, actualType, expectType reflect.Type) (
	reflect.Value, reflect.Value, string,
) {
	actualVal := reflect.ValueOf(actual)
	expectVal := reflect.ValueOf(expect)

	if actualNumeric := isNumeric(actualVal); actualNumeric != nil {
		if expectNumeric := isNumeric(expectVal); expectNumeric != nil {
			return safelyConvertNumberTypes(actualVal, expectVal, actualNumeric, expectNumeric)
		}

		// It shouldn't really be possible to get here since we checked the
		// results of ConvertibleTo for the 2 types earlier.
		return actualVal, expectVal, cannotConvertMessage(actualType, expectType)
	}

	// We should only end up here when we have two non-numeric types that are
	// identical in form but not name. This can happen when one type is an
	// alias for the other (`type StringLike string`) or when two structs have
	// identical fields but different type names.
	return actualVal, expectVal.Convert(actualType), ""
}

type numericInfo struct {
	baseType string
	bits     int
}

var numericKindRE = regexp.MustCompile(`^(int|uint|float|complex|rune|byte)(8|16|32|64|128)?$`)

// From Dave Cheney via
// https://grokbase.com/t/gg/golang-nuts/14c1mpnz2e/go-nuts-is-code-running-on-32-bit-or-64-bit-platform
const is64Bit = uint64(^uint(0)) == ^uint64(0)

const uintBase = "uint"
const intBase = "int"

func isNumeric(val reflect.Value) *numericInfo {
	var matches []string
	if matches = numericKindRE.FindStringSubmatch(val.Kind().String()); len(matches) <= 1 {
		return nil
	}

	base := matches[1]
	var bits int
	if matches[2] != "" {
		// This panic _should_ be impossible to reach. If our regexp matched we
		// know that the match contains a valid int.
		var err error
		bits, err = strconv.Atoi(matches[2])
		if err != nil {
			panic(fmt.Sprintf("Could not convert %s to int: %s", matches[2], err))
		}
	} else {
		// nolint: gocritic
		if matches[1] == "byte" {
			base = uintBase
			bits = 8
		} else if matches[1] == "rune" {
			base = intBase
			bits = 32
		} else if is64Bit {
			bits = 64
		} else {
			bits = 32
		}
	}

	return &numericInfo{
		baseType: base,
		bits:     bits,
	}
}

func safelyConvertNumberTypes(actualVal, expectVal reflect.Value, actualNumeric, expectNumeric *numericInfo) (
	reflect.Value, reflect.Value, string,
) {
	// If they have the same base type then the conversion is
	// simple. Just convert to the bigger type.
	if actualNumeric.baseType == expectNumeric.baseType {
		if actualNumeric.bits > expectNumeric.bits {
			return actualVal, expectVal.Convert(actualVal.Type()), ""
		}

		return actualVal.Convert(expectVal.Type()), expectVal, ""
	}

	// If one is complex and the other is not we cannot convert between the two.
	if actualNumeric.baseType == "complex" || expectNumeric.baseType == "complex" {
		return actualVal, expectVal, cannotConvertMessage(actualVal.Type(), expectVal.Type())
	}

	// The max float32 value is _much_ bigger than the max uint64 value so we
	// can always safely convert to a float.
	if actualNumeric.baseType == "float" {
		return actualVal, expectVal.Convert(actualVal.Type()), ""
	} else if expectNumeric.baseType == "float" {
		return actualVal.Convert(expectVal.Type()), expectVal, ""
	}

	if actualNumeric.baseType == intBase && expectNumeric.baseType == uintBase {
		return intUintConversion(actualVal, expectVal, actualNumeric, expectNumeric)
	} else if actualNumeric.baseType == uintBase && expectNumeric.baseType == intBase {
		expectVal, actualVal, desc := intUintConversion(expectVal, actualVal, actualNumeric, expectNumeric)
		return actualVal, expectVal, desc
	}

	panic(
		fmt.Sprintf(
			"Should never get here - convert between %s and %s",
			actualVal.Type().Name(),
			expectVal.Type().Name()),
	)
}

func intUintConversion(int, uint reflect.Value, intInfo, uintInfo *numericInfo) (reflect.Value, reflect.Value, string) {
	// If we have an int and uint of different sizes we can always convert to
	// the bigger size safely.
	if intInfo.bits < uintInfo.bits {
		return int.Convert(uint.Type()), uint, ""
	} else if intInfo.bits > uintInfo.bits {
		return int, uint.Convert(int.Type()), ""
	}

	var intMax uint64
	switch intInfo.bits {
	case 8:
		intMax = uint64(math.MaxInt8)
	case 16:
		intMax = uint64(math.MaxInt16)
	case 32:
		intMax = uint64(math.MaxInt32)
	case 64:
		intMax = uint64(math.MaxInt64)
	}

	if uint.Uint() > intMax {
		return int, uint, fmt.Sprintf(
			"Cannot convert %d-bit uint (%d) to %d-bit int without overflow",
			uintInfo.bits, uint.Uint(), intInfo.bits,
		)
	}

	return int, uint.Convert(int.Type()), ""
}
