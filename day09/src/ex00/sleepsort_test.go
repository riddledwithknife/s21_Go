package ex00

import (
	"testing"
)

func TestSleepSort(t *testing.T) {
	test := struct {
		input    []int
		expected []int
	}{
		[]int{2, 2, 10, 3, -2, 4, -5},
		[]int{-5, -2, 2, 2, 3, 4, 10},
	}

	receiver := sleepSort(test.input)
	var output []int
	for i := 0; i < len(test.input); i++ {
		output = append(output, <-receiver)
	}

	for i := range test.expected {
		if output[i] != test.expected[i] {
			t.Error("Wrong")
		}
	}
}
