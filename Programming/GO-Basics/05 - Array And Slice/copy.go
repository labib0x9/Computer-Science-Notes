package main

import "fmt"

func PrintSlice(s []int) {
	fmt.Println(s)
}

func main() {

	/** --- **/

	sl1 := []int{1, 2, 3, 4, 5}
	sl2 := make([]int, 3)

	// sl2 is a copy of sl1, they are not backed by the same array
	copy(sl2, sl1)

	sl1[0] = 1000
	PrintSlice(sl2)

	/** --- **/

	// copy function -> copy(dst, src)
	// Remove range [2, 3] zero based index
	s1 := []int{0, 1, 2, 3, 4, 5}
	PrintSlice(s1[2:])
	PrintSlice(s1[4:])
	copy(s1[2:], s1[4:])
	PrintSlice(s1[:len(s1)-(4-2)])

	// Insert at index i = 3
	s2 := []int{0, 1, 2, 4, 5, 6}
	s2 = append(s2, -1)
	PrintSlice(s2[4:])
	PrintSlice(s2[3:])
	copy(s2[4:], s2[3:])
	s2[3] = 3
	PrintSlice(s2)

	// clear a slice without reallocation
	copy(s2, s2[:0])
	PrintSlice(s2)

}
