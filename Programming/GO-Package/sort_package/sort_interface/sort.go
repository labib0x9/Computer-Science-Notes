package main

import "sort"

type Person struct {
	Name string
	Age  int
}

type NewType []Person

// Implement the sort.Interface
func (n NewType) Len() int           { return len(n) }
func (n NewType) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n NewType) Less(i, j int) bool { return n[i].Age < n[j].Age }

func main() {

	people := []Person{
		{"Alice", 30},
		{"Bob", 20},
		{"Charlie", 25},
	}

	sort.Sort(NewType(people))
	sort.Stable(NewType(people))
}
