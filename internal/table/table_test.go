package table

import (
	"strings"
	"testing"

	"github.com/houseabsolute/detest/internal/table/cell"
	"github.com/houseabsolute/detest/internal/table/row"
	"github.com/houseabsolute/detest/internal/table/style"
	"github.com/stretchr/testify/assert"
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
	{
		title: "Table",
		body: []testRow{{
			{content: "1"},
			{content: "2"},
			{content: "3"},
		}},
		expect: `
┍━━━━━━━━━━━┑
│   Table   │
┝━━━┯━━━┯━━━┥
│ 1 │ 2 │ 3 │
┕━━━┷━━━┷━━━┙
`,
		name: "table with title",
	},
	{
		title: "Table",
		header: []testRow{{
			{content: "H1"},
			{content: "H2"},
			{content: "H3"},
		}},
		body: []testRow{{
			{content: "1"},
			{content: "2"},
			{content: "3"},
		}},
		expect: `
┍━━━━━━━━━━━━━━┑
│    Table     │
┝━━━━┯━━━━┯━━━━┥
│ H1 │ H2 │ H3 │
├────┼────┼────┤
│ 1  │ 2  │ 3  │
┕━━━━┷━━━━┷━━━━┙
`,
		name: "table with title and header row",
	},
	{
		title: "Table",
		header: []testRow{{
			{content: "H1"},
			{content: "H2"},
			{content: "H3"},
		}},
		body: []testRow{{
			{content: "1"},
			{content: "2"},
			{content: "3"},
		}},
		footer: []testRow{{
			{content: "F1"},
			{content: "F2"},
			{content: "F3"},
		}},
		expect: `
┍━━━━━━━━━━━━━━┑
│    Table     │
┝━━━━┯━━━━┯━━━━┥
│ H1 │ H2 │ H3 │
├────┼────┼────┤
│ 1  │ 2  │ 3  │
┝━━━━┿━━━━┿━━━━┥
│ F1 │ F2 │ F3 │
┕━━━━┷━━━━┷━━━━┙
`,
		name: "table with title, header row, and footer row",
	},
	{
		title: "Table",
		header: []testRow{
			{
				{content: ""},
				{content: "ACTUAL", colSpan: 2, alignment: cell.AlignCenter},
				{content: ""},
				{content: "EXPECT", colSpan: 2, alignment: cell.AlignCenter},
				{content: ""},
			},
			{
				{content: "PATH", alignment: cell.AlignCenter},
				{content: "TYPE"},
				{content: "VALUE"},
				{content: "OP"},
				{content: "TYPE"},
				{content: "VALUE"},
				{content: "CALLER", alignment: cell.AlignCenter},
			},
		},
		body: []testRow{
			{
				{content: "map[string][]map[string][]string"},
				{content: "", colSpan: 5},
				{content: "t_test.go@134 called detest.(*D).Map"},
			},
			{
				{content: "[foo]"},
				{content: "", colSpan: 5},
				{content: "t_test.go@137 called detest.(*D).Key  "},
			},
			{
				{content: "[]map[string][]string"},
				{content: "", colSpan: 5},
				{content: "t_test.go@137 called detest.(*D).Slice"},
			},
			{
				{content: "[0]"},
				{content: "", colSpan: 5},
				{content: "t_test.go@138 called detest.(*D).Idx  "},
			},
			{
				{content: "map[string][]string"},
				{content: "", colSpan: 5},
				{content: "t_test.go@138 called detest.(*D).Map  "},
			},
			{
				{content: "[bar]"},
				{content: "", colSpan: 5},
				{content: "t_test.go@139 called detest.(*D).Key  "},
			},
			{
				{content: "[]string"},
				{content: "", colSpan: 5},
				{content: "t_test.go@139 called detest.(*D).Slice"},
			},
			{
				{content: "[2]"},
				{content: "", colSpan: 5},
				{content: "t_test.go@141 called detest.(*D).Idx  "},
			},
			{
				{content: "string"},
				{content: "", colSpan: 5},
				{content: "t_test.go@141 called detest.(*D).Equal"},
			},
			{
				{content: ""},
				{content: "string"},
				{content: "quux"},
				{content: "=="},
				{content: "string"},
				{content: "not quux"},
			},
		},
		expect: `
┍━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┑
│                                                        Table                                                        │
┝━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━┯━━━━┯━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥
│                                  │     ACTUAL     │    │      EXPECT       │                                        │
│               PATH               │ TYPE   │ VALUE │ OP │ TYPE   │ VALUE    │                 CALLER                 │
├──────────────────────────────────┼────────┴───────┴────┴────────┴──────────┼────────────────────────────────────────┤
│ map[string][]map[string][]string │                                         │ t_test.go@134 called detest.(*D).Map   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ [foo]                            │                                         │ t_test.go@137 called detest.(*D).Key   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ []map[string][]string            │                                         │ t_test.go@137 called detest.(*D).Slice │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ [0]                              │                                         │ t_test.go@138 called detest.(*D).Idx   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ map[string][]string              │                                         │ t_test.go@138 called detest.(*D).Map   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ [bar]                            │                                         │ t_test.go@139 called detest.(*D).Key   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ []string                         │                                         │ t_test.go@139 called detest.(*D).Slice │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ [2]                              │                                         │ t_test.go@141 called detest.(*D).Idx   │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ string                           │                                         │ t_test.go@141 called detest.(*D).Equal │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌┬╌╌╌╌╌╌╌┬╌╌╌╌┬╌╌╌╌╌╌╌╌┬╌╌╌╌╌╌╌╌╌╌┼╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│                                  │ string │ quux  │ == │ string │ not quux │
┕━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━┷━━━━┷━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┙
`,
		name: "table with multiple headers, multiple body rows, and cells with >1 col span",
	},
	{
		title: "Title is longer than cell contents in the body",
		body: []testRow{
			{
				{content: "C1"},
				{content: "C2"},
				{content: "C3"},
			},
		},
		expect: `
┍━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┑
│  Title is longer than cell contents in the body  │
┝━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━┥
│ C1             │ C2             │ C3             │
┕━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━┙
`,
		name: "table where title is longer than all other cells combined",
	},
}

func TestRender(t *testing.T) {
	for _, test := range tests {
		test := test
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

			output, err := tab.Render(style.New(test.includeANSI))
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
