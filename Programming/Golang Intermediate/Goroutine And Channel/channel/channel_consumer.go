package main

import "fmt"

/*
 */

func main() {
	// has responsibility
	// 1. instantiate the channel
	// 2. performs writes or pass ownership
	// 3. close the channel
	//
	// 4. or encapsulate those 3, here is the example
	chanOwner := func() <-chan int { // read-only
		resultStream := make(chan int, 5) // 1
		go func() {
			defer close(resultStream) // 3
			for i := 0; i < 6; i++ {
				resultStream <- i // 2
			}
		}()
		return resultStream
	}

	// consumer -> blocking or closed channel
	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Println(result)
	}
}
