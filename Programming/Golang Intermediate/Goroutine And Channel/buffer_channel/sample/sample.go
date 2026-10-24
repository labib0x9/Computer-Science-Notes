package main

import "fmt"

// send -> insert at the last of queue
// receive -> removes elements from the front of queue
// Channel is full ? sender blocks
// Channel is empty ? receiver blocks

func BufferChannel() {

	// three sized buffer channel
	ch := make(chan int, 3)

	// FIFO
	ch <- 1 // {1}
	ch <- 2 // {1, 2}
	ch <- 3 // {1, 2, 3}

	x := <-ch // {2, 3}
	y := <-ch // {3}
	z := <-ch // {}

	fmt.Println("cap:", cap(ch))
	fmt.Println("len:", len(ch))
	fmt.Println(x, y, z)
}
