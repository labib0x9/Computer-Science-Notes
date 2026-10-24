package main

import (
	"fmt"
	"sync"
)

// Only Sender channel
func SendMessage(ch chan<- string, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 3; i++ {
		ch <- fmt.Sprintf("Worker-%d : %d", n, i)
	}
}

// Only Receiver channel
func PrintMessage(ch <-chan string) {
	for msg := range ch {
		fmt.Println(msg)
	}
}

func main() {

	ch := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go SendMessage(ch, 1, &wg)
	go SendMessage(ch, 2, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	PrintMessage(ch)

}
