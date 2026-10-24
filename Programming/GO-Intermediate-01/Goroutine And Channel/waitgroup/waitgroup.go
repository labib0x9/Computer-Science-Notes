package main

import (
	"fmt"
	"sync"
)

func Print(i int, wg *sync.WaitGroup) {
	defer wg.Done() // mark goroutine as done
	fmt.Println(i, i*i)
}

func main() {

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1) // add one goroutine
		go Print(i, &wg)
	}

	// wait for all goroutine to finish
	wg.Wait()
}
