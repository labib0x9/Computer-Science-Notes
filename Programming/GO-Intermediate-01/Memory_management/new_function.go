package main

import "fmt"

func EscapeToHeap() {
	n := new(int)
	fmt.Println(n) // escape
}

func NoEscapeToHeap() {
	m := new(int)
	_ = m
}

func main() {

	EscapeToHeap()
	NoEscapeToHeap()
}
