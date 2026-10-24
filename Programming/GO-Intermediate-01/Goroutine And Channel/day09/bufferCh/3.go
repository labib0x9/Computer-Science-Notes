// package main

// import (
// 	"fmt"
// 	"time"
// )

// func fakeAPI(response chan string) {
// 	// time.Sleep(3 * time.Second)
// 	time.Sleep(1 * time.Second)
// 	response <- "Received"
// }

// func fakeAPI2(response2 chan string) {
// 	time.Sleep(3 * time.Second)
// 	response2 <- "Received2"
// }

// func main() {

// 	response := make(chan string)
// 	response2 := make(chan string)

// 	go fakeAPI(response)
// 	go fakeAPI2(response2)

// 	select {
// 	case res := <-response:
// 		fmt.Println(res)
// 	case res := <-response2:
// 		fmt.Println(res)
// 	case <-time.After(5 * time.Second):
// 		fmt.Println("TimeOut")
// 	}

// }