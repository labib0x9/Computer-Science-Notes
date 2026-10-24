package main

import (
	"fmt"
	"time"
)

/**
	What is happening here ?
	run the code and observe the output ....

	Why this is happening ?
	because the sender and receiver must be ready at the same time, it's a
	blocking mechanism.
	until execution of both <-ch and ch<- in different goroutine,
	both of them are blocked.


	Solution - 01:
	- use buffered channel
	- use another goroutine for receiving
**/

func Run(ch chan int) {
	fmt.Println("Running !!!")
	ch <- 10          // Blocks here, waiting to send
	fmt.Println(<-ch) // this line never executes
}

func Run2(ch chan int) {
	fmt.Println("Running 2 !!!")
	go func() {
		fmt.Println(<-ch) // receiver goroutine is created and waits here
	}()
	ch <- 10 // sender sends, now that a receiver is waiting in another goroutine
}

func main() {

	ch := make(chan int)

	go Run(ch)
	// go Run2(ch)

	time.Sleep(5 * time.Second)
	fmt.Println("OK I AM DONE !!")
}
