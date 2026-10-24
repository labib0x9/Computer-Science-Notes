package main

import "fmt"

// Pipelines -> one goroutines output is another goroutines input

func Pipeline() {

	naturals := make(chan int)
	squares := make(chan int)

	// // Counter
	// go func() {
	// 	for i := 0; i < 20; i++ {
	// 		naturals <- i
	// 	}
	// 	// Close the channel
	// 	close(naturals)
	// }()

	// // Squares
	// go func ()  {
	// 	for x := range naturals {
	// 		squares <- x * x
	// 	}
	// 	close(squares)
	// }()

	// for x := range squares {
	// 	fmt.Println(x)
	// }

	go Counter(naturals)
	go Square(squares, naturals)
	PrintNumber(squares)
}

func Counter(out chan<- int) {
	for i := 0; i < 20; i++ {
		out <- i
	}
	close(out)
}

func Square(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x * x
	}
	close(out)
}

func PrintNumber(in <-chan int) {
	for x := range in {
		fmt.Println(x)
	}
}