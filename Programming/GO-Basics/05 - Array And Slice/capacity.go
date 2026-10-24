package main

import "fmt"

func main() {
	// len = 8, cap = 8
	a1 := []int{1, 2, 3, 4, 5, 6, 7, 8}

	a := a1[3:5]
	fmt.Println(len(a), cap(a)) // 2 5

	a = a1[3:5:8]
	fmt.Println(len(a), cap(a)) // 2 5

	a = a1[3:5:5]
	fmt.Println(len(a), cap(a)) // 2 2
}
