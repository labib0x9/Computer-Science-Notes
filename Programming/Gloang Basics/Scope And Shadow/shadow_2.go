package main

import "fmt"

var x = 100

func main() {
	x := 50
	{
		x := 10
		fmt.Println(x)
	}
	fmt.Println(x)
	fmt.Println(global())
}

func global() int {
	return x
}
