package main

import (
	"fmt"
	"reflect"
)

type MyInt int

func main() {

	var x MyInt = 42
	v := reflect.ValueOf(x)

	fmt.Println(v.Kind() == reflect.Int)    // true
	fmt.Println(v.Type().Name() == "MyInt") // true

	fmt.Println(v.Int()) // 42

	fmt.Println("Can interface:", v.CanInterface()) // true

	i := v.Interface()          // i is of type interface{} and holds a MyInt
	fmt.Printf("%T %v\n", i, i) // main.MyInt 42

	// // panic: interface conversion: interface {} is main.MyInt, not int
	// f := i.(int)

	f := i.(MyInt) // Is i a MyInt?
	fmt.Println(f)

	fmt.Println(reflect.TypeOf(i).Kind()) // int

	fmt.Println(v.CanSet()) // false
	// v.SetInt(100) // panic: reflect: reflect.Value.SetInt using unaddressable value

	v = reflect.ValueOf(&x)
	fmt.Println(v.CanSet()) // false
	v = v.Elem()
	fmt.Println(v.CanSet()) // true
	v.SetInt(100)

	// We changed the value of x through reflection
	fmt.Println(x) // 100

	var y int = 42
	v = reflect.ValueOf(y)
	// _ = v.Elem() // panic: reflect: call of reflect.Value.Elem on int Value
	fmt.Println(v.CanSet()) // false
	fmt.Println(v.Type())   // int

	v = reflect.ValueOf(&y)
	fmt.Println(v.CanSet())                     // false
	fmt.Println(v.Type())                       // *int
	fmt.Println(v.Type().Kind() == reflect.Ptr) // true

	v = v.Elem()               // dereference
	fmt.Println(v.Type())      // int
	fmt.Println(v.Interface()) // 42

	fmt.Println(v.Addr())  // 0xc0000180b8
	fmt.Printf("%p\n", &y) // 0xc0000180b8

	fmt.Println(v.CanSet()) // true
	v.SetInt(100)           // set through reflection, y = 100 now
	fmt.Println(y)          // 100
}
