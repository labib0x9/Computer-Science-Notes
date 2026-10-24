package main

import "fmt"

type Teacher struct {
	Name string
	Age  int
	Id string
}

// Receiver type is pointer
// Return type is pointer
func (t *Teacher) setName(Name string) *Teacher {
	t.Name = Name
	return t
}

// Receiver type is value
// Return type is value
func (t Teacher) setAge(Age int) Teacher {
	t.Age = Age
	return t
}

// func (t *Teacher) setId(Id string) *Teacher {
// 	t.Id = Id
// 	return t
// }

// // Receiver type and pointer type should be same, else there is error
// // Receiver type is value
// // Return type is pointer of copied value 
func (t Teacher) setId(Id string) *Teacher {
	t.Id = Id
	return &t
}

func main() {

	t1 := Teacher{
		Name: "Labib",
		Age: 19,
		Id: "123",
	}

	fmt.Println(t1)

	// Name updates as it is pointer
	// Age doesn't update as it is passed as value
	// Id doesn't update as it's receiver type is value
	t1.setName("Faisal").setAge(30).setId("321")

	// // Error
	// // setName reciever is a pointer
	// // setAge return is a value
	// t1.setAge(231).setName("Labib").setId("32r2")


	// setId returns pointer, which works as a reciever in setName
	t1.setId("365").setName("FAISAL")

	fmt.Println(t1)
}
