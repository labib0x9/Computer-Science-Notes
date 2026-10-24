// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var counter = 0
// var mu sync.Mutex

// func increment(ch chan bool) {
// 	for i := 0; i < 100; i++ {
// 		mu.Lock()
// 		counter++
// 		mu.Unlock()
// 	}

// 	ch <-true
// }

// func main() {

// 	ch := make(chan bool, 2)

// 	go increment(ch)
// 	go increment(ch)

// 	<-ch
// 	<-ch

// 	fmt.Println(counter)
// }