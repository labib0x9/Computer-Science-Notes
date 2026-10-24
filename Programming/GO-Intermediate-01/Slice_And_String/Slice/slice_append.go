package main

import (
	"fmt"
	"unsafe"
)

// Need to check again, misinformation may occurs here

func main() {

	s := make([]int, 1)

	fmt.Println()
	// slice-header on stack
	ptr := unsafe.Pointer(&s) // 0x1400000c018
	fmt.Println("Slice Header:", ptr)

	// Backed-array on heap
	bPtr := unsafe.SliceData(s) // 0x1400000e0b8
	fmt.Println("First Elem  :", bPtr)

	s = append(s, 10)
	fmt.Println()

	// slice-header on stack
	ptr = unsafe.Pointer(&s) // 0x1400000c018
	fmt.Println("Slice Header:", ptr)

	// new backed-array on heap
	aPtr := unsafe.SliceData(s) // 0x1400000e0d0
	fmt.Println("First Elem  :", aPtr)

	// How far the new allocation happend
	fmt.Println("diff        :", uintptr(unsafe.Pointer(aPtr))-uintptr(unsafe.Pointer(bPtr)), "bytes")
}
