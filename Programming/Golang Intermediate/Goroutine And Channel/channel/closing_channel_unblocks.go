package main

import "time"

func main() {
	ch := make(chan struct{})

	go func() {
		time.Sleep(10 * time.Second)
		close(ch) // closing the channel wakes waiting/blocking
	}()
	<-ch
}
