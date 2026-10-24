package main

import (
	"fmt"
	"sync"
)

/**
	Fan-OUT : distributing work across multiple goroutines.
	Fan-IN  : merging results from multiple goroutine into a single channel.
**/

// Create n output channel for per worker (worker = n)
func FanOut(n int, input <-chan int) []<-chan string {
	output := make([]<-chan string, n)

	for i := 0; i < n; i++ {
		output[i] = worker(i, input)
	}
	return output
}

// Combine multiple channels into a single output channel
func FanIn(channel []<-chan string) <-chan string {
	var wg sync.WaitGroup
	combined := make(chan string)

	for _, outChan := range channel {
		wg.Add(1)
		go func(out <-chan string) {
			defer wg.Done()
			for r := range out {
				combined <- r
			}
		}(outChan)
	}

	go func() {
		wg.Wait()
		close(combined)
	}()
	return combined
}

func worker(n int, input <-chan int) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		for i := range input {
			result <- fmt.Sprintf("%d : %d", n, i)
		}
	}()
	return result
}

func main() {

	input := make(chan int, 3)
	go func() {
		defer close(input)
		for i := 0; i < 14; i++ {
			input <- i
		}
	}()

	chanOut := FanOut(3, input)
	res := FanIn(chanOut)

	for r := range res {
		fmt.Println(r)
	}
}
