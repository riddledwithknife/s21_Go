package ex00

import "time"

func sleepSort(nums []int) <-chan int {
	out := make(chan int, len(nums))
	done := make(chan struct{})

	go func() {
		defer close(out)
		defer close(done)

		for _, num := range nums {
			go func(n int) {
				time.Sleep(time.Duration(n) * time.Millisecond)
				out <- n
				done <- struct{}{}
			}(num)
		}

		for i := 0; i < len(nums); i++ {
			<-done
		}
	}()

	return out
}
