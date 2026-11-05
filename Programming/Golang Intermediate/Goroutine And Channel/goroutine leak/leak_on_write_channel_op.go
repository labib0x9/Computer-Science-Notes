package main

import (
	"fmt"
	"math/rand"
)

/*
abib0x9@TestGoCode $ go run main.go
7291885921076212459
3854131851901820103
8004090447373233516
*/

func main() {
	chanOwner := func() <-chan int {
		intStream := make(chan int)
		go func() {
			defer fmt.Println("CLOSING...")
			defer close(intStream)
			for { // no channel close
				intStream <- rand.Int()
			}
		}()
		return intStream
	}

	randReader := chanOwner()
	for i := 0; i < 3; i++ { // reads only three., but goroutine continues in background
		fmt.Println(<-randReader)
	}
}
