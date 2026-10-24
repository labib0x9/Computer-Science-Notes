package main

import "fmt"

func main() {
	
	// basic datatypes
	var name string
	var counter int
	var flag bool
	fmt.Println(name, counter, flag)

	// pointer
	var pointer *int
	fmt.Println(pointer)

	// slice
	var arr []string
	brr := make([]string, 0)
	fmt.Println(arr, brr)

	// struct
	var Struct = struct {
		a int
		s string
	}{}
	fmt.Println(Struct)

	// map
	var seen map[string]bool
	visited := make(map[string]bool)
	fmt.Println(seen, visited)

	// channel
	ch := make(chan int)
	fmt.Println(ch)

	// function
	var fun = func () {}
	fmt.Println(fun)
}