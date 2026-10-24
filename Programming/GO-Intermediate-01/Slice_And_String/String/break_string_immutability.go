package main

import (
	"fmt"
	"unsafe"
)

func main() {

	// s := "hello"	// read-only memory, Results SIGBUS/SIGSEGV
	s := string([]byte("hello")) // heap

	// s[0] = 'K' // No permission

	sPtr := unsafe.StringData(s) // pointer to backed-array of s
	fmt.Println(sPtr)

	b := unsafe.Slice(sPtr, len(s))
	fmt.Println(unsafe.Pointer(&b))

	/** -- **/
	bPtr := unsafe.SliceData(b) // pointer to backed-array of b
	fmt.Println(bPtr)

	b[0] = 'H' // changed s

	fmt.Println(s)
	fmt.Println(b)
}
