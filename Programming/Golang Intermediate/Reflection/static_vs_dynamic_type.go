package main

import (
	"fmt"
	"reflect"
)

/**
	t = reflect.Type
	t.Name() // static name of the type
	t.Kind() // kind of the type (underlying type) - dynamic/concrete type
**/

type MyInt int
type AnotherMyInt MyInt

func main() {

	var x MyInt = 42
	t := reflect.TypeOf(x)
	fmt.Println(t.Name()) // MyInt
	fmt.Println(t.Kind()) // int

	var y AnotherMyInt = AnotherMyInt(x)
	t = reflect.TypeOf(y)
	fmt.Println(t.Name()) // AnotherMyInt
	fmt.Println(t.Kind()) // int

	a := []MyInt{MyInt(1), MyInt(2), MyInt(3)}
	t = reflect.TypeOf(a)

	fmt.Println(t.Name())	// ""
	fmt.Println(t.Kind())	// slice
	fmt.Println(t.Elem().Name())	// MyInt
	fmt.Println(t.Elem().Kind())	// int
}
