package ex00

import (
	"reflect"
	"testing"
)

func TestMinCoins2(t *testing.T) {
	tests := []struct {
		value    int
		coins    []int
		expected []int
	}{
		{13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
		{13, []int{1, 3, 4, 7, 13, 15}, []int{13}},
		{11, []int{1, 2, 5, 10}, []int{10, 1}},
		{0, []int{1, 5, 10}, []int{}},
		{13, []int{1, 10, 1}, []int{10, 1, 1, 1}},
		{11, []int{10, 5, 1, 2}, []int{10, 1}},
		{0, []int{}, []int{}},
	}

	for _, test := range tests {
		result := MinCoins2(test.value, test.coins)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For value %d and coins %v, expected %v but got %v", test.value, test.coins, test.expected, result)
		}
	}
}
