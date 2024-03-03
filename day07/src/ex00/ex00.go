package ex00

import "sort"

func MinCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

// MinCoins2 calculates the minimum number of coins needed to reach a certain value.
//
// This function accepts a value and a slice of coin denominations. It returns a slice
// containing the minimum number of coins needed to reach the given value.
//
// Compared to the original minCoins function, this implementation has the following differences:
// 1. Handles cases with duplicate denominations.
// 2. Handles unsorted denominations by sorting them in descending order.
// 3. Handles empty input by returning an empty slice.
//
// Optimization:
// The function sorts the denominations in descending order before processing, allowing it to
// efficiently iterate through the denominations and select the largest possible coin at each step.
func MinCoins2(val int, coins []int) []int {
	if len(coins) == 0 || val == 0 {
		return []int{}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(coins)))

	res := make([]int, 0)
	i := 0

	for val > 0 && i < len(coins) {
		if coins[i] <= val {
			res = append(res, coins[i])
			val -= coins[i]
		} else {
			i++
		}
	}

	return res
}

// Generate documentation:
// To generate HTML documentation, run the following command in the terminal:
//   godoc -http=:6060
// Then open a web browser and navigate to http://localhost:6060/pkg/day07/ex00
// Replace <package-name> with the name of the package containing your code.
