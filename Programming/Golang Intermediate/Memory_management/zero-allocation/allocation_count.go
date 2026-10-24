package main

import (
	"fmt"
	"strings"
	"testing"
)

func main() {

	words := []string{"Hello", "world", "\n", "I", "am", "learning", "goooooooooooooooooooooooo"}

	allocCount := testing.AllocsPerRun(1000, func() {
		out := ""
		for _, w := range words {
			out += w
		}
	})
	fmt.Println("Allocations using += :", allocCount)

	allocCount = testing.AllocsPerRun(1000, func() {
		var sb strings.Builder
		sb.Grow(100) // preallocate enough capacity
		for _, w := range words {
			sb.WriteString(w)
		}
		out := sb.String()
		_ = out
	})

	fmt.Println("Allocations using Builder:", allocCount)

}
