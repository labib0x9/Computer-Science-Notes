package main

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

func main() {

	var wg sync.WaitGroup
	printSum := func(id string, v1, v2 value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(20 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Println("ID=", id, " | a + b =", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)

	// copies lock value :), no deadlock... Why is that ??

	go printSum("a", a, b)
	go printSum("b", b, a)

	wg.Wait()
}
