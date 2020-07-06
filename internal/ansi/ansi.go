package ansi

import "regexp"

type Scheme struct {
	Strong    func(string) string
	Em        func(string) string
	Correct   func(string) string
	Incorrect func(string) string
	Warning   func(string) string
}

var DefaultScheme = Scheme{
	Strong:    bold,
	Em:        em,
	Correct:   green,
	Incorrect: red,
	Warning:   orange,
}

const endEscape = "\033[0m"

func bold(s string) string {
	return "\033[1m" + s + endEscape
}

func em(s string) string {
	return "\033[7m" + s + endEscape
}

func green(s string) string {
	return "\033[38:5:2m" + s + endEscape
}

func red(s string) string {
	return "\033[38:5:9m" + s + endEscape
}

func orange(s string) string {
	return "\033[38:5:208m" + s + endEscape
}

// Copied from github.com/apcera/termtables/cell.go with a fix to allow a
// colon as a number separator (for colors)
var (
	// Must match SGR escape sequence, which is "CSI Pm m", where the Control
	// Sequence Introducer (CSI) is "ESC ["; where Pm is "A multiple numeric
	// parameter composed of any number of single numeric parameters, separated
	// by ; character(s).  Individual values for the parameters are listed with
	// Ps" and where Ps is A single (usually optional) numeric parameter,
	// composed of one of [sic] more digits."
	//
	// In practice, the end sequence is usually given as \e[0m but reading that
	// definition, it's clear that the 0 is optional and some testing confirms
	// that it is certainly optional with MacOS Terminal 2.3, so we need to
	// support the string \e[m as a terminator too.
	ansiFilter = regexp.MustCompile(`\033\[(?:\d+(?:[:;]\d+)*)?m`)
)

func Strip(s string) string {
	return ansiFilter.ReplaceAllString(s, "")
}
