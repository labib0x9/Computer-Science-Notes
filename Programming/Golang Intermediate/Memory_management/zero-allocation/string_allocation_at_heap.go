package main

import (
	"fmt"
	"unsafe"
)

/**
go build -gcflags="-m" main.go

./main.go:20:2: moved to heap: out
./main.go:26:9: moved to heap: w

**/

func main() {

	words := []string{"Hello", "world", "\n", "I", "am", "learning", "goooooooooooooooooooooooo"}

	out := ""

	ptr := unsafe.Pointer(&out)
	fmt.Println("Allocation at:", ptr)

	// heap allocation for w
	for _, w := range words {
		out += w
		ptr = unsafe.Pointer(&w)
		fmt.Println("Allocation at:", ptr)
	}

	out = ""
	for i := range words {
		out += words[i]
		fmt.Println(cap([]byte(out)), len(out)) // cap == len
	}

}
