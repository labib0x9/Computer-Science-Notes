package main

import "fmt"

func main() {
	// insertion, n elements, n + 1 room
	n := 7
	arr := make([]int, n+1)

	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}

	fmt.Println(arr)

	idx, x := 0, 100
	copy(arr[idx+1:], arr[idx:]) // right shift
	arr[idx] = x                 // insert

	fmt.Println(arr)
}
