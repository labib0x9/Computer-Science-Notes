package main

import "fmt"

func main() {
	a1 := [...]int{} // This is Array of zero size

	var a2 [3]int

	days := [...]string{
		1 : "Saturday",
		2 : "Sunday",
		3 : "Friday",
	}
	fmt.Println(days[1])

	arr := [5]int{0, 1, 2, 3, 4}

	fmt.Println(a1, a2, arr)

	var x [][]int
	fmt.Println(x)
}
