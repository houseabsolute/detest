package main

import (
	"testing"

	"github.com/houseabsolute/detest"
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
		d.Slice(func(d *detest.D) {
			d.Idx(0, d.Equal(1))
			d.Idx(1, d.Equal(2))
			d.Idx(2, d.Equal(4))
			d.Idx(3, d.Equal(4))
		}),
		"slice out of bounds",
	)
}

func Test4(t *testing.T) {
	d := detest.New(t)
	d.Is(
		42,
		d.Slice(func(d *detest.D) {
			d.Idx(0, d.Equal(1))
			d.Idx(1, d.Equal(2))
			d.Idx(2, d.Equal(4))
		}),
		"slice value wrong",
	)
}

func Test4_1(t *testing.T) {
	d := detest.New(t)
	d.Is(
		42,
		d.Slice(func(d *detest.D) {
			d.Idx(0, 1)
			d.Idx(1, 1)
			d.Idx(2, 3)
		}),
		"slice compared to bare value",
	)
}

func Test5(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[][]int{{1, 2, 3}, {3, 4, 5}},
		d.Slice(func(d *detest.D) {
			d.Idx(0, d.Slice(func(d *detest.D) {
				d.Idx(0, d.Equal(1))
				d.Idx(1, d.Equal(2))
				d.Idx(2, d.Equal(3))
			}))
			d.Idx(1, d.Slice(func(d *detest.D) {
				d.Idx(0, d.Equal(2))
				d.Idx(1, d.Equal(4))
				d.Idx(2, d.Equal(5))
			}))
		}),
		"nested slice",
	)
}

func Test6(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[]int{1, 2, 3},
		d.Slice(func(d *detest.D) {
			d.AllSliceValues(func(v int) bool {
				return v > 2
			})
			d.Idx(0, d.Equal(1))
			d.Idx(1, d.Equal(2))
			d.Idx(2, d.Equal(4))
		}),
		"AllSliceValues",
	)
}

func Test7(t *testing.T) {
	d := detest.New(t)
	d.Is(
		map[string]string{"foo": "bar", "baz": "buz", "quux": "fnord"},
		d.Map(func(d *detest.D) {
			d.AllMapValues(func(v string) bool {
				return len(v) <= 3
			})
			d.Key("foo", d.Equal("bar"))
			d.Key("baz", d.Equal("buz"))
			d.Key("quux", d.Equal("fnord!"))
		}),
		"AllMapValues",
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
		d.Map(func(d *detest.D) {
			d.Key("foo", d.Slice(func(d *detest.D) {
				d.Idx(0, d.Map(func(d *detest.D) {
					d.Key("bar", d.Slice(func(d *detest.D) {
						d.Idx(1, "buz")
						d.Idx(2, "not quux")
					}))
				}))
				d.Idx(1, d.Map(func(d *detest.D) {
					d.Key("nosuchkey", d.Slice(func(d *detest.D) {
						d.Idx(1, "buz")
						d.Idx(2, "not quux")
					}))
				}))
			}))
		}),
		"map of slice of map of slice",
	)
}
