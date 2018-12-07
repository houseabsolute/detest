package detest_test

import (
	"testing"

	"github.com/houseabsolute/detest"
)

func ExampleIs(t *testing.T) {
	d := detest.New(t)
	d.Is(
		[]int{1, 2, 3},
		d.Slice(func(d *detest.D) {
			d.Idx(0, d.Equal(1))
			d.Idx(1, d.Equal(2))
			d.Idx(2, d.Equal(4))
		}),
		"Slice contains expected values",
	)
}
