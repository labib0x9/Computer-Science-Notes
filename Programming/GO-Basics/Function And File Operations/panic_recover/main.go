package main

import "fmt"

func divide(a, b int) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recieved :")
		}
	}()

	if b == 0 {
		panic("Can't divide by zero")
	}

	fmt.Println(a / b)
}

func main() {
	divide(10, 2)
	divide(4, 0)
}
