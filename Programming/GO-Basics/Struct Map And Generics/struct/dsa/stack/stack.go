package main

import (
	"errors"
	"fmt"
)

type Stack struct {
	st []int
}

func NewStack() *Stack {
	return &Stack{st: make([]int, 0)}
}

func (s *Stack) Empty() bool {
	return len(s.st) == 0
}

func (s *Stack) Push(x int) {
	s.st = append(s.st, x)
}

func (s *Stack) Top() (int, error) {
	if len(s.st) == 0 {
		return 0, errors.New("error: empty stack")
	}
	return s.st[len(s.st) - 1], nil
}

func (s *Stack) Pop() (x int, err error) {
	if len(s.st) == 0 {
		err = errors.New("error: empty stack")
		return
	}
	x = s.st[len(s.st) - 1]
	s.st = s.st[:len(s.st) - 1]
	return
}

func (s *Stack) String() string {
	return fmt.Sprintf("%v", s.st)
}

func main() {

	st := NewStack()

	st.Push(10)

	_, err := st.Top()
	x, err := st.Pop()

	if err != nil {

	}
	fmt.Println(x)

	st.Push(10)
	st.Push(20)

	fmt.Println(st.String())

}