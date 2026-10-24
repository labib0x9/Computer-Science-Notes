package main

type Stack[T any] struct {
	st []T
}

func (s *Stack[T]) Push(val T) {
	s.st = append(s.st, val)
}

func main() {

}