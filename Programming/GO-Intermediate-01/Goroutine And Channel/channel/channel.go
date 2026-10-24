package main

import "fmt"

func Channel() {

	// channel
	// Only string types data's are sent to goroutine
	// unbuffered channel
	ch := make(chan string)

	// ch = make(chan string, 0) // unbuffered channel
	// ch = make(chan string, 3) // capacity 3, buffered channel

	// runs a goroutine
	go func(ch chan string) {
		ch <- "Completed"
	}(ch)

	// waits for background goroutine to finish
	ok := <-ch
	fmt.Println(ok)

	// If we just want to finish the goroutine
	done := make(chan struct{})
	go func() {
		fmt.Println("Running goroutine")
		done <- struct{}{}
	}()
	<-done
}

func CHAN() {

	// Boom Deadlock Why ???
	//  See send_receive_issue code in unbuffer_channel folder
	// The main goroutine is trying to send data, before the receiver is ready
	ch := make(chan string)
	ch <- "String"
	<-ch

	// NO DEADLOCK
	ch1 := make(chan string, 2)
	ch1 <- "STRing"
	<-ch1

	// BOOM AGAIN DEADLOCK, wHy ??
	ch2 := make(chan string, 2)
	ch2 <- "STRing"
	ch2 <- "STR"
	ch2 <- "SGG"
	<-ch2
}
