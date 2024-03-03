package ex01

import (
	"context"
	"testing"
	"time"
)

func TestCrawlWeb(t *testing.T) {
	testURLs := []string{"https://google.com", "https://yandex.net", "https://www.wikipedia.org/", "https://www.youtube.com/",
		"https://www.amazon.com/", "https://www.reddit.com/", "https://www.apple.com/", "https://www.ebay.com/"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := make(chan string, len(testURLs))
	for _, url := range testURLs {
		urls <- url
	}
	close(urls)

	results := crawlWeb(ctx, urls)

	var receivedBodies []string
	for body := range results {
		receivedBodies = append(receivedBodies, body)
	}

	if len(receivedBodies) != len(testURLs) {
		t.Errorf("Expected %d bodies, but received %d", len(testURLs), len(receivedBodies))
	}
}
