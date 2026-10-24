package main

import (
	"fmt"
	"reflect"
)

type MyInt int
type AnotherMyInt MyInt

func Add(a, b interface{}) interface{} {
	vA := reflect.ValueOf(a)
	vB := reflect.ValueOf(b)
	if vA.Kind() != vB.Kind() {
		return fmt.Errorf("type mismatch: a=%s, b=%s", vA.Kind(), vB.Kind())
	}

	switch vA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vA.Int() + vB.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vA.Uint() + vB.Uint()
	case reflect.Float32, reflect.Float64:
		return vA.Float() + vB.Float()
	case reflect.String:
		return vA.String() + vB.String()
	default:
		return fmt.Errorf("unsupported type: %s", vA.Kind())
	}
}

func main() {

	var x MyInt = 42
	var y AnotherMyInt = 40

	fmt.Println(Add(x, y))
	fmt.Println(Add(12.5, 3.4))
	fmt.Println(Add("Hello, ", "World!"))
	fmt.Println(Add(x, "SSS"))
}
