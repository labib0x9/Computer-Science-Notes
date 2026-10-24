package main

import (
	"fmt"
	"reflect"
)

/**
	At compile time:
	var x interface{}   // static type is interface{}
	x = 42 				// valid because interface{} can hold any type
	x = "Hello, World!" // valid because interface{} can hold any type

	At runtime :
	x = {type: nil, value: nil} 			   // dynamic type is nil
	x = {type: int, value: 42}				   // dynamic type is int
	x = {type: string, value: "Hello, World!"} // dynamic type is string
**/

func main() {

	var x interface{}
	fmt.Println("Type:", reflect.TypeOf(x), "Value:", x) // {nil, nil}

	x = 42
	fmt.Println("Type:", reflect.TypeOf(x), "Value:", x) // {int, 42}

	x = "Hello, World!"
	fmt.Println("Type:", reflect.TypeOf(x), "Value:", x) // {string, Hello, World!}

	fmt.Println("x =", x)
}
