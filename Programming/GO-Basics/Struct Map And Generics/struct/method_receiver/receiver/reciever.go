package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func methodReciever() {

	var a User = User{
		Name: "Labib",
		Age:  100,
	}
	a.SayHello()

	// Method chaining
	b := &User{}
	fmt.Println(b)
	b.SetName("Faisal").SetAge(10).SayHello()
	fmt.Println(b)
}

// Received by value
func (u User) SayHello() {
	fmt.Println("hello", u.Name)
}

// Received by reference
func (u *User) SayHelloAgain() {
	fmt.Println("hello", u.Name)
}

// Receiver is a pointer (u *USer)
// Return type is a pointer *User
func (u *User) SetName(name string) *User {
	u.Name = name
	return u
}

func (u *User) SetAge(age int) *User {
	u.Age = age
	return u
}
