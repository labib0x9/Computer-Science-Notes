package main

import "fmt"

type Person struct {
	Name string
	Age int
}

type Employee struct {
	Person	// embeded type
	Position string
	EmployeeId string
}

func main() {
	e := Employee {
		Person: Person{
			Name: "Labib",
			Age: 24,
		},
		Position: "CEO",
		EmployeeId: "ADMIN",
	}

	fmt.Println(e.Name)
}