package main

import (
	"fmt"
	"sync"
)

// Shared queue -> fixed sized worker pool
// NOTE: Distribution isn't fair.
// That means one worker may run more often than others,
// e.g., worker-0 might process 10 tasks while worker-1 and worker-2 handle fewer (2 tasks).

/**
	Worker Pool Architecture:

			[Main Threat]
				|
				|	[Enqueue tasks]
				v
			[Job Queue] [Pool]
				|
			 /  |  \
			W1	W2  W3
			 \	|  /	[Send Result]
				|
		[Result Channel]
				|
				|
		  [WaitGroup sync]
				|
		  [Process Output]

**/

// pool starts the workers and distributes the tasks to them via a shared job queue.
func pool(numWorker int, task []int, result chan string, wg *sync.WaitGroup) {
	job := make(chan int, len(task)) // Shared job queue (buffered to length of tasks)

	// Create numWorker goroutines, each acting as a worker.
	for i := 0; i < numWorker; i++ {
		wg.Add(1) // Track each worker's completion
		go worker(i, job, result, wg)
	}

	// Submit all tasks into the shared job queue
	for _, j := range task {
		job <- j
	}
	close(job) // Signal to workers: no more tasks
}

// worker continuously pulls tasks from the job queue and sends results to the result channel.
func worker(curWorker int, job chan int, result chan string, wg *sync.WaitGroup) {
	defer wg.Done() // Mark worker as done once it exits

	for ch := range job {
		result <- fmt.Sprintf("%d -> %d", curWorker, ch) // Report task handling
	}
}

func main() {
	var wg sync.WaitGroup

	// Prepare 14 tasks (integers 0 to 13)
	task := make([]int, 0, 14)
	for i := 0; i < 14; i++ {
		task = append(task, i)
	}

	// Buffered result channel to collect output from all workers
	result := make(chan string, len(task))

	// Run the worker pool (this blocks until all jobs are enqueued)
	pool(3, task, result, &wg)

	// Close result channel once all workers finish
	go func() {
		wg.Wait()     // Wait for all worker goroutines to complete
		close(result) // Safe to close after all sends are done
	}()

	// Print results from result channel
	for r := range result {
		fmt.Println(r)
	}
}
