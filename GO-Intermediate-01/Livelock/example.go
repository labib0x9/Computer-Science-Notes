package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	codance := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			codance.Broadcast()
		}
	}()

	takeStep := func() {
		codance.L.Lock()
		codance.Wait() // Isn't this blocking ???
		codance.L.Unlock()
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(dir, 1)

		takeStep()

		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprintf(out, " success...")
			return true
		}

		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool {
		return tryDir("left", &left, out)
	}
	tryRight := func(out *bytes.Buffer) bool {
		return tryDir("right", &right, out)
	}

	// --- //

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer walking.Done()
		defer func() { fmt.Println(out.String()) }()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in expere!", name)
	}

	var pepopleInHallway sync.WaitGroup
	pepopleInHallway.Add(2)

	go walk(&pepopleInHallway, "Alice")
	go walk(&pepopleInHallway, "Bob")

	pepopleInHallway.Wait()
}
