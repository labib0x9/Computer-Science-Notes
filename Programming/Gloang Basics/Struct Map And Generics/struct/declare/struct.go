package main

import "fmt"

type Person struct {
	Name string
	Age int
	FavSongs []string
}

func main() {

	p1 := Person{
		Name: "Labib",
		Age: 12,
		FavSongs: []string{"A", "B"},
	}

	p1.Name = "Faisal"

	fmt.Println(p1)

	c1 := struct {
		Name string
		LastName string
	} {
		Name: "Cat",
		LastName: "Catty",
	}

	fmt.Println(c1)

}