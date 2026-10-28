package main

/*

labib0x9@TestGoCode $ go run main.go
1 2 3 4 5 6 7 8 9

*/

import "fmt"

func main() {
	intStream := make(chan int)

	go func() {
		defer close(intStream)
		for i := 1; i < 10; i++ {
			intStream <- i
		}
	}()

	// how does it know that channel is closed ????
	for value := range intStream {
		fmt.Printf("%d ", value)
	}
	fmt.Println()
}
