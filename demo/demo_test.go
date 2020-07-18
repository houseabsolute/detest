//+build demo

// Run this with `go test -tags demo -v ./demo/`

package demo

import (
	"strings"
	"testing"

	"github.com/houseabsolute/detest/pkg/detest"
)

func Test1(t *testing.T) {
	d := detest.New(t)
	d.Is(1, 1, "1 == 1")
	d.Is(2, 2, "2 == 2")
}

func Test2(t *testing.T) {
	d := detest.New(t)
	d.Is(1, 0, "1 == 0")
	d.Is(2, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, `2 == []int{...}`)
}

func Test3(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[]int{1, 2, 3},
		d.Slice(func(st *detest.SliceTester) {
			st.Idx(0, d.Equal(1))
			st.Idx(1, d.Equal(2))
			st.Idx(2, d.Equal(4))
			st.Idx(3, d.Equal(4))
		}),
		"slice out of bounds",
	)
}

func Test4(t *testing.T) {
	d := detest.New(t)
	d.Is(
		42,
		d.Slice(func(st *detest.SliceTester) {
			st.Idx(0, d.Equal(1))
			st.Idx(1, d.Equal(2))
			st.Idx(2, d.Equal(4))
		}),
		"slice value wrong",
	)
}

func Test4_1(t *testing.T) {
	d := detest.New(t)
	d.Is(
		42,
		d.Slice(func(st *detest.SliceTester) {
			st.Idx(0, 1)
			st.Idx(1, 1)
			st.Idx(2, 3)
		}),
		"slice compared to bare value",
	)
}

func Test5(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[][]int{{1, 2, 3}, {3, 4, 5}},
		d.Slice(func(st *detest.SliceTester) {
			st.Idx(0, d.Slice(func(st *detest.SliceTester) {
				st.Idx(0, d.Equal(1))
				st.Idx(1, d.Equal(2))
				st.Idx(2, d.Equal(3))
			}))
			st.Idx(1, d.Slice(func(st *detest.SliceTester) {
				st.Idx(0, d.Equal(2))
				st.Idx(1, d.Equal(4))
				st.Idx(2, d.Equal(5))
			}))
		}),
		"nested slice",
	)
}

func Test6(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[]int{1, 2, 3},
		d.Slice(func(st *detest.SliceTester) {
			st.AllValues(func(v int) bool {
				return v > 2
			})
			st.Idx(0, d.Equal(1))
			st.Idx(1, d.Equal(2))
			st.Idx(2, d.Equal(4))
		}),
		"slice AllValues",
	)
}

func Test7(t *testing.T) {
	d := detest.New(t)
	d.Is(
		map[string]string{"foo": "bar", "baz": "buz", "quux": "fnord"},
		d.Map(func(mt *detest.MapTester) {
			mt.AllValues(func(v string) bool {
				return len(v) <= 3
			})
			mt.Key("foo", d.Equal("bar"))
			mt.Key("baz", d.Equal("buz"))
			mt.Key("quux", d.Equal("fnord!"))
		}),
		"map AllValues",
	)
}

func Test8(t *testing.T) {
	d := detest.New(t)
	actual := map[string][]map[string][]string{
		"foo": {
			{
				"bar": {"baz", "buz", "quux"},
				"x":   {"y", "z"},
			},
			{
				"mar":   {"maz", "muz"},
				"Blort": {"42"},
			},
		},
		"bar": {},
		"baz": {
			{
				"buz": {"quux"},
			},
		},
	}

	d.Is(
		actual,
		d.Map(func(mt *detest.MapTester) {
			mt.Key("foo", d.Slice(func(st *detest.SliceTester) {
				st.Idx(0, d.Map(func(mt *detest.MapTester) {
					mt.Key("bar", d.Slice(func(st *detest.SliceTester) {
						st.Idx(1, "buz")
						st.Idx(2, "not quux")
					}))
				}))
				st.Idx(1, d.Map(func(mt *detest.MapTester) {
					mt.Key("nosuchkey", d.Slice(func(st *detest.SliceTester) {
						st.Idx(1, "buz")
						st.Idx(2, "not quux")
					}))
				}))
			}))
		}),
		"map of slice of map of slice",
	)
}

type HasMethods struct {
	size int
	Name string
}

func (hm HasMethods) Size() int {
	return hm.size
}

func (hm HasMethods) SizePlus(plus int) int {
	return hm.size + plus
}

func (hm HasMethods) UCName() string {
	return strings.ToUpper(hm.Name)
}

func Test9(t *testing.T) {
	hm := &HasMethods{42, "Arthur"}
	d := detest.New(t)
	d.Is(
		hm,
		d.Struct(func(st *detest.StructTester) {
			st.Field("size", 43)
			st.Field("Name", "Douglas")
		}),
		"struct fields in struct pointer",
	)
}

func Test10(t *testing.T) {
	hm := HasMethods{42, "Arthur"}
	d := detest.New(t)
	d.Is(
		hm,
		d.Struct(func(st *detest.StructTester) {
			st.Field("size", 43)
			st.Field("Name", "Douglas")
		}),
		"struct fields in struct",
	)
}
