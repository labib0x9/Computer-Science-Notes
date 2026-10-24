package main

import "sort"

func main() {

	sort.Ints([]int{1, 2, 3, 4})

	arr := []int{1, 2, 3, 4, 5, 6, 8, 9}
	sort.Search(len(arr), func(i int) bool {
		return arr[i] >= 7
	})

	people := []struct {
		Name string
		Age  int
	}{
		{"Alice", 30},
		{"Bob", 25},
		{"Eve", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})

	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})

}
