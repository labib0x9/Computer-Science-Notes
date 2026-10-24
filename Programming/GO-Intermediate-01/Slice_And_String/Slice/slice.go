package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	sl := make([]int, 10, 20)

	ptr := unsafe.Pointer(&sl) // memory of slice header - 0x1400000c018
	fmt.Println(ptr)

	ptr = unsafe.Pointer(&sl[0]) // memory of backed-array - 0x14000106000
	fmt.Println(ptr)

	/** ---- **/

	sliceHeaderSize := unsafe.Sizeof(sl) // 24-bytes
	fmt.Println(sliceHeaderSize)

	elemSize := unsafe.Sizeof(sl[0])
	fmt.Println(cap(sl) * int(elemSize))

	fmt.Println("Total Memory: ", cap(sl)*int(elemSize)+int(sliceHeaderSize))

	/** --- **/
	ptr = unsafe.Pointer(&sl)
	header := (*reflect.SliceHeader)(ptr)
	fmt.Println(header.Data, header.Len, header.Cap) // 1374390607872 10 20
	fmt.Printf("0x%x\n", header.Data)

	/** --- **/

	tlPtr := unsafe.SliceData(sl) // same as header.Data
	fmt.Println(tlPtr)
	fmt.Println(reflect.TypeOf(tlPtr).Kind())

	slPtr := unsafe.Slice(tlPtr, 20)
	_ = slPtr
	// fmt.Println(slPtr)

}
