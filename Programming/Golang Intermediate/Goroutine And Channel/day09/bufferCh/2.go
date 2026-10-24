// package main

// import (
// 	"fmt"
// 	"time"
// )

// // func(ch chan<- int) // Sender
// // func(ch <-chan int) // Receiver
// // func(ch chan int) // Both

// func producer(ch chan<- int) {
// 	for i := 0; i <= 5; i++ {
// 		ch <- i
// 		fmt.Println("Produed: ", i)
// 		// time.Sleep(time.Millisecond * 500)
// 	}
// 	close(ch)
// }

// func consumer(ch <-chan int) {
// 	for data := range ch {
// 		fmt.Println("Consumed: ", data)
// 		time.Sleep(time.Second * 2)
// 	}
// }

// func main() {

// 	// works like asynchonously upto 2 receivers
// 	ch := make(chan int, 2)

// 	go producer(ch)
// 	go consumer(ch)

// 	time.Sleep(time.Second * 15)
// }