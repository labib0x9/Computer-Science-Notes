package main

import (
	"fmt"
	"time"
)

func main() {
	chanOwner := func(done <-chan struct{}, stream <-chan int) <-chan struct{} {
		complete := make(chan struct{})
		go func() {
			defer close(complete)
			for {
				select {
				case <-done: // Tells the goroutine to exit, leak prevented
					return
				case val := <-stream: // blocks on read operation....
					fmt.Println(val)
				}
			}
		}()
		return complete
	}

	done := make(chan struct{})
	complete := chanOwner(done, nil)
	go func() {
		time.Sleep(10 * time.Second)
		close(done) /// HEERE <-
	}()
	<-complete
}
