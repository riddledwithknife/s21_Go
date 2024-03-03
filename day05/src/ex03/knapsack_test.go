package knapsack

import (
	"day05/ex02"
	"testing"
)

func TestGrabPresents(t *testing.T) {
	presents := presentsheap.PresentHeap{{5, 1}, {4, 5}, {3, 1}, {5, 2}}

	got := grabPresents(presents, 3)
	expect := []presentsheap.Present{{5, 1}, {5, 2}}

	for i, val := range got {
		if expect[i] != val {
			t.Error("Wrong answer for capacity = 3")
		}
	}

	got = grabPresents(presents, 8)
	expect = []presentsheap.Present{{5, 1}, {4, 5}, {5, 2}}

	for i, val := range got {
		if expect[i] != val {
			t.Error("Wrong answer for capacity = 3")
		}
	}
}
