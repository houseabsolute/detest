package detest

import (
	"testing"
)

func TestIs(t *testing.T) {
	mockT := new(mockT)
	d := NewWithOutput(mockT, mockT)
	d.Is(1, 1, "1 == 1")
	mockT.AssertNotCalled(t, "Fail")
	mockT.AssertCalled(t, "WriteString", "Passed test: 1 == 1\n")
}
