package main

import (
	"strings"
)

func main() {

	// Slice reuse
	s := make([]int, 10, 10000)
	s = s[:0]

	// sync.Pool

	// strings.Builder
	var x strings.Builder
	x.Grow(10000)
	x.WriteString("Hello")
}
