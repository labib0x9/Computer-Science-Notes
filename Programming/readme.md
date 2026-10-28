---
title: Programming Languages
date: 2026-10-24
---

Now that we know the concepts, how do we implement them? With programming, of course.

---

## Overview

I will cover only three programming languages (C, C++ and Golang), because those are the only ones I know. I will update this if I ever switch to another.

- **Go**: I may use a debugger and a profiler to look at what is happening at runtime.
- **C and C++**: I cover only syntax and concepts, maybe some OOP for C++.


## Roadmaps
<details>
<summary><strong>POSIX Threads in C Roadmap</strong></summary>
The roadmap below was generated with ChatGPT. I will tick off items as I learn them.

## Stage 1: Foundations (threads and memory model)

**Goal:** understand how threads work: stacks, sharing, creation and lifecycle.

**Concepts**

- [ ] What is a thread?
- [ ] Process vs thread
- [ ] POSIX threads overview
- [ ] `pthread_create`, `pthread_join`, `pthread_exit`
- [ ] Thread function signatures
- [ ] Passing arguments to threads
- [ ] Returning values from threads
- [ ] Thread-local storage (`__thread` or `pthread_key_t`)

**Resources**

- TLPI (*The Linux Programming Interface*), the chapters on threads
- POSIX Threads Programming tutorial (Lawrence Livermore National Laboratory)
- GeeksforGeeks: POSIX Threads basics (quick reference)

**Mini projects (do all)**

- [ ] Spawn 10 threads that print their IDs.
- [ ] Vector addition using threads (split the array into N chunks).
- [ ] Compute the sum of N numbers using multiple threads.
- [ ] Thread-local counter vs shared counter (observe the race).

**Outcome:** I understand thread creation, joining, stacks and passing arguments.

## Stage 2: Synchronization primitives

**Goal:** learn every core primitive that protects shared state.

**Concepts**

- [ ] Mutex (`pthread_mutex_t`)
  - [ ] Lock and unlock
  - [ ] Deadlocks
  - [ ] Re-entrant (recursive) vs non-reentrant mutex
  - [ ] Contention impact
- [ ] Condition variables (`pthread_cond_t`)
  - [ ] Waiting and signaling
  - [ ] Blocking and waking up threads
  - [ ] Producer-consumer pattern
- [ ] Semaphores
  - [ ] `sem_wait`, `sem_post`, `sem_init`
  - [ ] Binary semaphore
  - [ ] Counting semaphore
- [ ] Read-write locks (`pthread_rwlock_t`)
  - [ ] When to use an RW lock instead of a mutex
- [ ] Barriers (`pthread_barrier_t`)
  - [ ] Thread phase synchronization

**Resources**

- TLPI, the synchronization chapters
- MIT 6.824 notes on concurrency (conceptual)
- POSIX man pages (for example `man pthread_cond_wait`)

**Projects**

- [ ] Producer-consumer queue (mutex + condition variable)
- [ ] Reader-writer lock example: many readers, few writers on shared data
- [ ] Semaphore-based job counter
- [ ] Barrier: wait for 10 threads to finish their partial work

**Outcome:** I understand how real systems manage shared state.

## Stage 3: Real concurrency patterns

**Goal:** know how real systems use threads.

**Patterns to study and implement**

- [ ] Thread pool (my own project)
  - [ ] Task queue
  - [ ] N worker threads
  - [ ] Condition variable for wake-ups
- [ ] Fan-out / fan-in
  - [ ] Threads compute tasks
  - [ ] The main thread aggregates the results
- [ ] Pipeline pattern
  - [ ] Stage 1 → Stage 2 → Stage 3
  - [ ] Each stage is a thread plus a queue
- [ ] Asynchronous tasks using `pthread` and a future
  - [ ] A thread that returns a result later
  - [ ] `pthread_join` retrieves the value
- [ ] Work stealing (conceptual only; no need to implement)
  - [ ] Basic introduction to schedulers

**Resources**

- Intel Threading Building Blocks docs (conceptual)
- Linux kernel threads documentation (conceptual exposure)
- My own thread-pool implementation

**Projects**

- [ ] Complete my thread pool.
- [ ] Concurrent logger (a worker handles the log messages).
- [ ] Mini pipeline:
  - [ ] Thread 1: read a file
  - [ ] Thread 2: count words
  - [ ] Thread 3: output the result

**Outcome:** I think in systems patterns, not only low-level APIs.

## Stage 4: Advanced and real-world systems

**Goal:** understand the complex topics used in databases, runtimes and kernels.

**Topics**

- [ ] Lock-free programming and atomics
  - [ ] C11 atomic types
  - [ ] Compare-and-swap (CAS)
  - [ ] Memory ordering basics
  - [ ] Implement a simple lock-free stack
- [ ] Thread scheduling
  - [ ] Linux scheduler basics
  - [ ] CPU affinity (`pthread_setaffinity_np`)
  - [ ] Thread priorities
- [ ] Thread models in large systems
  - [ ] Nginx worker model
  - [ ] Redis I/O thread model
  - [ ] Database thread models (MySQL/InnoDB)
- [ ] Deadlock and race detection tools
  - [ ] Valgrind Helgrind
  - [ ] ThreadSanitizer (clang)

**Resources**

- C11 atomic operations tutorial
- Linux kernel documentation on RCU (conceptual)
- The concurrency chapter of *Computer Systems: A Programmer's Perspective* (CS:APP)

**Projects**

- [ ] Lock-free stack or queue using atomics
- [ ] Parallel file processing with a thread pool
- [ ] Thread benchmark:
  - [ ] Increase the thread count
  - [ ] Measure performance
  - [ ] Observe scheduler behavior

**Outcome:** I am ready to build multi-threaded servers, databases and tools in C.

## Capstone project

Build a mini multi-threaded HTTP server:

- [ ] The main thread accepts connections.
- [ ] A thread pool handles the requests.
- [ ] A shared log, protected by a mutex.
- [ ] Each worker parses the request, serves a static file, and logs the request.

**This project uses:**

- Threads
- Mutex
- Condition variable
- Queue
- Pipeline
- I/O multiplexing (optional)
- Thread pool
- Real-world concurrency design

This resembles the threaded model of older web servers such as Apache's worker mode.
</details>


<details>
<summary><strong>Pointers in C Roadmap</strong></summary>
The roadmap below was generated with ChatGPT. I will tick off items as I learn them.

## Stage 1: Memory model and basics

**Concepts**

- [ ] How variables are stored (stack vs static vs heap)
- [ ] What an address is, and how pointers store addresses
- [ ] `&` (address-of) and `*` (dereference)
- [ ] Pointer types: `int*`, `char*`, `void*`
- [ ] Null pointer, wild pointer, dangling pointer

**Practice**

- [ ] Print the addresses of variables.
- [ ] Use `int*` and `char*` to modify values indirectly.
- [ ] Write small examples, such as swap using pointers.
- [ ] Build a "memory inspector" program that prints addresses and values for a few variables.

## Stage 2: Pointer arithmetic

**Concepts**

- [ ] Pointer increment and decrement
- [ ] Scaling by type size (`ptr + 1` moves by `sizeof(type)` bytes)
- [ ] Traversing arrays using pointers
- [ ] Arrays vs pointers (`int arr[]` vs `int *ptr`)

**Practice**

- [ ] Implement `strlen`, `strcpy` and `memcpy` using only pointer arithmetic.
- [ ] Iterate through integer arrays without indexing.
- [ ] Write an array iterator that takes a base pointer, an element size and a count, and prints the elements.

## Stage 3: Pointers and functions

**Concepts**

- [ ] Passing pointers to functions
- [ ] `const` correctness (`const int *p` vs `int *const p`)
- [ ] Double pointers (`int **`)
- [ ] Dynamic 2D arrays using pointers-to-pointers

**Practice**

- [ ] Implement dynamic matrix creation using `int**`.
- [ ] Write functions that modify multiple variables using double pointers.
- [ ] Implement a dynamic table: create rows, add values, and free it properly.

## Stage 4: Pointers and the heap

**Concepts**

- [ ] `malloc`, `calloc`, `realloc`, `free`
- [ ] Ownership and lifetime
- [ ] Memory leaks vs dangling pointers
- [ ] `sizeof` and alignment rules

**Practice**

- [ ] Allocate arrays dynamically.
- [ ] Write custom resize functions for arrays.
- [ ] Detect memory leaks with Valgrind or AddressSanitizer.
- [ ] Implement my own dynamic array type, like C++ `vector<int>`.

## Stage 5: Strings and char pointers

**Concepts**

- [ ] A string as `char*`
- [ ] String literals (read-only) vs mutable buffers
- [ ] Buffer overflows and bounds checking

**Practice**

- [ ] Implement my own `strtok`, `strdup` and `atoi`.
- [ ] Validate input to avoid overflow.
- [ ] Build a simple tokenizer using only pointers.

## Stage 6: Complex pointer declarations

**Concepts**

- [ ] Function pointers, for example `int (*f)(int, int)`
- [ ] Arrays of pointers
- [ ] Pointers to arrays
- [ ] `typedef` for readability
- [ ] Callback functions

**Practice**

- [ ] Implement callback-based sorting (like `qsort`).
- [ ] Use an array of function pointers for a command dispatcher.
- [ ] Create a "plugin" system: store function pointers in a table and call them dynamically.

## Stage 7: Low-level pointer applications

**Concepts**

- [ ] Struct pointers and `->`
- [ ] Offset calculations and `offsetof`
- [ ] Casting pointers (`void*`, `char*`)
- [ ] Memory layout, padding and alignment

**Practice**

- [ ] Implement a packed struct and examine its binary layout.
- [ ] Parse a custom binary file format using pointer casts.
- [ ] Build a simple binary serializer and deserializer for a struct.

## Stage 8: Advanced patterns

**Concepts**

- [ ] Smart-pointer patterns in C (manual reference counting)
- [ ] Pointer tagging (using the low bits)
- [ ] Memory pools and arenas
- [ ] Buffers with pointer cursors

**Practice**

- [ ] Build an arena allocator using one large `malloc`.
- [ ] Perform pointer bump allocation.
- [ ] Implement an allocator that returns aligned blocks.

## Stage 9: System programming with pointers

**Concepts**

- [ ] Pointers in system calls (buffers for `read` and `write`)
- [ ] `struct sockaddr*` in networking
- [ ] Pointers in file I/O and `mmap`
- [ ] Accessing hardware-like memory through pointers (only in OS-level code)

**Practice**

- [ ] Implement a small TCP client using pointer-based buffers.
- [ ] Implement my own `getline` using `realloc` and pointers.
- [ ] Use `mmap` to create a memory-backed file and navigate it with pointers.
- [ ] Build a tiny key-value store where both the data and the index are pointer-based structures.

## Stage 10: Capstone project

Implement a minimal storage engine in C:

- [ ] `malloc`- or `mmap`-based pages
- [ ] Pointer arithmetic to navigate pages
- [ ] Pointer-to-struct parsing
- [ ] Serialization and deserialization
- [ ] Dynamic resizing
- [ ] Free lists or pointer-based indexes

This consolidates everything: alignment, pointer math, dynamic memory, casts, function pointers and buffer management.

## How to study (1 hour daily)

- 20 minutes of theory (K&R, TLPI chapters 5–7, Molnar's *System Programming with C*)
- 20 minutes of practice tasks
- 20 minutes of incremental project building
</details>


<details>
<summary><strong>Concurrency Roadmap in Go</strong></summary>
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

</details>

## Resources

- [Stack Frame and GDB](https://www.youtube.com/watch?v=7dMTCdFM2ss)
- [Garbage Collector](https://journal.stuffwithstuff.com/2013/12/08/babys-first-garbage-collector/)
- [Memory Management in C](https://www.youtube.com/watch?v=rJrd2QMVbGM&t=30s)
- [Socket to Me: Where Do Sockets Live in Go? - Gabbi Fisher](https://www.youtube.com/watch?v=pGR3r0UhoS8)
- [Asynchronous Networking Patterns - Filippo Valsorda](https://www.youtube.com/watch?v=afSiVelXDTQ)