package main

import (
	"errors"
	"fmt"
)

func division(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("error: can't divide by zero")
	}
	return a / b, nil
}

func add(a int, b int) (sum int) {
	sum = a + b
	return 
}

func throwError() error {
	err := fmt.Errorf("Error 2")
	err = errors.New("Error 3")
	return err
}

func main() {
	err := throwError()
	if err == nil {
		fmt.Println("No error")
	}

	div, _ := division(10, 20)
	fmt.Println(div)

	fmt.Println(add(10, 20))

	mult := func (a, b int) int {
		return a * b
	}
	fmt.Println(mult(10, 2))
}