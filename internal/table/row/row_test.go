package row_test

// We put these tests in their own package because we want to use the
// renderinfo package as well, but that package references the row package.

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/renderinfo/builder"
	"github.com/houseabsolute/detest/internal/table/row"
	"github.com/houseabsolute/detest/internal/table/style"
)

func TestColumnNumbers(t *testing.T) {
	r := row.New()
	r.AddCell(cell.NewWithParams("", 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("", 2, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("", 3, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("", 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("", 5, cell.AlignLeft))
	assert.Equal(t, r.ColumnNumbers(), []int{0, 1, 3, 6, 7}, "got expected column numbers")
}

func TestColumnSeparatorPositions(t *testing.T) {
	r := testRow()
	widths := map[int]int{}
	for _, cn := range r.ColumnNumbers() {
		widths[cn] = r.Cell(cn).DisplayWidth()
	}
	ri := builder.BuildRenderInfo([]*row.Row{r})
	assert.Equal(
		t,
		[]int{6, 14, 27, 31},
		r.ColumnSeparatorPositions(ri),
		"got expected column separator positions",
	)
}

func TestRender(t *testing.T) {
	r := testRow()
	ri := builder.BuildRenderInfo([]*row.Row{r})
	output, err := r.Render(ri, style.Default)
	assert.Nil(t, err, "no error calling Render()")
	assert.Equal(
		t,
		"│ foo │ blort │ 你好，世界 │ \x1b[1mA\x1b[0m │ ding dong │",
		output,
		"Render() returned expected output",
	)
}

func testRow() *row.Row {
	r := row.New()
	r.AddCell(cell.NewWithParams("foo", 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("blort", 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("你好，世界", 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams(ansi.DefaultScheme.Strong("A"), 1, cell.AlignLeft))
	r.AddCell(cell.NewWithParams("ding dong", 1, cell.AlignLeft))
	return r
}
