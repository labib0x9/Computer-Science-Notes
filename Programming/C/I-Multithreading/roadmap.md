---
title: POSIX Threads in C Roadmap
date: 2026-10-25
tags:
  - c
  - threads
  - pthreads
  - concurrency
  - roadmap
summary: A four-stage study roadmap for POSIX threads in C, from creating threads to a multi-threaded HTTP server
---

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