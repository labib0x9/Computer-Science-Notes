package main

import (
	"fmt"
	"time"
)

// main goroutine = parent
// test() creates a goroutine from main goroutine
// parent dies -> spawn goroutine dies

func test() {
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("GOROUTINE")
	}()
}

func main() {
	test()
	fmt.Println("MAIN")
	time.Sleep(30 * time.Second)
}

// OUTPUT: MAIN, GOROUTINE
