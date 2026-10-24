package main

import (
	"fmt"
	"reflect"
)

func sum(v ...int) int {
	total := 0
	for i := range v {
		total += v[i]
		fmt.Println(v[i], total)
	}
	return total
}

func whatType(v ...int) {
	t := reflect.TypeOf(v)

	fmt.Println(t)               // []int
	fmt.Println(t.Kind())        // Slice
	fmt.Println(t.Elem().Kind()) // int
}

func main() {

	fmt.Println(sum(1, 2, 3, 4))

	s := []int{1, 2, 3, 4}
	fmt.Println(sum(s...))

	whatType(1, 2, 3, 4)
}
