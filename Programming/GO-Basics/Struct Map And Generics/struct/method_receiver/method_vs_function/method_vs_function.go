
package main

import "fmt"

// xyz custom data type
type xyz int

func Blabla(){

	var a xyz
	a = 10

	a.SayHelloAgain()
}

// Function
func SayHello(){
	fmt.Println("Gopher!!")
}

// Method
func (b xyz) SayHelloAgain(){
	fmt.Println("GOOGO")
}