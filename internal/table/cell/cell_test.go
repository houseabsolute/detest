package cell

import (
	"fmt"
	"testing"

	"github.com/autarch/testify/assert"
	"github.com/houseabsolute/detest/internal/ansi"
	"github.com/houseabsolute/detest/internal/table/style"
)

func TestDisplayWidth(t *testing.T) {
	type testCase struct {
		content string
		expect  int
		name    string
	}
	tests := []testCase{
		{
			content: "foo",
			expect:  5,
			name:    "ASCII",
		},
		{
			content: ansi.DefaultScheme.Strong("foo"),
			expect:  5,
			name:    "ASCII with ANSI",
		},
		{
			content: "你好，世界",
			expect:  12,
			name:    "Chinese characters",
		},
		{
			content: ansi.DefaultScheme.Strong("你好，世界"),
			expect:  12,
			name:    "Chinese characters with ANSI",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, New(test.content).DisplayWidth())
		})
	}
}

func Test_pad(t *testing.T) {
	type testCase struct {
		content   string
		width     int
		alignment Alignment
		expect    string
		name      string
	}
	tests := []testCase{
		{
			content:   "foo",
			width:     5,
			alignment: AlignLeft,
			expect:    " foo ",
			name:      "simple left alignment",
		},
		{
			content:   "foo",
			width:     8,
			alignment: AlignLeft,
			expect:    " foo    ",
			name:      "left alignment with 1 extra space to fill - pads right first",
		},
		{
			content:   "foo",
			width:     9,
			alignment: AlignLeft,
			expect:    " foo     ",
			name:      "left alignment with many extra spaces to fill",
		},
		{
			content:   "foo",
			width:     5,
			alignment: AlignRight,
			expect:    " foo ",
			name:      "simple right alignment",
		},
		{
			content:   "foo",
			width:     8,
			alignment: AlignRight,
			expect:    "    foo ",
			name:      "right alignment with 1 extra space to fill - pads left first",
		},
		{
			content:   "foo",
			width:     9,
			alignment: AlignRight,
			expect:    "     foo ",
			name:      "right alignment with many extra spaces to fill",
		},
		{
			content:   "foo",
			width:     5,
			alignment: AlignCenter,
			expect:    " foo ",
			name:      "simple center alignment",
		},
		{
			content:   "foo",
			width:     8,
			alignment: AlignCenter,
			expect:    "  foo   ",
			name:      "center alignment with 1 extra space to fill - pads left first",
		},
		{
			content:   "foo",
			width:     9,
			alignment: AlignCenter,
			expect:    "   foo   ",
			name:      "center alignment with many extra spaces to fill",
		},
		{
			content:   ansi.DefaultScheme.Strong("foo"),
			width:     9,
			alignment: AlignLeft,
			expect:    " " + ansi.DefaultScheme.Strong("foo") + "     ",
			name:      "left alignment with many extra spaces to fill and ansi content",
		},
		{
			content:   ansi.DefaultScheme.Strong("foo"),
			width:     9,
			alignment: AlignRight,
			expect:    "     " + ansi.DefaultScheme.Strong("foo") + " ",
			name:      "left alignment with many extra spaces to fill and ansi content",
		},
		{
			content:   ansi.DefaultScheme.Strong("foo"),
			width:     9,
			alignment: AlignCenter,
			expect:    "   " + ansi.DefaultScheme.Strong("foo") + "   ",
			name:      "center alignment with many extra spaces to fill and ansi content",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := NewWithParams(test.content, 1, test.alignment).pad(style.Default, test.width)
			assert.Equal(t, test.expect, output)
		})
	}
}

func TestRender(t *testing.T) {
	// The core of this code is already tested by our tests of DisplayWidth()
	// and pad(), but we should test error handling.
	for _, w := range []int{1, 2, 3, 4} {
		_, err := New("foo").Render(w, style.Default)
		assert.EqualError(
			t,
			err,
			fmt.Sprintf("Cell needs width of 5 but was only allowed %d", w),
			"width of %d is too small for Render", w,
		)
	}
	output, err := New("foo").Render(5, style.Default)
	assert.Nil(t, err, "no error from Render when given enough width")
	assert.Equal(t, " foo ", output, "got expected rendering")
}
