package main

import (
	"fmt"
	"unsafe"
)

func main() {

	var x int = 256
	intSize := unsafe.Sizeof(x)

	ptr := (*byte)(unsafe.Pointer(&x)) // why type cast to *byte ??

	fmt.Println(ptr, intSize)
	xByte := unsafe.Slice(ptr, intSize) // int -> []byte

	fmt.Println(xByte)

	x = 123456 // changes xByte..

	fmt.Println(xByte)

	y := *(*int)(unsafe.Pointer(&xByte[0])) // []byte -> int

	fmt.Println(y)
}
