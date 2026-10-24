package main

import (
	"fmt"
	"unsafe"
)

type s struct {
	f float32 // 4 + padding(4) = 8bytes
	x int     // 8 bytes
}

func main() {

	v := s{3.1416, 2}

	vPtr := unsafe.Pointer(&v) // starting addr of v
	fPtr := unsafe.Pointer(&v.f)
	xPtr := unsafe.Pointer(&v.x)

	fmt.Println(vPtr) // vPtr == fPtr
	fmt.Println(fPtr)
	fmt.Println(xPtr)

	vSize := unsafe.Sizeof(v)
	fSize := unsafe.Sizeof(v.f)
	xSize := unsafe.Sizeof(v.x)

	fmt.Println(vSize)
	fmt.Println(fSize)
	fmt.Println(xSize)

	fOffset := unsafe.Offsetof(v.f)
	xOffset := unsafe.Offsetof(v.x)

	fmt.Println(fOffset) // 0
	fmt.Println(xOffset) // 0 + 8 (float32 size) = 8

	fAddr := unsafe.Pointer(uintptr(vPtr) + fOffset)
	xAddr := unsafe.Pointer(uintptr(vPtr) + xOffset)

	fmt.Println(fAddr)
	fmt.Println(xAddr)

	vPtrEnd := unsafe.Pointer(uintptr(vPtr) + vSize)
	fmt.Println(vPtrEnd)

	/****/

}
