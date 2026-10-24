// package main

// import (
// 	"fmt"
// 	"sync"
// )

/** What is the problem here ?
race condition
**/

// var counter = 0

// func increment(wg *sync.WaitGroup) {
// 	defer wg.Done()		// wg = wg - 1
// 	for i := 0; i < 100; i++ {
// 		counter++ // Race
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