package detest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringish string
type structy struct {
	// nolint: megacheck, structcheck
	ignored string
}

func Test_typeOf(t *testing.T) {
	assert.Equal(t, "string", describeTypeOfActualValue("const string"))
	s := "var string"
	assert.Equal(t, "string", describeTypeOfActualValue(s))
	assert.Equal(t, "*string", describeTypeOfActualValue(&s))
	p := &s
	assert.Equal(t, "**string", describeTypeOfActualValue(&p))
	assert.Equal(t, "stringish", describeTypeOfActualValue(stringish("foo")))

	assert.Equal(t, "int", describeTypeOfActualValue(42))
	i := 42
	assert.Equal(t, "int", describeTypeOfActualValue(i))

	assert.Equal(t, "uint8", describeTypeOfActualValue(uint8(4)))

	assert.Equal(t, "[]string", describeTypeOfActualValue([]string{}))
	assert.Equal(t, "[4]string", describeTypeOfActualValue([4]string{}))

	assert.Equal(t, "structy", describeTypeOfActualValue(structy{}))

	assert.Equal(t, "map[string]string", describeTypeOfActualValue(map[string]string{}))

	assert.Equal(t, "nil", describeTypeOfActualValue(nil))
}
