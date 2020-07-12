package testhelper

import (
	"regexp"
	"runtime"
)

func Callback(f func() interface{}) interface{} {
	return f()
}

// This is mostly copied from detest.go but we cannot import detest in this
// package since this package is imported by the detest package's tests.
var packageRE = regexp.MustCompile(`((?:[^/]+/)*[^\.]+)\.`)

func PackageName() string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(1, pc)
	if n == 0 {
		panic("Cannot get New() from runtime.Callers!")
	}
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()

	m := packageRE.FindStringSubmatch(frame.Function)
	if len(m) == 1 {
		return ""
	}
	return m[1]
}
