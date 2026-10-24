package main

import "fmt"

func main() {
	// n elements
	n := 7
	arr := make([]int, n)

	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}

	fmt.Println(arr)

	// deletion
	idx := 0
	copy(arr[idx:], arr[idx+1:]) // left shift
	arr = arr[:n-1]              // remove one room

	fmt.Println(arr)
}
