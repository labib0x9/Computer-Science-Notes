package main

import "fmt"

func main() {
	// through variables and short declaration
	var s1 []int  // nil slice
	s2 := []int{} // empty slice

	// make(type, length, capacity)
	s3 := make([]int, 0)
	s4 := make([]int, 0, 10)

	// array slicing -> [low:high:capacity]
	arr := []int{1, 2, 4, 5, 6, 7, 8, 9, 10}
	s5 := arr[2:4]
	s6 := arr[3:5:5]
}
