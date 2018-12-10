package detest

import (
	"testing"
)

func TestIs(t *testing.T) {
	t.Run("Passing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(1, 1, "1 == 1")
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: 1 == 1\n")
	})

	t.Run("Failing test", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(1, 2, "1 == 2")
		mockT.AssertCalled(t, "Fail")
	})

	t.Run("Equivalent values do not compare as true", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(int32(1), int64(2), "int32(1) == int64(2)")
		mockT.AssertCalled(t, "Fail")
	})

	t.Run("second argument is Comparer - pass", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(42, GTComparer(41), "42 > 41")
		mockT.AssertNotCalled(t, "Fail")
		mockT.AssertCalled(t, "WriteString", "Passed test: 42 > 41\n")
	})

	t.Run("second argument is Comparer - fail", func(t *testing.T) {
		mockT := new(mockT)
		d := NewWithOutput(mockT, mockT)
		d.Is(42, GTComparer(43), "42 > 43")
		mockT.AssertCalled(t, "Fail")
	})

}
