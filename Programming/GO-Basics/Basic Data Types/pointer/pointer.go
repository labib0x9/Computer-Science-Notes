package main

import (
	"fmt"
)

func main() {

	var x int
	p := &x

	*p = 100 // dereferencing

	var q *int = p
	*q = 200 // also p is changed
	fmt.Println(*p)

	/** **/
	var z *int
	fmt.Println(z) // nil, because z is initialized with zero value of ptr, which is nil

}
