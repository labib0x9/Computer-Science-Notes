package main

import (
	"fmt"
	"reflect"
)

func whatType(t reflect.Type) {
	fmt.Println(t)
	fmt.Println(t.Kind())
}

func main() {

	var x *int
	var y = new(int)

	t := reflect.TypeOf(x)
	whatType(t)

	t = reflect.TypeOf(y)
	whatType(t)

	// So, both x, y is *int

	fmt.Println(x) // nil
	fmt.Println(y) // has value

	// fmt.Println(*x) // nil -> panic in dereferencing
	fmt.Println(*y) // int -> zero value = 0

	/** --------------------- **/

	// make() function only used for slice, map, channel
	// what if we want to create them using new() ?

	/** -------------------- **/

	// slice
	s := new([]int)
	fmt.Println(s, *s, *s == nil)

	// chan
	c := new(chan int)
	fmt.Println(c, *c, *c == nil)

	// map
	m := new(map[int]string)
	fmt.Println(m, *m, m == nil, *m == nil)

	t = reflect.TypeOf(m)
	whatType(t) // kind -> ptr

	/** ---------------------- **/
	mp := make(map[int]string)
	fmt.Println(mp, mp == nil)

	t = reflect.TypeOf(mp)
	whatType(t) // kind -> map
}
