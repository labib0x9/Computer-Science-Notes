package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Memory of each goroutine

func memConsumed() uint64 {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

func main() {
	var c <-chan any
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c } // dont confuse here, waitgroup count decreases, but the goroutine stays alive..

	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := 0; i < numGoroutines; i++ {
		go noop()
	}
	wg.Wait()
	after := memConsumed()

	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)
}
