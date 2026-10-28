package main

import "fmt"

/*

Each one has the same probability.

labib0x9@TestGoCode $ go run main.go
338 345 316
*/

func main() {

	c1 := make(chan any)
	close(c1)
	c2 := make(chan any)
	close(c2)
	c3 := make(chan any)
	close(c3)

	var count1, count2, count3 int
	for i := 0; i < 999; i++ {
		select {
		case <-c1:
			count1++
		case <-c2:
			count2++
		case <-c3:
			count3++
		}
	}

	fmt.Println(count1, count2, count3)

}
