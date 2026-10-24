---
title: Concurrency Roadmap in Go
date: 2026-10-24
tags:
  - golang
  - concurrency
  - roadmap
summary: A study roadmap for Go concurrency, from OS theory to real-world projects
---

The roadmap below was generated with ChatGPT. I will tick off items as I learn them.

## 1. OS and theory foundations (before goroutines)

**Concepts**

- [ ] Processes vs threads
- [ ] User-level threads vs kernel-level threads
- [ ] M:N scheduler concepts
- [ ] Context switching: cost, registers, stacks
- [ ] Concurrency vs parallelism
- [ ] Stack vs heap, and dynamic stack growth
- [ ] Blocking I/O vs non-blocking I/O
- [ ] Race conditions and deadlocks
- [ ] Mutual exclusion and semaphores

**Goal:** understand what happens beneath goroutines.

## 2. Go runtime fundamentals

Learn Go's concurrency model and runtime scheduler.

**Concepts**

- [ ] The G, M, P model:
  - G is a goroutine.
  - M is a machine (an OS thread).
  - P is a processor (the logical context that runs a G).
- [ ] `GOMAXPROCS`
- [ ] How goroutines are scheduled
- [ ] Stack growth model
- [ ] Preemption (Go 1.14+)
- [ ] I/O poller (network poller)

**Tasks**

- [ ] Read about the Go scheduler (very important).
- [ ] Compare goroutine count and thread count using `runtime.NumGoroutine()` and `runtime.NumCPU()`.

## 3. Goroutine basics

Master the foundation of running concurrent tasks.

**Concepts**

- [ ] Creating goroutines
- [ ] Communication via channels
- [ ] Blocking behavior
- [ ] Buffered vs unbuffered channels
- [ ] How goroutines exit or get stuck

**Experiments**

- [ ] Launch 100k goroutines and observe memory usage.
- [ ] Build my first worker task using goroutines.

## 4. Synchronization primitives in Go

These map to OS-level primitives, but they are optimized and integrated with the runtime.

**Concepts**

- [ ] `sync.Mutex`
- [ ] `sync.RWMutex`
- [ ] `sync.WaitGroup`
- [ ] `sync.Cond`
- [ ] `sync.Once`
- [ ] `sync.Map`
- [ ] Semaphores using channels (a classic Go trick)

**Tasks**

- [ ] Demonstrate a race condition using `go run -race`.
- [ ] Fix the race using a mutex and channel synchronization.

## 5. Channel mastery

Channels are Go's primary concurrency construct.

**Concepts**

- [ ] Unbuffered channel synchronization
- [ ] Buffered channel throughput
- [ ] Directional channels (`chan<-`, `<-chan`)
- [ ] Closing channels and detecting closure
- [ ] Deadlocks and channel misuse
- [ ] The `select` statement (the most important skill)

**Tasks**

- [ ] Write a producer-consumer using buffered channels.
- [ ] Use `select` with a timeout (`time.After`).
- [ ] Build a job dispatcher that feeds workers.

## 6. Key concurrency patterns (Go-specific)

This is where Go becomes powerful.

**Patterns**

- [ ] Worker pool
- [ ] Fan-out / fan-in
- [ ] Producer → pipeline → consumer
- [ ] Tee pattern
- [ ] Pub-sub using channels
- [ ] Rate limiting using `time.Ticker`
- [ ] Reusable goroutine pools
- [ ] Atomic counters and lock-free techniques

**Tasks**

- [ ] Implement a general-purpose worker pool.
- [ ] Build a pipeline with 3–4 stages.
- [ ] Build a fan-in aggregator from N workers.

## 7. Context and cancellation (real-world essential)

I must learn how to stop goroutines cleanly.

**Concepts**

- [ ] `context.Context`
- [ ] `WithCancel`, `WithTimeout`, `WithDeadline`
- [ ] Propagating cancellation in pipelines
- [ ] Avoiding goroutine leaks

**Tasks**

- [ ] Implement a cancellable worker.
- [ ] Add cancellation to my pipeline.
- [ ] Build an HTTP worker with a timeout.

## 8. Multithreading and parallelism in Go

Understand how goroutines scale across CPU cores.

**Concepts**

- [ ] `GOMAXPROCS`
- [ ] Runtime scheduler and thread parking
- [ ] Parallel map-reduce patterns
- [ ] CPU-bound vs I/O-bound goroutines
- [ ] Atomics (`sync/atomic`)

**Tasks**

- [ ] Benchmark goroutines with 1 CPU vs N CPUs.
- [ ] Create a parallel sum using goroutines.
- [ ] Use atomic counters safely.

## 9. Debugging and profiling concurrency

A critical skill for production Go.

**Tools**

- [ ] `runtime.NumGoroutine()`
- [ ] `pprof`
- [ ] Goroutine dump (`debug.Stack()`)
- [ ] `go test -race`
- [ ] `go tool trace`
- [ ] Visualizing scheduler events

**Tasks**

- [ ] Create a leak on purpose and detect it with `pprof`.
- [ ] Profile a goroutine-heavy program.
- [ ] Use `trace` to see scheduler behavior.

## 10. Real-world concurrency projects

Building these will cement my understanding.

**Beginner**

- [ ] Concurrent web fetcher
- [ ] Parallel file hasher
- [ ] Worker pool that executes tasks

**Intermediate**

- [ ] Log processing pipeline
- [ ] Concurrent file searcher
- [ ] Chat server using goroutines

**Advanced**

- [ ] Job scheduler with timeouts
- [ ] In-memory database with transaction locks
- [ ] Network load balancer
- [ ] Actor system using goroutine mailboxes
- [ ] Mini Kubernetes-like task scheduler

## Resources

- [Exploring the depths of Golang channels (Medium)](https://medium.com/%40ravikumar19997/exploring-the-depths-of-golang-channels-a-comprehensive-guide-53e1a97cafe6)
- [7 powerful Golang concurrency patterns (Cristian Curteanu)](https://cristiancurteanu.com/7-powerful-golang-concurrency-patterns-that-will-transform-your-code-in-2025/)
- [Goroutines: Under the Hood, Vicki Niu (video)](https://www.youtube.com/watch?v=S-MaTH8WpOM&t=1194s)
- [Understanding Channels, Kavya Joshi (video)](https://www.youtube.com/watch?v=KBZlN0izeiY&t=191s)