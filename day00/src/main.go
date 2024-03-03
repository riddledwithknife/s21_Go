package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
)

func mean(numbers []int) float64 {
	sum := 0

	for _, num := range numbers {
		sum += num
	}

	return float64(sum) / float64(len(numbers))
}

func median(numbers []int) float64 {
	n := len(numbers)

	if n%2 == 0 {
		return float64(numbers[n/2-1]+numbers[n/2]) / 2.0
	}

	return float64(numbers[n/2])
}

func mode(numbers []int) int {
	frequency := make(map[int]int)
	var maxFreq int
	var mode int

	for _, num := range numbers {
		frequency[num]++
		if frequency[num] > maxFreq {
			maxFreq = frequency[num]
			mode = num
		} else if frequency[num] == maxFreq && num < mode {
			mode = num
		}
	}

	return mode
}

func standardDeviation(numbers []int) float64 {
	m := mean(numbers)

	var sum float64
	for _, num := range numbers {
		sum += math.Pow(float64(num)-m, 2)
	}

	return math.Sqrt(sum / float64(len(numbers)))
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func main() {
	meanFlag := flag.Bool("mean", false, "Calculate mean, cmd + D")
	medianFlag := flag.Bool("median", false, "Calculate median, cmd + D")
	modeFlag := flag.Bool("mode", false, "Calculate mode, cmd + D")
	sdFlag := flag.Bool("sd", false, "Calculate standard deviation, cmd + D")
	flag.Parse()

	anyFlagsSet := *meanFlag || *medianFlag || *modeFlag || *sdFlag

	scanner := bufio.NewScanner(os.Stdin)
	var numbers []int

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil || num < -100000 || num > 100000 {
			handleError(fmt.Errorf("invalid input: %s", scanner.Text()))
		}
		numbers = append(numbers, num)
	}

	if err := scanner.Err(); err != nil {
		handleError(err)
	}

	slices.Sort(numbers)

	if !anyFlagsSet || *meanFlag {
		fmt.Printf("Mean: %.2f\n", mean(numbers))
	}
	if !anyFlagsSet || *medianFlag {
		fmt.Printf("Median: %.2f\n", median(numbers))
	}
	if !anyFlagsSet || *modeFlag {
		fmt.Printf("Mode: %d\n", mode(numbers))
	}
	if !anyFlagsSet || *sdFlag {
		fmt.Printf("SD: %.2f\n", standardDeviation(numbers))
	}
}
