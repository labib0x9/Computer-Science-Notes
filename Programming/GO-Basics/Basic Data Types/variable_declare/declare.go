package main

import "fmt"

func print(arr ...int) {
	fmt.Println(arr)
}

func explicitImplicit() {
	// Explicit
	var a int
	a = 10

	// Implicit
	var b = 10

	print(a, b)

	var name string = "Labib Al Faisal"
	fmt.Println(name)
}

func blockDeclare() {
	// Block Declare
	var (
		c int
		d int
	)

	print(c, d)
}

func shortDeclare() {
	// short declare
	e := "string"
	fmt.Println(e)

	a, b := 10, 20	// a,b both declares and assign
	b, c := 30, 40	// b assign, c declare and assign
	print(a, b, c)
}

func constType() {
	const a = 10
	print(a)

	const (
		low = iota
		medium
		high
	)
	print(low, medium, high)

	const (
		b = iota + 10
		c
		d
	)
	print(b, c, d)
}

func main() {
	explicitImplicit()
	blockDeclare()
	shortDeclare()
	constType()


	var a, b = 10, 20
	sum := a + b
	fmt.Println(sum)

	c, d := b, a
	print(a, b, c, d)
}
