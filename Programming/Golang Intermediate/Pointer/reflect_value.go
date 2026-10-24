package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	var pInt *int = new(int)
	rv := reflect.ValueOf(pInt)

	// ptr
	fmt.Println(rv)
	fmt.Println(rv.Type())
	fmt.Println(rv.Kind())

	// elem -> int
	fmt.Println(rv.Elem())
	fmt.Println(rv.Elem().Type())
	fmt.Println(rv.Elem().Kind())

	/***/

	// do not store uintptr :)
	// rvp := rv.Pointer()
	// fmt.Println(rvp)
	// fmt.Println(unsafe.Pointer(rvp))

	// No, we didn't use stored uintptr
	// newPInt := (*int)(unsafe.Pointer(rvp))
	newPInt := (*int)(unsafe.Pointer(rv.Pointer()))
	*newPInt = 100

	fmt.Println(*pInt)
	fmt.Println(*newPInt)


	/***/

	// fmt.Println(rv.UnsafeAddr())
	// fmt.Println(rv.Elem().UnsafeAddr())
}
