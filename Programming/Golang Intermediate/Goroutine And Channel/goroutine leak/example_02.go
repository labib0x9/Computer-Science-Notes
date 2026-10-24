package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var mu sync.Mutex
	mu.Lock()

	go func() {
		fmt.Println("worker: waiting for lock...")
		mu.Lock()
		fmt.Println("worker: got lock")
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println("main: still alive")
	}
}
