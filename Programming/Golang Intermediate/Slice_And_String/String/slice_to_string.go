package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	buf := []byte("Hello World")

	// s is a string, string is immutable
	// but here the underlying array is the reference of the underlying array of buf
	// so any changes to buf (within capacity), changes s, breaks string immutability.
	s := unsafe.String(unsafe.SliceData(buf), len(buf))

	fmt.Println(s)	// Hello World

	fmt.Println(reflect.TypeOf(s)) // string

	buf[0] = 'h'

	fmt.Println(s)	// hello World
}
