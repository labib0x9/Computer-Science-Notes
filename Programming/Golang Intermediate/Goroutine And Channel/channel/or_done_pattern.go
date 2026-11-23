package main

import (
	"fmt"
	"time"
)

/*

labib0x9@TestGoCode $ go run main.go
done after 1.001201875s

**/

func main() {

	// think like an array of int...

	// takes a slice of channels
	// returns a channel

	var or func(channels ...<-chan any) <-chan any

	or = func(channels ...<-chan any) <-chan any {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan any)

		// recursive calll....
		go func() {
			defer close(orDone)

			// you can add more, if you want.. every call spawns a new goroutine..
			// and each goroutine blocks on select for listen..
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(time.Second),
		sig(time.Hour),
		sig(3*time.Hour))
	fmt.Println("done after", time.Since(start))
}
