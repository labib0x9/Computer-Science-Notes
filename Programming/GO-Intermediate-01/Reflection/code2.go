package main

import (
	"fmt"
	"reflect"
)

func Print(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
	fmt.Println(reflect.TypeOf(v).Name())
	fmt.Println(reflect.TypeOf(v).Kind())
	fmt.Println(reflect.TypeOf(v).Elem().Kind())

	fmt.Println("-----")
}

func main() {

	sl := make([]int, 0, 10)
	Print(sl) // []int, "", slice, int

	al := [10]int{}
	Print(al) // [10]int, "", array, int

	mp := make(map[string]int)
	Print(mp) // map[string]int, "", map, int

	ch := make(chan int)
	Print(ch) // chan int, "", chan, int

	// fn := func(a int) int { return a }
	// Print(fn) // func(int) int, "", func, Panic: reflect: Elem of non-pointer type func(int) int

	var p *int
	Print(p) // *int, "", ptr, int

	// var i int
	// Print(i) // int, "int", int, Panic: reflect: Elem of non-pointer type int

	// var i interface{}
	// Print(i) // <nil>, "", invalid, Panic: reflect: Elem of non-pointer type interface

	// var i interface{} = "GEllo"
	// Print(i) // string, "string", string, Panic: reflect: Elem of non-pointer type string
}
