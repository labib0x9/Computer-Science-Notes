### Goroutines And Concurrency - This chapter is incomplete. :)
---
ChatGpt Roadmap
---

# OS & Theory Foundations (Before Goroutines)
**Concepts**
- Processes vs Threads
- User-level threads vs kernel-level threads
- M:N scheduler concepts
- Context switching: cost, registers, stacks
- Concurrency vs Parallelism
- Stack vs Heap; dynamic stack growth
- Blocking I/O vs Non-blocking I/O
- Race conditions, deadlocks
- Mutual exclusion, semaphores
**Goal**
- Understand what happens beneath goroutines.

# Go Runtime Fundamentals
Learn Go’s concurrency model and runtime scheduler.
**Concepts**
- G, M, P model <br>
G → goroutine <br>
M → machine (OS thread) <br>
P → processor (logical context for running G) <br>
- GOMAXPROCS
- How goroutines are scheduled
- Stack growth model
- Preemption (Go 1.14+)
- I/O poller (network poller)
**Tasks**
- Read about the Go scheduler (very important).
- Observe goroutine count vs thread count using runtime.NumGoroutine() and runtime.NumCPU().

# Goroutine Basics
Master the foundation of running concurrent tasks.
**Concepts**
- Creating goroutines
- Communication via channels
- Blocking behavior
- Buffered vs unbuffered channels
- How goroutines exit or get stuck
**Experiments**
- Launch 100k goroutines → observe memory usage.
- Build your first worker task using goroutines.

# Synchronization Primitives in Go
These map to OS-level primitives but are optimized and integrated with the runtime.
**Concepts**
- sync.Mutex
- sync.RWMutex
- sync.WaitGroup
- sync.Cond
- sync.Once
- sync.Map
- Semaphores using channels (classic Go trick)
**Tasks**
- Race condition demo using go run -race.
- Fix races using mutex and channel synchronization.

# Channel Mastery
Channels are Go’s primary concurrency construct.
**Concepts**
- Unbuffered channel synchronization
- Buffered channel throughput
- Directional channels (chan<-, <-chan)
- Closing channels, detecting closure
- Deadlocks and channel misuse
- Select statement (the most important skill)
**Tasks**
- Write producer-consumer using buffered channels.
- Use select with timeout (time.After).
- Build a job dispatcher feeding workers.

# Key Concurrency Patterns (Go-Specific)
This is where Go becomes powerful.
**Learn patterns**
- Worker pool
- Fan-out / Fan-in
- Producer → pipeline → consumer
- Tee pattern
- Pub-sub using channels
- Rate limiting using time.Ticker
- Reusable goroutine pools
- Atomic counters + lock-free techniques
**Tasks**
- Implement a general-purpose worker pool.
- Build a pipeline with 3–4 stages.
- Build a fan-in aggregator from N workers.

# Context and Cancellation (Real-World Essential)
You must learn how to stop goroutines cleanly.
**Concepts**
- context.Context
- WithCancel, WithTimeout, WithDeadline
- Propagating cancellation in pipelines
- Avoiding goroutine leaks
**Tasks**
- Implement a cancellable worker.
- Add cancellation to your pipeline.
- Build a timeout-enabled HTTP worker.

# Multithreading & Parallelism in Go
Understanding how goroutines scale across CPU cores.
**Concepts**
- GOMAXPROCS
- Runtime scheduler and thread parking
- Parallel map-reduce patterns
- CPU-bound vs I/O-bound goroutines
- Using atomics (sync/atomic)
**Tasks**
- Benchmark goroutines with 1 vs N CPUs.
- Create a parallel sum using goroutines.
- Use atomic counters safely.

# Debugging & Profiling Concurrency
Critical skill for production Go.
**Tools**
- runtime.NumGoroutine()
- pprof
- Goroutine dump (debug.Stack())
- go test -race
- go tool trace
- Visualizing scheduler events
**Tasks**
- Create a leak on purpose → detect it with pprof.
- Profile a goroutine-heavy program.
- Use trace to see scheduler behavior.

# Build Real-World Concurrency Projects
This cements understanding.
**Beginner**
- Concurrent web fetcher
- Parallel file hasher
- Worker pool executing tasks
**Intermediate**
- Log processing pipeline
- Concurrent file searcher
- Chat server using goroutines
**Advanced**
- Job scheduler with timeouts
- In-memory database with transaction locks
- Networking load balancer
- Actor system using goroutine-mailboxes
- Mini-kubernetes-like scheduler for tasks

## Resources 
- https://medium.com/%40ravikumar19997/exploring-the-depths-of-golang-channels-a-comprehensive-guide-53e1a97cafe6
- https://cristiancurteanu.com/7-powerful-golang-concurrency-patterns-that-will-transform-your-code-in-2025/?utm_source=chatgpt.com
- [Goroutines: Under the Hood | Vicki Niu](https://www.youtube.com/watch?v=S-MaTH8WpOM&t=1194s)
- [Understanding Channels - Kavya Joshi](https://www.youtube.com/watch?v=KBZlN0izeiY&t=191s)