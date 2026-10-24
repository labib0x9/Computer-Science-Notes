package main

import "errors"

type Deque struct {
	dq []int
}

func NewDeque() *Deque {
	return &Deque{
		dq: make([]int, 0),
	}
}

func (d *Deque) Len() int {
	return len(d.dq)
}

func (d *Deque) PushFront(x int) {
	d.dq = append([]int{x}, d.dq...)
}

func (d *Deque) PushBack(x int) {
	d.dq = append(d.dq, x)
}

func (d *Deque) PopFront() (int, error) {
	if len(d.dq) == 0 {
		return 0, errors.New("deque is empty")
	}
	val := d.dq[0]
	d.dq = d.dq[1:]
	return val, nil
}

func (d *Deque) PopBack() (int, error) {
	if len(d.dq) == 0 {
		return 0, errors.New("deque is empty")
	}
	val := d.dq[len(d.dq)-1]
	d.dq = d.dq[:len(d.dq)-1]
	return val, nil
}

func main() {

}
