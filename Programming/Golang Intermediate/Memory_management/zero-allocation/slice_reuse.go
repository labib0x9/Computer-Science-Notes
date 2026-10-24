package main

import (
	"fmt"
	"unsafe"
)

func check(x []int) {
	fmt.Println("len:", len(x), "cap:", cap(x))
}

func main() {

	sl := make([]int, 0)

	ptr := unsafe.Pointer(&sl)
	fmt.Println("Allocates at:", ptr)

	check(sl) // {0, 0}

	for i := 0; i < 20; i++ {
		sl = append(sl, i)
	}

	check(sl) // {20, 32}

	/**
		length decreases, but the capacity remains the same..
		No allocation
	**/
	sl = sl[:0]

	ptr = unsafe.Pointer(&sl)
	fmt.Println("Allocates at:", ptr)

	check(sl) // {0, 32}
}
