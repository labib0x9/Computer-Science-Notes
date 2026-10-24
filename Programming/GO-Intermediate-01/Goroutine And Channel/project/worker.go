package main

import (
	"fmt"
	"sync"
)

const (
	MAX_QUEUE_ITEM   = 10
	MAX_WORKER_COUNT = 4
)

type Task struct {
	a, b int
}

func execTask(t Task) {
	fmt.Println(t.a, "+", t.b, "=", t.a+t.b)
}

func worker(tch chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tch {
		execTask(t)
	}
}

func main() {

	var wg sync.WaitGroup
	queue := make(chan Task, MAX_QUEUE_ITEM)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go worker(queue, &wg)
	}

	for i := 1; i <= 20; i++ {
		t := Task{
			a: i,
			b: i * i,
		}
		queue <- t
	}

	close(queue)
	wg.Wait()
}
