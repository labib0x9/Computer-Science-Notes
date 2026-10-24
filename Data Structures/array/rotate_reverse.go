package main

import "fmt"

/*
	Rotate an array by k using reverse technique.

	TC : O(n)
	MC : O(1)
*/

func reverse_array(arr []int, lo int, hi int) {
	for lo < hi {
		arr[lo] ^= arr[hi]
		arr[hi] ^= arr[lo]
		arr[lo] ^= arr[hi]
		lo, hi = lo+1, hi-1
	}
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	n, k := len(arr), 3
	// Left - Rotatioin
	reverse_array(arr, 0, k-1)
	reverse_array(arr, k, n-1)
	reverse_array(arr, 0, n-1)
	fmt.Println(arr)

	// Right - Rotation
	reverse_array(arr, 0, n - k - 1)
	reverse_array(arr, n - k, n - 1)
	reverse_array(arr, 0, n - 1)
	fmt.Println(arr)
}
