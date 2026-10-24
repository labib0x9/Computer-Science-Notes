package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// What does this code do ?? What is happening here ??
// Open all the files in goroutine and find max of these numbers

/**
But what happens if i try to open 100000 files ??
- a lots of goroutine, the system curshes ...
**/

func main() {
	files := []string{"a.txt", "b.txt", "c.txt"}
	// ch := make(chan int, len(files))
	ch := make(chan int)

	// Creates three goroutine and ends the loop
	for _, filename := range files {
		go findMax(filename, ch)
	}

	// Waiting for goroutines to send some result
	for i := 0; i < len(files); i++ {
		fmt.Println(<-ch)
	}
}

func findMax(filename string, ch chan int) {
	f, err := os.Open(filename)
	if err != nil {
		ch <- 0
		return
	}
	defer f.Close()

	arr := make(chan int)
	maxx, line := 0, 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		go convert(scanner.Text(), arr)
		line++
	}

	for i := 0; i < line; i++ {
		x := <-arr
		maxx = max(maxx, x)
	}

	ch <- maxx
}

func convert(value string, arr chan int) {
	value = strings.TrimSpace(value)
	x, err := strconv.Atoi(value)
	if err != nil {
		arr <- 0
		return
	}
	arr <- x
}
