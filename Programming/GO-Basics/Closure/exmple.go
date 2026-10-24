package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sal := "hello"
	wg.Add(1)

	// closure ....
	go func() {
		defer wg.Done()
		sal = "welcome"
	}()

	wg.Wait()
	fmt.Println(sal)
}
