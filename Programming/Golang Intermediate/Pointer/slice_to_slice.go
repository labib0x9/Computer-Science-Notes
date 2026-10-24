package main

import (
	"fmt"
	"unsafe"
)

func main() {

	var bigBuf [1024]byte
	copy(bigBuf[:], []byte("hello labib al faisal, how is your learning going ?"))

	l, r := 6, 21
	s1 := bigBuf[l:r]

	fmt.Println(s1)
	fmt.Println(string(s1))

	fmt.Println(unsafe.Pointer(&bigBuf[l]))
	fmt.Println(unsafe.Pointer(&s1[0]))

	/***/

	dataPtr := (*byte)(unsafe.Pointer(&bigBuf[l]))
	s2 := unsafe.Slice(dataPtr, r-l)

	fmt.Println(s2)
	fmt.Println(string(s2))
	fmt.Println(unsafe.SliceData(s2))
}
