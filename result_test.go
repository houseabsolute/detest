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
	assert.Equal(t, "string", describeTypeOfValue("const string"))
	s := "var string"
	assert.Equal(t, "string", describeTypeOfValue(s))
	assert.Equal(t, "*string", describeTypeOfValue(&s))
	p := &s
	assert.Equal(t, "**string", describeTypeOfValue(&p))
	assert.Equal(t, "stringish", describeTypeOfValue(stringish("foo")))

	assert.Equal(t, "int", describeTypeOfValue(42))
	i := 42
	assert.Equal(t, "int", describeTypeOfValue(i))

	assert.Equal(t, "uint8", describeTypeOfValue(uint8(4)))

	assert.Equal(t, "[]string", describeTypeOfValue([]string{}))
	assert.Equal(t, "[4]string", describeTypeOfValue([4]string{}))

	assert.Equal(t, "structy", describeTypeOfValue(structy{}))

	assert.Equal(t, "map[string]string", describeTypeOfValue(map[string]string{}))
}
