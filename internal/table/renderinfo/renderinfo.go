package renderinfo

type RI map[int]int

func (ri RI) TotalWidth() int {
	// The total width is the width of all the cells plus all the inner
	// separators. The starting value here is a count of all the separators
	// that will appear in the table.
	total := len(ri) - 1
	for _, w := range ri {
		total += w
	}
	return total
}

func (ri RI) ColumnWidth(c int) int {
	return ri[c]
}
