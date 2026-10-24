package main

import "slices"

func main() {

	s := []int{1, 2, 3, 4, 5, 6}

	slices.BinarySearch(s, 10)
	slices.Clone(s)
	slices.Compare(s, s)
	slices.Sort(s)
	slices.Concat(s, s, s)
	slices.Contains(s, 10)
	slices.Equal(s, s)
	slices.Delete(s, 1, 2)
	slices.Insert(s, 1, 10)
	slices.Compact(s)
	slices.Reverse(s)
}