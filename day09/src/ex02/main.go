package ex02

import "sync"

func multiplex(channels ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	copyval := func(ch <-chan interface{}) {
		defer wg.Done()
		for msg := range ch {
			out <- msg
		}
	}

	for _, ch := range channels {
		wg.Add(1)
		go copyval(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
