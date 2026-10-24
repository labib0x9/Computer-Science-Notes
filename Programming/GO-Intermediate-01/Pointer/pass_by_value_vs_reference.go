package main

import (
	"fmt"
	"unsafe"
)

type A struct {
	v1 int
	v2 string
}

func PassByValue(v A) {
	ptr := unsafe.Pointer(&v)
	fmt.Println("PassByValue=", ptr)
}

func PassByReference(v *A) {
	ptr := unsafe.Pointer(v)
	fmt.Println("PassByReference=", ptr)
}

func main() {

	v := A{v1: 10, v2: "key"}

	ptr := unsafe.Pointer(&v)
	fmt.Println("main=", ptr)

	PassByValue(v)
	PassByReference(&v)
}
