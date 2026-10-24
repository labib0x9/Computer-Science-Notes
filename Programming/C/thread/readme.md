### Threads - POSIX Threads
---
ChatGpt Roadmap to learn thread in C
---
# Foundations: Threads & Memory Model
**Goal:** Understand how threads work, stack, sharing, creation, and lifecycle.
**Concepts**
- What is a thread?
- Process vs thread
- POSIX threads overview
- pthread_create, pthread_join, pthread_exit
- Thread function signatures
- Passing arguments to threads
- Returning values from threads
- Thread-local storage (__thread or pthread_key_t)
**Resources**
- TLPI (The Linux Programming Interface) – Chapters on Threads
- POSIX Threads Tutorial — Lawrence Livermore Lab (free PDF)
- GeeksForGeeks: POSIX Threads Basics (quick reference)
**Mini Projects (Do all):**
- Spawn 10 threads printing their ID.
- Vector addition using threads (split array into N chunks).
- Compute sum of N numbers using multi-threading.
- Thread-local counter vs shared counter (observe race).
**Outcome:** You fully understand thread creation, joining, stack, arguments.

# Synchronization Primitives
**Goal:** Learn every core primitive that protects shared state.
**Concepts**
- Mutex (pthread_mutex_t)
- Lock/unlock
- Deadlocks
- Re-entrant vs non-reentrant mutex
- Contention impact
- Condition Variables (pthread_cond_t)
- Waiting/signaling
- Blocking & waking up threads
- Producer–consumer pattern
- Semaphores
- sem_wait, sem_post, sem_init
- Binary semaphore
- Counting semaphore
- Read-Write Locks
- pthread_rwlock_t
- When to use RW-lock instead of mutex
- Barriers
- pthread_barrier_t
- Thread phase synchronization
**Resources**
- TLPI — Synchronization Chapters
- MIT 6.824 notes on concurrency (conceptual)
- POSIX man pages (man pthread_cond_wait, etc.)
**Projects**
- Producer–Consumer queue (Mutex + Condvar)
- Reader–Writer lock example:
- Many readers, few writers on shared data.
- Semaphore-based job counter
- Barrier: wait for 10 threads to compute partial work
**Outcome:** You will fully understand how real systems manage shared state.

# Real Concurrency Patterns
**Goal:** Know how modern systems architecture uses threads in reality.
- Patterns to study and implement
- Thread Pool (your current project)
- Task queue
- N worker threads
- Condition variable for wakeups
- Fan-Out / Fan-In
- Threads compute tasks
- Main thread aggregates results
- Pipeline Pattern
- Stage1 → Stage2 → Stage3
- Each stage = a thread + queue
- Asynchronous Tasks using pthread + future
- Thread that returns a result later
- pthread_join retrieves value
- Work Stealing (Conceptual)
- Basic intro to schedulers
- Only conceptual, no need to implement
**Resources**
- Intel Threading Building Blocks Docs (conceptual)
- Linux kernel threads documentation (conceptual exposure)
- Your own thread-pool implementation
**Projects**
- Complete your thread pool
- Concurrent Logger (worker handles log messages)
- Mini-pipeline:
    - Thread1: read file
    - Thread2: count words
    - Thread3: output result
**Outcome:** You now think in systems patterns, not low-level APIs.

# Advanced & Real-World Systems
**Goal:** Understand complex topics used in databases, runtimes, kernels.
**Topics**
- Lock-Free & Atomics
- C11 atomic types
- Compare-and-Swap (CAS)
- Memory ordering basics
- Implement lock-free stack (simple)
- Thread Scheduling
- Linux scheduler basics
- CPU affinity (pthread_setaffinity_np)
- Thread priorities
- Thread Pools in Large Systems
- Nginx worker model
- Redis I/O thread model
- Database thread models (MySQL/InnoDB)
- Deadlocks & Race Detection Tools
- Valgrind Helgrind
- ThreadSanitizer (clang)
**Resources**
- C11 Atomic operations tutorial
- Linux Kernel Documentation on RCU (conceptual)
- Concurrency chapter from Computer Systems: A Programmer’s Perspective (CS:APP)
**Projects**
- Lock-free stack or queue using atomics
- Parallel file processing with thread pool
- Thread benchmark:
- Increase thread count
- Measure performance
- Observe scheduler behavior
**Outcome:** You are ready to build multi-threaded servers, databases, and tools in C.

# Final Capstone Projects
- Build a mini multi-threaded HTTP server:
    - Main thread accepts connections
    - A thread pool handles requests
    - Shared log with mutex
    - Each worker:
        - Parse request
        - Serve static file
        - Log request
**This project activates:**
- threads
- mutex
- condvar
- queue
- pipeline
- I/O multiplexing (optional)
- thread pool
- real-world concurrency design <br>
This is exactly how early versions of Nginx/Apache worked.