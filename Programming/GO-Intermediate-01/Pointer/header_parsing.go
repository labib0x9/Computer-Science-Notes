package main

import (
	"fmt"
	"unsafe"
)

type Header struct {
	src uint16 // 2bytes
	dst uint16 // 2bytes
	tmp uint16 // 2bytes
}

func main() {

	data := []byte{
		0x12, 0x34, // 13330
		0x34, 0x12, // 4660
		0x52, 0x18, // -> why xByte to uint16 is 6226 ??? because of endianess
	}

	fmt.Println(len(data), unsafe.Sizeof(Header{})) // 6 6

	hdr := (*Header)(unsafe.Pointer(&data[0]))
	fmt.Println(hdr)

	/***/
	x := uint16(4660)

	ptr := (*byte)(unsafe.Pointer(&x))
	// xByte := unsafe.Slice(ptr, unsafe.Sizeof(x))
	xByte := unsafe.Slice(ptr, 2)

	fmt.Println(xByte) // [52, 18]

	y := *(*uint16)(unsafe.Pointer(&xByte[0])) // 4660
	fmt.Println(y)

	// reverse 2byte xbyte
	xByte[0], xByte[1] = xByte[1], xByte[0]   // [18, 52]
	y = *(*uint16)(unsafe.Pointer(&xByte[0])) // 13330
	fmt.Println(y)
}
