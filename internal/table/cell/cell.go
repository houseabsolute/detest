package cell

import (
	"fmt"
	"math"
	"strings"

	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/houseabsolute/detest/internal/table/style"
	"github.com/mattn/go-runewidth"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

type Cell struct {
	content   string
	ColSpan   int
	alignment Alignment
}

// All cells should have at least one space to the left and right of their
// content.
const minPaddingWidth = 2

func New(content string) *Cell {
	return &Cell{
		content:   content,
		ColSpan:   1,
		alignment: AlignLeft,
	}
}

func NewWithParams(content string, span int, alignment Alignment) *Cell {
	return &Cell{
		content:   content,
		ColSpan:   span,
		alignment: alignment,
	}
}

func (c *Cell) Render(width int, sty style.Style) (string, error) {
	if width < c.DisplayWidth() {
		return "", fmt.Errorf("Cell needs width of %d but was only allowed %d", c.DisplayWidth(), width)
	}

	return c.pad(sty, width), nil
}

func (c *Cell) pad(sty style.Style, width int) string {
	var display string
	if sty.IncludeANSI {
		display = c.content
	} else {
		display = ansi.Strip(c.content)
	}

	padded := " "
	padded += c.align(display, width-minPaddingWidth)
	padded += " "
	return padded
}

func (c *Cell) align(content string, width int) string {
	extra := width - displayWidth(content)
	if c.alignment == AlignLeft {
		return content + strings.Repeat(" ", extra)
	} else if c.alignment == AlignRight {
		return strings.Repeat(" ", extra) + content
	} else if c.alignment == AlignCenter {
		// If this division by two is X.5 we want the extra space _after_ the
		// content, not before. It looks better that way.
		before := int(math.Floor(float64(extra) / 2))
		after := int(math.Ceil(float64(extra) / 2))
		return strings.Repeat(" ", before) + content + strings.Repeat(" ", after)
	} else {
		panic(fmt.Sprintf("Very surprising alignmnent: %d", c.alignment))
	}
}

func (c *Cell) DisplayWidth() int {
	return displayWidth(c.content) + minPaddingWidth
}

func displayWidth(content string) int {
	return runewidth.StringWidth(ansi.Strip(content))
}
