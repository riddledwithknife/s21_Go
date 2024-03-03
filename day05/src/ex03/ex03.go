package knapsack

import "day05/ex02"

func grabPresents(presents []presentsheap.Present, capacity int) []presentsheap.Present {
	n := len(presents)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			if presents[i-1].Size <= w {
				include := presents[i-1].Value + dp[i-1][w-presents[i-1].Size]
				dp[i][w] = max(include, dp[i-1][w])
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	selected := make([]presentsheap.Present, 0)
	i, j := n, capacity
	for i > 0 && j > 0 {
		if dp[i][j] != dp[i-1][j] {
			selected = append(selected, presents[i-1])
			j -= presents[i-1].Size
		}
		i--
	}

	reverse(selected)
	return selected
}

func reverse(presents []presentsheap.Present) {
	for i, j := 0, len(presents)-1; i < j; i, j = i+1, j-1 {
		presents[i], presents[j] = presents[j], presents[i]
	}
}
