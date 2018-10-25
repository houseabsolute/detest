package row

import (
	"fmt"
	"os"
	"sort"

	"github.com/houseabsolute/detest/internal/table/border"
	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/debug"
	"github.com/houseabsolute/detest/internal/table/renderinfo"
	"github.com/houseabsolute/detest/internal/table/style"
)

type Row struct {
	// The key is the column where the cell starts. This isn't an array in
	// order to handle cells which span >1 column.
	cells map[int]*cell.Cell
}

func New() *Row {
	return &Row{cells: map[int]*cell.Cell{}}
}

func (r *Row) AddCell(c *cell.Cell) {
	i := 0
	for _, c := range r.cells {
		i += c.ColSpan
	}
	r.cells[i] = c
}

func (r *Row) ColumnNumbers() []int {
	nums := []int{}
	for n := range r.cells {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	return nums
}

func (r *Row) Cell(i int) *cell.Cell {
	return r.cells[i]
}

func (r *Row) ColumnSeparatorPositions(ri renderinfo.RI) []int {
	positions := []int{}

	for _, i := range r.ColumnNumbers() {
		c := r.Cell(i)
		pos := 0
		for j := 0; j <= i+(c.ColSpan-1); j++ {
			// The "+ 1" is to account for the vertical bar that starts the
			// table.
			pos += ri.ColumnWidth(j) + 1
		}
		// The very last column's separator is the right border of the table.
		if pos >= ri.TotalWidth() {
			continue
		}
		positions = append(positions, pos)
	}
	return positions
}

func (r *Row) Render(ri renderinfo.RI, sty style.Style) (string, error) {
	rendered := border.Vertical.String()

	if debug.Debug {
		os.Stderr.WriteString("\n  Rendering row\n")
	}

	for _, i := range r.ColumnNumbers() {
		c := r.Cell(i)

		var w int
		from := i
		to := i + (c.ColSpan - 1)
		if debug.Debug {
			fmt.Fprintf(os.Stderr, "    cell %d spans from column %d to %d (cs: %d)\n", i, from, to, c.ColSpan)
		}
		if from == to {
			w = ri.ColumnWidth(from)
		} else {
			for x := from; x <= to; x++ {
				w += ri.ColumnWidth(x)
			}
		}
		// We need to add one more for each separator that the cell spans.
		w += (c.ColSpan - 1)
		i += c.ColSpan

		if debug.Debug {
			fmt.Fprintf(os.Stderr, "      w = %d\n", w)
		}

		s, err := c.Render(w, sty)
		if err != nil {
			return "", err
		}

		rendered += s
		rendered += border.Vertical.String()
	}

	return rendered, nil
}
