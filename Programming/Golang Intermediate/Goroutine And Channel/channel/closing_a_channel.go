package main

/*

labib0x9@TestGoCode $ go run main.go
(true): Hello!!
(true): 99
(false): 0

*/

import "fmt"

func main() {
	stringStream := make(chan string)

	go func() {
		stringStream <- "Hello!!"
	}()

	salutation, ok := <-stringStream
	fmt.Printf("(%v): %v\n", ok, salutation)

	// --- //

	intStream := make(chan int)

	go func() {
		intStream <- 'c'
		close(intStream)
	}()

	value, ok := <-intStream
	fmt.Printf("(%v): %v\n", ok, value)

	value, ok = <-intStream
	fmt.Printf("(%v): %v\n", ok, value)
}
