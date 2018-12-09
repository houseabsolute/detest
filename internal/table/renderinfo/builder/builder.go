package builder

// This needs to be its own package so renderinfo can reference row and both
// table and the row packages can build a renderinfo. The row package only
// needs this for testing but this builder code is not something we want to
// copy and paste.

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/houseabsolute/detest/internal/table/debug"
	"github.com/houseabsolute/detest/internal/table/renderinfo"
	"github.com/houseabsolute/detest/internal/table/row"
)

func BuildRenderInfo(rows []*row.Row) renderinfo.RI {
	if debug.Debug {
		// nolint: errcheck
		os.Stderr.WriteString("  calculating column widths\n")
		// nolint: errcheck
		os.Stderr.WriteString("  -------------------------\n")
	}

	widths := map[int]int{}
	for i, r := range rows {
		for _, j := range r.ColumnNumbers() {
			c := r.Cell(j)

			w := c.DisplayWidth()
			quo := int(math.Floor(float64(w) / float64(c.ColSpan)))
			rem := int(math.Remainder(float64(w), float64(c.ColSpan)))

			if debug.Debug {
				fmt.Fprintf(
					os.Stderr,
					"    [%d][%d] (cs: %d) = %d (q = %d, r = %d)\n",
					i, j, c.ColSpan, w, quo, rem,
				)
			}

			for x := j; x <= (j + c.ColSpan - 1); x++ {
				cw := quo
				// For multi-column cells, we want to evenly apply the
				// remainder to each cell in turn until it's gone. So if we
				// have a display width of 11 across 3 cells, we end up with
				// 4, 4, and 3. If we did something more naive like round up
				// we'd end up with 4,4, and 4.
				if rem > 0 {
					cw++
					rem--
				}
				if cw > widths[x] {
					widths[x] = cw
				}
			}
		}
	}

	if debug.Debug {
		// nolint: errcheck
		os.Stderr.WriteString("  final column widths\n")
		// nolint: errcheck
		os.Stderr.WriteString("  -------------------\n")

		nums := []int{}
		for n := range widths {
			nums = append(nums, n)
		}
		sort.Ints(nums)
		for i := range nums {
			fmt.Fprintf(os.Stderr, "    %d: %d\n", i, widths[i])
		}
	}

	return renderinfo.RI(widths)
}
