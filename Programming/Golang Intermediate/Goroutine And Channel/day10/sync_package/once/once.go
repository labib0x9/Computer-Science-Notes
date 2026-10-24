package main

import (
	"fmt"
	"sync"
)

// Instance must be used once, no matter how many goroutines calls it

func Init() {
	fmt.Println("Init ")
}

func main() {

	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int, wg *sync.WaitGroup) {
			defer wg.Done()
			once.Do(Init)
			fmt.Printf("Goroutine %d running\n", n)
		}(i, &wg)
	}

	wg.Wait()
	fmt.Println("Done !!")
}
