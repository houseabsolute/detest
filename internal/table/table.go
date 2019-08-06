package table

import (
	"errors"
	"fmt"
	"os"

	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/debug"
	"github.com/houseabsolute/detest/internal/table/renderinfo/builder"
	"github.com/houseabsolute/detest/internal/table/row"
	"github.com/houseabsolute/detest/internal/table/style"
)

type Table struct {
	title  string
	header []*row.Row
	body   []*row.Row
	footer []*row.Row
}

func New() *Table {
	return &Table{}
}

func NewWithTitle(title string) *Table {
	return &Table{title: title}
}

func (t *Table) AddRow(cells ...interface{}) {
	t.body = append(t.body, maybeMakeRow(cells...))
}

func (t *Table) AddHeaderRow(cells ...interface{}) {
	t.header = append(t.header, maybeMakeRow(cells...))
}

func (t *Table) AddFooterRow(cells ...interface{}) {
	t.footer = append(t.footer, maybeMakeRow(cells...))
}

func maybeMakeRow(cells ...interface{}) *row.Row {
	if len(cells) == 1 {
		if r, ok := cells[0].(*row.Row); ok {
			return r
		}
	}

	row := row.New()
	for _, pc := range cells {
		if c, ok := pc.(*cell.Cell); ok {
			row.AddCell(c)
		} else if c, ok := pc.(string); ok {
			row.AddCell(cell.New(c))
		} else {
			row.AddCell(cell.New(fmt.Sprintf("%v", c)))
		}
	}
	return row
}

func (t *Table) Render(style style.Style) (string, error) {
	if len(t.body) == 0 {
		return "", errors.New("cannot render a table without a body")
	}

	if debug.Debug {
		// nolint: errcheck
		os.Stderr.WriteString("Rendering table\n")
		if t.title != "" {
			// nolint: errcheck
			os.Stderr.WriteString("  has a title\n")
		} else {
			// nolint: errcheck
			os.Stderr.WriteString("  no title\n")
		}
		fmt.Fprintf(os.Stderr, "  %d header row(s)\n", len(t.header))
		fmt.Fprintf(os.Stderr, "  %d body row(s)\n", len(t.body))
		fmt.Fprintf(os.Stderr, "  %d footer row(s)\n", len(t.footer))
	}

	rowsWithTitle := []*row.Row{}
	if t.title != "" {
		rowsWithTitle = append(rowsWithTitle, t.titleRow())
	}
	rowsWithTitle = append(rowsWithTitle, t.allRows()...)
	ri := builder.BuildRenderInfo(rowsWithTitle)
	items := t.rowsWithSeparators()
	rendered := ""
	for j, item := range items {
		var s string
		var err error
		if r, ok := item.(*row.Row); ok {
			s, err = r.Render(ri, style)
		} else if sep, ok := item.(Separator); ok {
			var before, after *row.Row
			if j > 0 {
				before = items[j-1].(*row.Row)
			}
			if j < len(items)-1 {
				after = items[j+1].(*row.Row)
			}
			s, err = sep.Render(ri, before, after)
		}

		if err != nil {
			return "", err
		}
		rendered += s + "\n"
	}

	return rendered, nil
}

func (t *Table) titleRow() *row.Row {
	title := row.New()
	title.AddCell(cell.NewWithParams(t.title, t.maxCells(), cell.AlignCenter))
	return title
}

func (t *Table) rowsWithSeparators() []interface{} {
	rows := []interface{}{}
	rows = append(rows, Separator{sepType: Start})

	if t.title != "" {
		rows = append(rows, t.titleRow())
		rows = append(rows, Separator{sepType: AfterTitle})
	}

	if len(t.header) != 0 {
		for _, h := range t.header {
			rows = append(rows, h)
		}

		rows = append(rows, Separator{sepType: AfterHeader})
	}

	for i, b := range t.body {
		rows = append(rows, b)
		if i < len(t.body)-1 {
			rows = append(rows, Separator{sepType: InBody})
		}
	}

	if len(t.footer) != 0 {
		rows = append(rows, Separator{sepType: AfterBody})
		for _, f := range t.footer {
			rows = append(rows, f)
		}
	}

	rows = append(rows, Separator{sepType: End})

	return rows
}

func (t *Table) maxCells() int {
	max := 0
	for _, r := range t.allRows() {
		count := 0
		for _, n := range r.ColumnNumbers() {
			c := r.Cell(n)
			count += c.ColSpan
		}
		if count > max {
			max = count
		}
	}

	return max
}

func (t *Table) allRows() []*row.Row {
	rows := []*row.Row{}
	rows = append(rows, t.header...)
	rows = append(rows, t.body...)
	rows = append(rows, t.footer...)
	return rows
}
