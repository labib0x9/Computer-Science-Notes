package main

import (
	"fmt"
	"reflect"
)

func main() {

	var x int
	p := &x
	t := reflect.TypeOf(p)

	fmt.Println(t)               // *int
	fmt.Println(t.Kind())        // ptr
	fmt.Println(t.Elem().Kind()) // int

}
