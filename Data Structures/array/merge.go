package main

import "fmt"

func main() {
	/* merge two array */
	arr := []int{1, 3, 4}
	brr := []int{2, 5, 7, 8}

	n, m := len(arr), len(brr)

	N := n + m
	crr := make([]int, N)

	// merge :)
	copy(crr, arr)
	copy(crr[n:], brr)

	fmt.Println(crr)

	/*  merge two sorted array  */
	i, j, idx := 0, 0, 0

	// copy one array entirely
	for i < n && j < m {
		if arr[i] < brr[j] {
			crr[idx] = arr[i]
			i, idx = i+1, idx+1
		} else {
			crr[idx] = brr[j]
			j, idx = j+1, idx+1
		}
	}

	// if arr remains
	for i < n {
		crr[idx] = arr[i]
		i, idx = i+1, idx+1
	}

	// if brr remains
	for j < m {
		crr[idx] = brr[j]
		j, idx = j+1, idx+1
	}

	fmt.Println(crr)
}