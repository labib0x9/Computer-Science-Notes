package main

import "errors"

type Qeque struct {
	que []int
}

func NewDeque() *Qeque {
	return &Qeque{
		que: make([]int, 0),
	}
}

func (q *Qeque) Len() int {
	return len(q.que)
}

func (q *Qeque) Push(x int) {
	q.que = append(q.que, x) // Add to back
}

func (q *Qeque) Pop() (int, error) {
	if len(q.que) == 0 {
		return 0, errors.New("queue is empty")
	}
	val := q.que[0]   // Take from front
	q.que = q.que[1:] // Remove from front
	return val, nil
}

func main() {

}
