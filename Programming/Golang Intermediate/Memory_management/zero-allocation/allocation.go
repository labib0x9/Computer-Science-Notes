package main

import (
	"fmt"
	"unsafe"
)

func main() {

	b := []byte("Hello World")

	ptr := unsafe.Pointer(&b)
	fmt.Println(ptr)

	s := string(b)
	ptr = unsafe.Pointer(&s)
	fmt.Println(ptr)

	/** --- **/

	sl := []int{1, 2, 3}

	ptr = unsafe.Pointer(&sl)
	fmt.Println(ptr)

	sll := sl[:0]
	ptr = unsafe.Pointer(&sll)
	fmt.Println(ptr)

}
