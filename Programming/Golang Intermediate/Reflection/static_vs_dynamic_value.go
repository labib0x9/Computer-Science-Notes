package main

import (
	"fmt"
	"reflect"
)

/**
	t = reflect.Value
	t.Type() = reflect.Type // static type
	t.Kind() = reflect.Kind // dynamic type
**/

type MyInt int
type AnotherMyInt MyInt

func main() {

	var x MyInt = 42
	v := reflect.ValueOf(x)
	fmt.Println(v)
	fmt.Println(v.Type())	// MyInt
	fmt.Println(v.Kind()) 	// int

	var y AnotherMyInt = AnotherMyInt(x)
	v = reflect.ValueOf(y)
	fmt.Println(v)
	fmt.Println(v.Type())	// AnotherMyInt
	fmt.Println(v.Kind())	// int

}
