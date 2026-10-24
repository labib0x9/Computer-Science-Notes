// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var counter = 0
// var mu sync.Mutex

// func increment(wg *sync.WaitGroup) {
// 	defer wg.Done()		// wg = wg - 1
// 	for i := 0; i < 100; i++ {
// 		mu.Lock()
// 		counter++
// 		mu.Unlock()
// 	}
// }

// func main() {

// 	var wg sync.WaitGroup	// wg = 0

// 	wg.Add(2)  // wg = wg + 2

// 	go increment(&wg)
// 	go increment(&wg)

// 	wg.Wait()	// wg == 0

// 	fmt.Println(counter)
// }