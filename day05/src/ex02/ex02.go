package presentsheap

import (
	"container/heap"
	"errors"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (ph PresentHeap) Len() int { return len(ph) }

func (ph PresentHeap) Less(i, j int) bool {
	if ph[i].Value > ph[j].Value {
		return true
	} else if ph[i].Value == ph[j].Value {
		return ph[i].Size < ph[j].Size
	}
	return false
}

func (ph PresentHeap) Swap(i, j int) { ph[i], ph[j] = ph[j], ph[i] }

func (ph *PresentHeap) Push(x interface{}) {
	*ph = append(*ph, x.(Present))
}

func (ph *PresentHeap) Pop() interface{} {
	old := *ph
	n := len(old)
	x := old[n-1]
	*ph = old[0 : n-1]
	return x
}

func getNCoolestPresents(presents []Present, n int) ([]Present, error) {
	if n < 0 || n > len(presents) {
		return nil, errors.New("invalid n value")
	}

	h := &PresentHeap{}
	heap.Init(h)

	for _, p := range presents {
		heap.Push(h, p)
	}

	coolest := make([]Present, 0, n)
	for i := 0; i < n; i++ {
		if h.Len() == 0 {
			break
		}
		coolest = append(coolest, heap.Pop(h).(Present))
	}

	return coolest, nil
}
