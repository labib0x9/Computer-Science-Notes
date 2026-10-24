package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	// var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(I int) {
			defer wg.Done()
			// mu.Lock()
			_, min, sec := time.Now().Clock()
			fmt.Println(min, sec , " :: ", rand.Intn(1000))
			// mu.Unlock()
		}(i)
	}

	wg.Wait()
}