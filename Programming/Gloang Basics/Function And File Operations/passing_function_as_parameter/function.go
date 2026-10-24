package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age int
}

func main() {

	people := []Person {
		{Name: "Labib", Age: 25},
		{Name: "Faisal", Age: 23},
		{Name: "Tareq", Age: 15},
	}

	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})

	fmt.Println(people)

}