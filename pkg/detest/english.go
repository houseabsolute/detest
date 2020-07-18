package detest

import "regexp"

var vowelRE = regexp.MustCompile(`^[aeiou]`)

func articleize(noun string) string {
	if vowelRE.MatchString(noun) {
		return "an " + noun
	}
	return "a " + noun
}
