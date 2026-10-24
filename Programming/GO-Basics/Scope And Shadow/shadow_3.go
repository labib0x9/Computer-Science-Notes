package main

import "fmt"

func foo(x int) {
	fmt.Println(x)
	{
		x := x + 10
		fmt.Println(x)
	}
	fmt.Println(x)
}

func main() {
	foo(5)
}
