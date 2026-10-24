package main

import (
	"fmt"
	"unsafe"
)

func main() {

	// what happens here ??
	// type punning..
	var f float64 = 3.1416
	var x uint64 = *(*uint64)(unsafe.Pointer(&f)) // reinterpret the memory as uint64 -> reinterpret_cast<uint64*>

	fmt.Println(x)
}
