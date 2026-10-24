package main

import "fmt"

func forLoop() {

	for i := 0; i <= 15; i++ {
		if i % 5 == 0 {
			continue
		}
		if i & 1 == 1 {
			fmt.Println("Odd")
		} else {
			fmt.Println("Even")
		}
	}


	for index, value := range []string{"a", "b", "c"} {
		fmt.Println(index, ":", value)
	}
}

func infiniteLoop() {
	for {
		// do operation
		break
	}
}

func main() {

	forLoop()
}