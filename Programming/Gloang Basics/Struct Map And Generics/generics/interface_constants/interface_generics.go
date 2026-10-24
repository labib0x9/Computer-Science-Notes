package main

type Adder interface {
	~int | ~float32
}

func Add[T Adder](a, b T) T {
	return a + b
}