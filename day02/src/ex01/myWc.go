package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func processFile(filename string, countLines, countWords, countChars bool) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var count int
	if countLines {
		for scanner.Scan() {
			count++
		}
	} else if countWords || countChars {
		for scanner.Scan() {
			if countWords {
				words := bufio.NewScanner(strings.NewReader(scanner.Text()))
				words.Split(bufio.ScanWords)
				for words.Scan() {
					count++
				}
			}
			if countChars {
				count += len(scanner.Text())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func main() {
	lines := flag.Bool("l", false, "Count lines")
	words := flag.Bool("w", false, "Count words")
	characters := flag.Bool("m", false, "Count characters")
	flag.Parse()

	flagsSet := 0
	if *lines {
		flagsSet++
	}
	if *words {
		flagsSet++
	}
	if *characters {
		flagsSet++
	}
	if flagsSet > 1 {
		fmt.Println("Error: Only one of -l, -w, or -m can be specified")
		os.Exit(1)
	}
	if flagsSet == 0 {
		*words = true
	}

	files := flag.Args()
	if len(files) == 0 {
		log.Fatal("No files provided")
	}

	var wg sync.WaitGroup
	results := make(chan string, len(files))
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			count, err := processFile(file, *lines, *words, *characters)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file %s: %s\n", file, err)
				return
			}
			results <- fmt.Sprintf("%d\t%s", count, file)
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
