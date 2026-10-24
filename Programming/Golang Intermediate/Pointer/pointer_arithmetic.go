package main

import (
	"fmt"
	"unsafe"
)

func main() {

	// pointer arithmetic
	arr := [5]int{1, 2, 3, 4, 5}

	ptr := &arr[0]
	elemSize := unsafe.Sizeof(arr[0])

	fmt.Println(ptr, elemSize)

	// now := uintptr(unsafe.Pointer(ptr)) // gc doesn't track, can be deallocate

	index := 0
	arr_0 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(index)*elemSize))

	fmt.Println(*arr_0)

	index = 3
	arr_3 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(index)*elemSize))

	fmt.Println(*arr_3)
}
