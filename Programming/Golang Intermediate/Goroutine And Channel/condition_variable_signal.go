package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal() // only one goroutine wakes up...
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			// c.L.Unlock()	// calls by wait, for visualization
			c.Wait() // suspends this goroutine, calls unlock when entered, calls lock when exits
			// c.L.Lock() // calls by wait, for visualization
		}
		queue = append(queue, struct{}{})
		fmt.Println("Added to queue")
		go removeFromQueue(time.Second)
		c.L.Unlock()
	}
}
