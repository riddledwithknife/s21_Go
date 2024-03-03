package ex01

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
)

func crawlWeb(ctx context.Context, urls <-chan string) <-chan string {
	results := make(chan string)

	var wg sync.WaitGroup

	sem := make(chan struct{}, 8)

	go func() {
		defer close(results)

		for url := range urls {
			select {
			case <-ctx.Done():
				return
			default:
			}

			sem <- struct{}{}
			wg.Add(1)

			go func(url string) {
				defer func() {
					<-sem
					wg.Done()
				}()

				resp, err := http.Get(url)
				if err != nil {
					log.Printf("Error fetching %s: %v\n", url, err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading body of %s: %v\n", url, err)
					return
				}

				results <- string(body)
			}(url)
		}

		wg.Wait()
	}()

	return results
}
