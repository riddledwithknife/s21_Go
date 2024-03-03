package ex01

import (
	"testing"

	"day07/ex00"
)

func BenchmarkMinCoins(b *testing.B) {
	val := 1000
	coins := []int{1, 5, 10, 50, 100, 500, 1000}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ex00.MinCoins(val, coins)
	}
}

func BenchmarkMinCoins2(b *testing.B) {
	val := 1000
	coins := []int{1, 5, 10, 50, 100, 500, 1000}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ex00.MinCoins2(val, coins)
	}
}
