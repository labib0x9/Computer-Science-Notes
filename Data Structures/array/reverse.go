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

	// two pointer reverse
	lo, hi := 0, n-1
	for lo < hi {
		arr[lo] ^= arr[hi]
		arr[hi] ^= arr[lo]
		arr[lo] ^= arr[hi]

		lo++
		hi--
	}

	fmt.Println(arr)
}
