package main

import "fmt"

// Only sender chan
func SendMessage(ch chan<- string, n int) {
	for i := 0; i < 5; i++ {
		ch <- fmt.Sprintf("Worker-%d : %d", n, i)
	}
	close(ch)
}

// Only receiver chan
func PrintMessage(ch <-chan string) {
	for msg := range ch {
		fmt.Println(msg)
	}
}

// Sender and receiver chan
// Another catch here : if the SendMessage() function finishes before Send() function.
// It will try to send data to a closed channel, which may cause panic / issue
func Send(ch chan string) {
	ch <- "Sender and receiver channel"
}

func main() {

	ch := make(chan string)

	go Send(ch)
	go SendMessage(ch, 1)
	PrintMessage(ch)
}
