package main

import (
	"fmt"
	"sync"
)

func main() {

	number := 0
	var wg sync.WaitGroup

	// // Using mutex.
	// var mu sync.Mutex

	// for i := 0; i < 10000; i++ {
	// 	wg.Add(1)
	// 	go func(wg *sync.WaitGroup, n *int) {
	// 		defer wg.Done()
	// 		mu.Lock()
	// 		*n++
	// 		mu.Unlock()
	// 	}(&wg, &number)
	// }

	// wg.Wait()
	// fmt.Println(number)

	// Using Confinement
	req := make(chan struct{})
	done := make(chan bool)

	go func(req <-chan struct{}, done <-chan bool, n *int) {
		for {
			select {
			case <-req:
				*n++
			case <-done:
				return
			}
		}
	}(req, done, &number)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			req <- struct{}{}
		}(&wg)
	}

	wg.Wait()
	done <- true

	fmt.Println(number)
}
