package main

import "fmt"

func worker1() {
	fmt.Println("Wroker 1")
}

func worker2() {
	fmt.Println("Wroker 2")
}

func worker3() {
	fmt.Println("Wroker 3")
}

func worker4() {
	fmt.Println("Wroker 4")
}

func main() {

	defer worker1()
	defer worker2()
	defer worker3()
	defer worker4()

}