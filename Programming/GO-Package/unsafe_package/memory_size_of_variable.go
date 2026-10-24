package main

import (
	"fmt"
	"unsafe"
)

type y struct {
	a int
	b int64
}

func main() {

	var x uint32
	yy := y{2, 3}

	// size in bytes
	fmt.Println(unsafe.Sizeof(x))
	fmt.Println(unsafe.Sizeof(yy))
}
