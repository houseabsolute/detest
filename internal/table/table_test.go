package table

import (
	"strings"
	"testing"

	"github.com/autarch/testify/assert"
	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/row"
	"github.com/houseabsolute/detest/internal/table/style"
)

type testCell struct {
	content   string
	colSpan   int
	alignment cell.Alignment
}

type testRow []testCell

type testCase struct {
	title       string
	includeANSI bool
	header      []testRow
	body        []testRow
	footer      []testRow
	expect      string
	name        string
}

var tests = []testCase{
	{
		body: []testRow{{{content: "one cell"}}},
		expect: `
┍━━━━━━━━━━┑
│ one cell │
┕━━━━━━━━━━┙
`,
		name: "one cell table",
	},
	{
		body: []testRow{{
			{content: "1"},
			{content: "2"},
			{content: "3"},
		}},
		expect: `
┍━━━┯━━━┯━━━┑
│ 1 │ 2 │ 3 │
┕━━━┷━━━┷━━━┙
`,
		name: "three cell table",
	},
	{
		body: []testRow{{
			{content: ""},
			{content: ""},
			{content: ""},
		}},
		expect: `
┍━━┯━━┯━━┑
│  │  │  │
┕━━┷━━┷━━┙
`,
		name: "three empty cell table",
	},
	{
		body: []testRow{
			{
				{content: "foo"},
				{content: ""},
				{content: "bar"},
			},
			{
				{content: ""},
				{content: "buz"},
				{content: ""},
			},
		},
		expect: `
┍━━━━━┯━━━━━┯━━━━━┑
│ foo │     │ bar │
├╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┤
│     │ buz │     │
┕━━━━━┷━━━━━┷━━━━━┙
`,
		name: "mix of empty and populated cells between rows",
	},
}

func TestRender(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tab *Table
			if test.title != "" {
				tab = NewWithTitle(test.title)
			} else {
				tab = New()
			}
			for _, h := range test.header {
				tab.AddHeaderRow(makeRow(h))
			}
			for _, b := range test.body {
				tab.AddRow(makeRow(b))
			}
			for _, f := range test.footer {
				tab.AddFooterRow(makeRow(f))
			}

			output, err := tab.Render(style.Default)
			assert.Nil(t, err, "No error from Render()")
			assert.Equal(t, strings.TrimPrefix(test.expect, "\n"), output, "got expected output")
		})
	}
}

func makeRow(testRow testRow) *row.Row {
	r := row.New()
	for _, c := range testRow {
		span := c.colSpan
		if span == 0 {
			span = 1
		}
		r.AddCell(cell.NewWithParams(c.content, span, c.alignment))
	}
	return r
}
