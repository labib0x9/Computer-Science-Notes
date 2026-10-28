---
title: "Multithreading in C"
date: 2026-10-25
tags:
  - c
  - pthreads
  - concurrency
  - systems
summary: POSIX threads in C, from create and join to mutexes, condition variables, and a small thread pool
---

# Overview

POSIX threads (**pthreads**) in C: create threads, wait for them, share data safely, and build a small thread pool. Every example compiles cleanly with:

```bash
gcc -Wall -Wextra -pthread main.c -o main
```

`-pthread` is needed at both compile and link time. On newer glibc (2.34+) the thread functions live in the main C library, so linking works without it, but the flag is still the portable way and also sets the right compile-time definitions.

Most snippets below are fragments. Add the includes they need (`<stdio.h>`, `<stdlib.h>`, `<pthread.h>`, and `<unistd.h>` for `sleep`).

## Threads vs processes

A **thread** is an independent flow of execution inside a process. All threads of a process share one address space, so they see the same global variables and heap, but each thread has its own stack.

```text
Shared by all threads:   text | rodata | data | bss | heap
Private to each thread:  its own stack (and registers, and thread-local variables)
```

This ties back to the memory layout note, and it gives the main rule of threading:
- **Local variables are private**, so they are safe, unless you hand their address to another thread.
- **Globals, `static` variables and heap memory are shared**, so any access that can overlap with a write needs synchronization.

A pthread is a 1:1 kernel thread, scheduled by the OS. Go's goroutines are M:N: many goroutines are multiplexed over a few OS threads by the runtime's GMP scheduler. The "spawn and run concurrently" feel is similar, but the cost is very different. A pthread reserves about 8 MB of virtual stack by default on Linux, and a goroutine starts at roughly 2 KB.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [Creating threads](#1-creating-threads) | `pthread_create`, the `void *(void *)` signature |
| 2 | [Waiting for threads](#2-waiting-for-threads) | `pthread_join`, why `sleep()` in `main` is wrong |
| 3 | [Thread IDs](#3-thread-ids) | `pthread_self` |
| 4 | [Passing arguments](#4-passing-arguments) | The `&i` loop pitfall |
| 5 | [Returning values](#5-returning-values) | Heap result and `pthread_join` |
| 6 | [Race conditions](#6-race-conditions) | Lost updates on `count++` |
| 7 | [Mutex](#7-mutex) | `pthread_mutex_t`, cost, batching |
| 8 | [Atomics](#8-atomics) | `stdatomic.h`, stop flags |
| 9 | [Condition variables](#9-condition-variables-bounded-queue-and-thread-pool) | Bounded queue and thread pool |
| 10 | [Deadlock](#10-deadlock) | Lock ordering |
| 11 | [Detached threads and attributes](#11-detached-threads-and-attributes) | `pthread_detach`, stack size |
| 12 | [Error handling](#12-error-handling) | pthread functions don't set `errno` |
| 13 | [Thread-local data and thread safety](#13-thread-local-data-and-thread-safety) | `_Thread_local`, `strtok_r` |
| 14 | [Other primitives](#14-other-primitives) | rwlock, semaphore, `pthread_once` |
| 15 | [Debugging](#15-debugging) | ThreadSanitizer, gdb |
| 16 | [Coming from Go](#16-coming-from-go) | The equivalent tools |
| 17 | [Platform notes](#17-platform-notes) | Linux vs macOS |
| 18 | [Exercises](#18-exercises) | Practice |

---

# 1. Creating threads

```c
int pthread_create(pthread_t *thread, const pthread_attr_t *attr,
                   void *(*start)(void *), void *arg);
```

The thread function must have the shape `void *f(void *)`. The `void *` parameter and result are how you pass data in and out (see sections 4 and 5).

```c
#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

void *AliceTurn(void *args) {
    (void)args;                      // unused: silences -Wextra
    for (int i = 0; i <= 10; i++) {
        sleep(1);
        printf("Alice = %d\n", i);
    }
    return NULL;
}

void *BobTurn(void *args) {
    (void)args;
    for (int i = 0; i <= 4; i++) {
        sleep(1);
        printf("Bob = %d\n", i);
    }
    return NULL;
}

int main(void) {
    pthread_t alice, bob;
    pthread_create(&alice, NULL, AliceTurn, NULL);
    pthread_create(&bob, NULL, BobTurn, NULL);

    sleep(10);   // BAD, see the next section
    return 0;
}
```

When `main` returns, the process exits and **every thread is killed**, finished or not. Go behaves the same way: when `main` exits, all goroutines die.

---

# 2. Waiting for threads

The `sleep(10)` above is a race against the workers. Alice needs 11 seconds (prints `0` to `10`), so she is killed before her last line. Guessing a duration is never correct. Use `pthread_join` instead.

```c
int main(void) {
    pthread_t alice, bob;
    pthread_create(&alice, NULL, AliceTurn, NULL);
    pthread_create(&bob, NULL, BobTurn, NULL);

    pthread_join(bob, NULL);     // blocks until Bob returns
    pthread_join(alice, NULL);   // blocks until Alice returns
    return 0;
}
```

- Join order doesn't matter for correctness. Joining a thread that already finished returns immediately.
- Every joinable thread must be joined (or detached), otherwise its resources stay allocated until the process exits.
- Alternative: ending `main` with `pthread_exit(NULL)` instead of `return` lets the other threads keep running. The process exits when the last thread finishes.

---

# 3. Thread IDs

```c
void *Calc(void *args) {
    (void)args;
    printf("Calc() = %lu\n", (unsigned long)pthread_self());
    return NULL;
}

int main(void) {
    printf("main() = %lu\n", (unsigned long)pthread_self());

    pthread_t t[10];
    for (int i = 0; i < 10; i++) pthread_create(&t[i], NULL, Calc, NULL);
    for (int i = 0; i < 10; i++) pthread_join(t[i], NULL);

    for (int i = 0; i < 10; i++)
        printf("%d ThreadId = %lu\n", i, (unsigned long)t[i]);
    return 0;
}
```

- `pthread_t` is an **opaque type**: its real type is up to the platform. The `unsigned long` cast works on Linux and macOS for printing, but isn't guaranteed portable. Compare IDs with `pthread_equal(a, b)`, never `==`.
- This is the pthread ID, not the kernel thread ID that `top -H` or gdb shows. For that, use `syscall(SYS_gettid)` (or `gettid()` on glibc 2.30+). macOS has `pthread_threadid_np`.
- After a thread is joined, the OS may reuse its ID for a new thread, so don't treat IDs as permanently unique.
- The order of the `Calc()` lines is nondeterministic. Never rely on it.

---

# 4. Passing arguments

The `void *arg` is a pointer, so whatever it points to must stay valid until the thread has read it.

```c
void *Worker(void *arg) {
    int id = *(int *)arg;
    printf("worker %d\n", id);
    return NULL;
}

int main(void) {
    pthread_t t[4];
    int ids[4];

    for (int i = 0; i < 4; i++) {
        ids[i] = i;                                    // one slot per thread
        pthread_create(&t[i], NULL, Worker, &ids[i]);
    }
    for (int i = 0; i < 4; i++) pthread_join(t[i], NULL);
    return 0;
}
```

**Pitfall:** `pthread_create(&t[i], NULL, Worker, &i)` is a bug. Every thread points at the same `i`, which keeps changing and goes out of scope after the loop. You will see duplicate or garbage IDs.

Options for passing data:

| Method | Good for | Watch out |
| --- | --- | --- |
| Pointer to an array element, as above | a few threads, `main` joins them | array must outlive the threads |
| `(void *)(intptr_t)i` and `(int)(intptr_t)arg` | one integer (`<stdint.h>`) | a hack; fine for small ints |
| `malloc`'d struct, thread frees it | several values, detached threads | someone must free it, exactly once |

For more than one value, put them in a struct and pass its address.

---

# 5. Returning values

A thread returns a `void *`, and `pthread_join` hands it back through a `void **`. Returning the address of a local variable is a use-after-return (the stack frame is gone), so allocate on the heap, or write into a struct the caller owns.

```c
void *Sum(void *arg) {
    int n = *(int *)arg;
    int *res = malloc(sizeof *res);
    if (res == NULL) return NULL;
    *res = 0;                           // malloc doesn't zero memory
    for (int i = 1; i <= n; i++) *res += i;
    return res;
}

int main(void) {
    int n = 100;                        // safe to share: main joins before n goes away
    pthread_t t;
    pthread_create(&t, NULL, Sum, &n);

    void *ret;
    pthread_join(t, &ret);
    if (ret != NULL) {
        int *result = ret;
        printf("%d\n", *result);        // 5050
        free(result);                   // the joiner owns the allocation
    }
    return 0;
}
```

Use a plain `void *ret` instead of casting `&result` to `void **`. Treating an `int **` as a `void **` is technically undefined behavior.

---

# 6. Race conditions

```c
int count;

void *Calc(void *args) {
    (void)args;
    for (int i = 1; i <= 10000000; i++) count++;
    return NULL;
}

int main(void) {
    pthread_t t;
    pthread_create(&t, NULL, Calc, NULL);
    Calc(NULL);                         // main works too
    pthread_join(t, NULL);
    printf("%d\n", count);              // expected 20000000, usually smaller
    return 0;
}
```

`count++` is really three steps: load, add, store. Two threads can both load `5`, both store `6`, and one increment is lost. The result changes from run to run, so run it a few times.

This is a **data race**: two threads access the same memory, at least one writes, and nothing orders them. In C that is undefined behavior, so the result isn't just "slightly off". With `-O2` the compiler may keep `count` in a register for the whole loop, or even fold the loop into one addition, so the program can print the right answer on some builds and still be broken.

---

# 7. Mutex

A **mutex** (mutual exclusion lock) lets only one thread at a time run the code between `lock` and `unlock`.

```c
pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
int count;

void *Calc(void *args) {
    (void)args;
    for (int i = 1; i <= 10000000; i++) {
        pthread_mutex_lock(&lock);
        count++;
        pthread_mutex_unlock(&lock);
    }
    return NULL;
}
```

Now it always prints `20000000`. Rules:

- Every access to the shared variable must take the **same** lock. One unguarded access breaks everything.
- Keep critical sections tiny, and never block (I/O, `sleep`) while holding a lock.
- Unlock on every exit path: early `return`, `break`, and error branches.
- Only the thread that locked a mutex may unlock it.
- Dynamic setup: `pthread_mutex_init(&m, NULL)` paired with `pthread_mutex_destroy(&m)`.
- `pthread_mutex_trylock` returns `EBUSY` instead of waiting, so you can do something else.

**Cost:** locking 20 million times is slow, because the threads fight over the lock. Reduce lock traffic by accumulating privately and merging once:

```c
void *Calc(void *args) {
    (void)args;
    int local = 0;
    for (int i = 1; i <= 10000000; i++) local++;   // private, no lock

    pthread_mutex_lock(&lock);
    count += local;                                // one lock per thread
    pthread_mutex_unlock(&lock);
    return NULL;
}
```

Time the three versions (racy, per-iteration lock, batched) with `time ./main` and compare.

---

# 8. Atomics

For a plain counter, a lock is overkill. C11 **atomic** operations do the whole read-modify-write as one indivisible step (one `lock xadd` instruction on x86):

```c
#include <stdatomic.h>

atomic_int count;                      // static storage: starts at 0

void *Calc(void *args) {
    (void)args;
    for (int i = 1; i <= 10000000; i++)
        atomic_fetch_add(&count, 1);
    return NULL;
}
```

For a pure counter you can use `atomic_fetch_add_explicit(&count, 1, memory_order_relaxed)`. That is faster because it asks only for atomicity, not for ordering with other memory operations. The default (`seq_cst`) is the strictest and the safest choice when unsure.

## Stop flags

A plain `int` used to tell a thread to stop is also a data race, and the compiler is allowed to read it once and loop forever:

```c
atomic_bool stop;                      // starts false

void *Worker(void *arg) {
    (void)arg;
    while (!atomic_load(&stop)) {
        /* do work */
    }
    return NULL;
}

/* elsewhere */
atomic_store(&stop, true);
```

## Limits

An atomic protects **one variable**. Rules that span several variables (for example a queue's `size`, `readAt` and `insertAt` that must change together) still need a mutex. Atomics also can't make a thread *wait*. Use a condition variable for that.

---

# 9. Condition variables: bounded queue and thread pool

A mutex answers "who may touch the data". A **condition variable** answers "sleep until the data is in the state I need". The producer pushes tasks into a fixed-size ring buffer. Workers pop and run them. A full queue blocks the producer, and an empty queue blocks the workers.

```c
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define QUEUE_CAP   10
#define NUM_WORKERS 4

typedef struct Task {
    int a, b;
} Task;

static Task queue[QUEUE_CAP];
static int insertAt = 0, readAt = 0, size = 0, finish = 0;

static pthread_mutex_t lock     = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t  notEmpty = PTHREAD_COND_INITIALIZER;  // workers wait on this
static pthread_cond_t  notFull  = PTHREAD_COND_INITIALIZER;  // producer waits on this

void push(Task task) {
    pthread_mutex_lock(&lock);
    while (size == QUEUE_CAP)                    // full: wait for a worker to pop
        pthread_cond_wait(&notFull, &lock);

    queue[insertAt] = task;
    insertAt = (insertAt + 1) % QUEUE_CAP;
    size++;

    pthread_cond_signal(&notEmpty);              // wake one worker
    pthread_mutex_unlock(&lock);
}

// returns 1 and fills *t, or 0 when the queue is drained and closed
int pop(Task *t) {
    pthread_mutex_lock(&lock);
    while (size == 0 && !finish)                 // empty: wait for push or shutdown
        pthread_cond_wait(&notEmpty, &lock);

    if (size == 0 && finish) {                   // nothing left, pool is closing
        pthread_mutex_unlock(&lock);
        return 0;
    }

    *t = queue[readAt];
    readAt = (readAt + 1) % QUEUE_CAP;
    size--;

    pthread_cond_signal(&notFull);               // wake the producer
    pthread_mutex_unlock(&lock);
    return 1;
}

void executeTask(Task task) {
    printf("ThreadId = %lu  ->  %d + %d = %d\n",
           (unsigned long)pthread_self(), task.a, task.b, task.a + task.b);
}

void *startThread(void *args) {
    (void)args;
    Task t;
    while (pop(&t)) executeTask(t);
    return NULL;
}

int main(void) {
    pthread_t workers[NUM_WORKERS];

    for (int i = 0; i < NUM_WORKERS; i++)
        pthread_create(&workers[i], NULL, startThread, NULL);

    for (int i = 1; i <= 20; i++)
        push((Task){ .a = i, .b = i * i });

    // shutdown: set the flag under the lock, wake EVERY sleeping worker
    pthread_mutex_lock(&lock);
    finish = 1;
    pthread_cond_broadcast(&notEmpty);
    pthread_mutex_unlock(&lock);

    for (int i = 0; i < NUM_WORKERS; i++)
        pthread_join(workers[i], NULL);
    return 0;
}
```

Why it works:

- **`while`, not `if`, around `pthread_cond_wait`.** Wakeups can be **spurious** (happen with no signal), and another thread may take the item between the signal and your wakeup. Always re-check the condition.
- `pthread_cond_wait` atomically releases the mutex and goes to sleep, then re-acquires the mutex before returning. That is why you must hold the lock when you call it.
- **Shutdown uses `broadcast`, not `signal`.** Every idle worker must wake up and see `finish`. A single `signal` wakes one worker, and the rest would sleep forever, so `pthread_join` would hang.
- Workers drain the queue before exiting, because `finish` only stops them when `size == 0`.
- `pthread_cond_signal` can be called before or after unlocking. Both are correct. Signaling after unlock can save the woken thread from immediately blocking on the mutex.
- For a wait with a deadline, use `pthread_cond_timedwait`, which takes an absolute time (by default measured on `CLOCK_REALTIME`).
- Use `#define` (or `enum`) for array sizes in C. `const int N = 10; int a[N];` at file scope does **not** compile in C, because a `const` variable is not a constant expression (this differs from C++).

Naming tip: name condition variables after the state you are waiting *for* (`notEmpty`, `notFull`), not the state you're leaving.

## Making it a general pool

The pool above only adds two numbers. To run any job, change the task to hold a function pointer (from your pointer notes):

```c
typedef struct Task {
    void (*fn)(void *);
    void *arg;
} Task;

// the worker then does:  t.fn(t.arg);
```

Now the same pool can run any work, and the caller decides who owns and frees `arg`.

---

# 10. Deadlock

A **deadlock** is when threads wait for each other forever.

```c
static pthread_mutex_t A = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t B = PTHREAD_MUTEX_INITIALIZER;

void *t1(void *arg) {
    (void)arg;
    pthread_mutex_lock(&A);
    sleep(1);                      // give t2 time to take B
    pthread_mutex_lock(&B);        // waits for t2, forever
    pthread_mutex_unlock(&B);
    pthread_mutex_unlock(&A);
    return NULL;
}

void *t2(void *arg) {
    (void)arg;
    pthread_mutex_lock(&B);
    sleep(1);
    pthread_mutex_lock(&A);        // waits for t1, forever
    pthread_mutex_unlock(&A);
    pthread_mutex_unlock(&B);
    return NULL;
}
```

Each thread holds one lock and waits for the other. A deadlock needs four things at once: mutual exclusion, hold-and-wait, no preemption (locks can't be taken away), and a circular wait. Break any one and it can't happen.

The standard fix breaks the circular wait: **every thread acquires locks in the same global order** (always `A` then `B`). Change `t2` to lock `A` first and the hang disappears.

Other rules:
- Never lock the same default mutex twice in one thread. It deadlocks or is undefined. A mutex created with `PTHREAD_MUTEX_ERRORCHECK` returns an error instead, which helps while debugging.
- Use `pthread_mutex_trylock` when you can back off and retry instead of blocking.
- Don't call unknown code (callbacks) while holding a lock, because it may try to take another lock.

---

# 11. Detached threads and attributes

If you don't need the result and don't want to join, **detach** the thread:

```c
pthread_t t;
pthread_create(&t, NULL, Worker, arg);
pthread_detach(t);     // resources are freed automatically when it exits
```

A detached thread can't be joined. Also make sure its `arg` stays valid. Passing the address of a local variable of a function that then returns gives the thread a dangling pointer. Use heap memory that the thread frees. And `main` must stay alive (or call `pthread_exit`) for it to finish, so detaching suits long-lived background workers, not a replacement for join.

## Attributes

An attribute object sets options when the thread is created:

```c
pthread_attr_t attr;
pthread_attr_init(&attr);
pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
pthread_attr_setstacksize(&attr, 256 * 1024);   // smaller stack for many threads

pthread_t t;
pthread_create(&t, &attr, Worker, NULL);
pthread_attr_destroy(&attr);
```

Each thread reserves a stack up front, so spawning thousands of threads with the default size uses a lot of virtual memory and can hit limits. Deep recursion or large local arrays in a thread overflow its stack the same way as in `main`, just at a smaller size.

---

# 12. Error handling

pthread functions return the error number directly (`0` means success) and **do not set `errno`**. So `perror()` after a failed `pthread_create` prints a misleading message.

```c
#include <string.h>

int rc = pthread_create(&t, NULL, Worker, NULL);
if (rc != 0) {
    fprintf(stderr, "pthread_create: %s\n", strerror(rc));
    exit(EXIT_FAILURE);
}
```

The same applies to `pthread_join`, `pthread_mutex_lock` and the rest. A common failure is `EAGAIN` (a thread or memory limit was reached). On older systems `strerror` may not be thread-safe, and `strerror_r` is the safe variant, but modern glibc and macOS are fine.

---

# 13. Thread-local data and thread safety

## Thread-local variables

A thread-local variable looks global, but each thread gets its own copy:

```c
_Thread_local int counter = 0;    // C11; GCC/Clang also accept __thread

void *Worker(void *arg) {
    (void)arg;
    counter++;                    // changes only this thread's copy
    printf("%d\n", counter);      // 1, in every thread
    return NULL;
}
```

This avoids locking for per-thread state, such as a buffer or an error code.

## Thread-unsafe library functions

Some standard functions keep hidden shared state, so calling them from several threads is a race. Use the reentrant versions (they take a buffer or state from you):

| Unsafe | Use instead |
| --- | --- |
| `strtok` | `strtok_r` |
| `localtime`, `gmtime` | `localtime_r`, `gmtime_r` |
| `rand` | `rand_r` (or a per-thread generator) |
| `asctime`, `ctime` | `asctime_r`, `ctime_r` |

`printf` is safe to call concurrently (each call locks the stream), but output from different threads can still interleave between calls. `malloc` and `free` are thread-safe.

---

# 14. Other primitives

| Tool | Use | Notes |
| --- | --- | --- |
| `pthread_rwlock_t` | many readers or one writer | helps when reads far outnumber writes |
| `sem_t` (`<semaphore.h>`) | counter that threads wait on and post to | `sem_init` is not supported on macOS |
| `pthread_once` | run an initializer exactly once | for lazy setup of shared state |
| `pthread_barrier_t` | all threads wait until everyone arrives | not available on macOS |
| `pthread_spinlock_t` | busy-wait lock for very short sections | wastes CPU if the lock is held long |

---

# 15. Debugging

```bash
gcc -g -O1 -fsanitize=thread -pthread main.c -o main && ./main
```

**ThreadSanitizer** reports data races with the stack traces of both accesses. Run the racy `count++` example from section 6 through it. You can't combine `-fsanitize=thread` with `-fsanitize=address` in one build.

Other tools:
- `valgrind --tool=helgrind` is slower and also finds lock-order problems. It works well on Linux.
- Hangs: in gdb, `info threads` and `thread apply all bt` show where every thread is stuck. That is the first thing to run on a deadlock. In lldb, use `thread list` and `bt all`.
- `top -H -p <pid>` (Linux) shows each thread's CPU use.

---

# 16. Coming from Go

| pthreads | Go |
| --- | --- |
| `pthread_create` | `go f()` |
| `pthread_join` | `sync.WaitGroup` (`Wait`) |
| `pthread_mutex_t` | `sync.Mutex` |
| `pthread_cond_t` | `sync.Cond` |
| `atomic_*` | `sync/atomic` |
| queue and condition variables | buffered channel |
| thread pool | worker pool reading from a channel |
| `-fsanitize=thread` | `go run -race` |

The key difference is that Go's runtime handles scheduling and gives you channels, while in C you build the queue, the shutdown logic and the error handling yourself. Writing the pool in section 9 shows what a buffered channel does for you.

---

# 17. Platform notes

- `gettid()` and `syscall(SYS_gettid)` are Linux only. On macOS use `pthread_threadid_np`.
- macOS has no `pthread_barrier_t`, and `sem_init` (unnamed semaphores) isn't supported there. Use named semaphores via `sem_open`, or build the feature from a mutex and a condition variable.
- Secondary threads on macOS get a much smaller default stack than Linux's 8 MB, so deep recursion in a thread fails sooner. Check with `pthread_attr_getstacksize`.
- ThreadSanitizer works with Apple Clang. Helgrind support on recent macOS is limited.
- Windows has its own threading API and isn't covered here.

---

# 18. Exercises

## Exercise 01: three counters

Fix the racy counter three ways: a mutex around each increment, a local counter merged once, and an atomic. Time each with `time ./main` for 4 threads. Explain the difference.

## Exercise 02: parallel sum

Sum a large array with N threads. Give each thread its own chunk with a struct, so no `malloc` is needed for results:

```c
typedef struct {
    const int *arr;
    size_t lo, hi;    // half-open range [lo, hi)
    long result;
} job_t;
```

`main` owns the array of `job_t`, joins every thread, and adds up the `result` fields. Think about why no lock is needed.

## Exercise 03: general thread pool

Extend section 9: tasks carry a function pointer and an argument, the pool size is a parameter, and there is a `pool_destroy` that finishes pending work and joins every worker. Decide what happens to `push` after shutdown.

## Exercise 04: producer and consumer with semaphores

Rebuild the bounded queue using two semaphores (`empty` and `full`) plus one mutex. Compare it with the condition-variable version.

## Exercise 05: dining philosophers

Five threads, five forks (mutexes). First write a version that deadlocks, then fix it with lock ordering.

## Exercise 06: concurrent server

Take the TCP server from your sockets note. First handle each client in a new thread, then switch to a fixed thread pool that takes accepted descriptors from a queue. Watch for who closes each client descriptor.

---

# Cheat sheet

| Need | API |
| --- | --- |
| Start a thread | `pthread_create` |
| Wait and get a result | `pthread_join` |
| Fire and forget | `pthread_detach` |
| Protect shared data | `pthread_mutex_lock` / `unlock` |
| Single counter or flag | `atomic_*` from `<stdatomic.h>` |
| Wait for a condition | `pthread_cond_wait` in a `while` loop |
| Wake one or all | `pthread_cond_signal` / `pthread_cond_broadcast` |
| Per-thread variable | `_Thread_local` |
| Check for races | `-fsanitize=thread` |

# Common mistakes

1. Using `sleep()` in `main` instead of `pthread_join`.
2. Passing `&i` of a loop variable to a thread.
3. Returning the address of a stack variable from a thread.
4. Forgetting to initialize `malloc`'d memory in the thread.
5. Touching shared state without the lock on one code path.
6. `if` instead of `while` around `pthread_cond_wait`.
7. `signal` instead of `broadcast` on shutdown.
8. Inconsistent lock ordering across threads.
9. Using `perror` for pthread errors.
10. Using a plain `int` as a stop flag instead of an atomic or a locked variable.
11. Calling thread-unsafe functions like `strtok` from several threads.
12. Passing the address of a local to a detached thread.

# Summary

- Threads share the heap and globals but each has its own stack.
- Create with `pthread_create`, and always join (or detach) each thread.
- Shared data needs a mutex or an atomic. A race is undefined behavior, not just a wrong number.
- Use condition variables to wait for a state, and always check the condition in a `while` loop.
- Take locks in one fixed order to avoid deadlock.
- pthread functions return error codes and don't set `errno`.
- Test with ThreadSanitizer.