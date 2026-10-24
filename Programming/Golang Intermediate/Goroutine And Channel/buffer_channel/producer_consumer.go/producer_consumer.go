package main

import "fmt"

// Send data to channel
func Producer(ch chan int) {

	for i := 0; i <= 10; i++ {
		ch <- i
	}

	// Run without closing the channel
	close(ch)
}

// Receieve data from channel
func Consumer(ch chan int) {
	for n := range ch { // this runs until the close of channel
		fmt.Println(n * n)
	}
}

func main() {
	ch := make(chan int, 3)
	go Producer(ch)
	Consumer(ch)
}
