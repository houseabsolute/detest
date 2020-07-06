package detest

import (
	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func tableWithTitle(title string, s ansi.Scheme) table.Writer {
	tw := table.NewWriter()
	tw.SetTitle(s.Strong(title))
	st := table.StyleDefault
	st.Box = table.StyleBoxLight
	st.Format.Header = text.FormatDefault
	st.Format.Footer = text.FormatDefault
	tw.SetStyle(st)
	return tw
}
