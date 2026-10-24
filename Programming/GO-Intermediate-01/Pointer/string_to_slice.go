package main

import (
	"fmt"
	"unsafe"
)

func main() {

	temp := "hello Labib"
	fmt.Println(unsafe.Pointer(&temp))
	fmt.Println(unsafe.StringData(temp))

	temp += " Al"
	fmt.Println(unsafe.Pointer(&temp))
	fmt.Println(unsafe.StringData(temp))

	temp += " Faisal"
	fmt.Println(unsafe.Pointer(&temp))
	fmt.Println(unsafe.StringData(temp))

	fmt.Println(temp)

	// string -> slice
	strData := unsafe.StringData(temp)
	lenn := len(temp)

	fmt.Println(strData, lenn)

	newS := unsafe.Slice(strData, lenn)
	fmt.Println(newS)

	fmt.Println(unsafe.SliceData(newS))
	fmt.Println(temp)
	newS = append(newS, []byte(" , I am learning")...)
	fmt.Println(unsafe.SliceData(newS))
	fmt.Println(temp)

	// slice -> string
	slcData := unsafe.SliceData(newS)
	lenn = len(newS)
	newStr := unsafe.String(slcData, lenn)

	fmt.Println(newStr)
	fmt.Println(slcData)
	fmt.Println(unsafe.StringData(newStr))
}
