package main

import "container/heap"

type PriorityQueue []int

func (pq PriorityQueue) Len() int { return len(pq) }

// min heap
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i] < pq[j]
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(int))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func main() {

	pq := &PriorityQueue{3, 2, 1}

	heap.Init(pq)

	heap.Push(pq, 2)
	
	if pq.Len() > 0 {
		heap.Pop(pq)
	}
}
