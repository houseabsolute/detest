package detest

import "testing"

type stringish string
type structy struct {
	ignored string
}

func Test_typeOf(t *testing.T) {
	d := New(t)

	strequal(t, d.typeOf("const string"), "string")
	s := "var string"
	strequal(t, d.typeOf(s), "string")
	strequal(t, d.typeOf(&s), "*string")
	p := &s
	strequal(t, d.typeOf(&p), "**string")
	strequal(t, d.typeOf(stringish("foo")), "stringish")

	strequal(t, d.typeOf(42), "int")
	i := 42
	strequal(t, d.typeOf(i), "int")

	strequal(t, d.typeOf(uint8(4)), "uint8")

	strequal(t, d.typeOf([]string{}), "[]string")
	strequal(t, d.typeOf([4]string{}), "[4]string")

	strequal(t, d.typeOf(structy{}), "structy")

	strequal(t, d.typeOf(map[string]string{}), "map[string]string")
}

func strequal(t *testing.T, got, expect string) {
	if got == expect {
		t.Logf("%s == %s", got, expect)
	} else {
		t.Logf("%s != %s", got, expect)
		t.Fail()
	}
}
