package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
labib0x9@TestGoCode $ go run main.go
5353358849993921561
2116591911802355215
1384995958192116083
CLOSING...

*/

func main() {
	chanOwner := func(done <-chan struct{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer fmt.Println("CLOSING...")
			defer close(intStream)
			for {
				select {
				case intStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return intStream
	}

	done := make(chan struct{})
	randReader := chanOwner(done)
	for i := 0; i < 3; i++ {
		fmt.Println(<-randReader)
	}
	close(done) // signal's the routine for cancellation..

	// spawns and terminate goroutine. same goroutine is responsible..

	time.Sleep(10 * time.Second)
}
