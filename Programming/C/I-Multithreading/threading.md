---
title: "Multithreading in C"
date: 2026-10-25
tags:
  - c
  - pthreads
  - concurrency
  - systems
---

# Overview

POSIX threads (pthreads) in C: create, join, share data safely, and build a small thread pool. Every example compiles with:

```bash
gcc -Wall -Wextra -pthread main.c -o main
```

`-pthread` is required at both compile and link time.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | Creating threads | `pthread_create`, the `void*(void*)` signature |
| 2 | Waiting for threads | `pthread_join`, why `sleep()` in `main` is wrong |
| 3 | Thread IDs | `pthread_self` |
| 4 | Passing arguments | the `&i` loop pitfall |
| 5 | Returning values | heap result + `pthread_join` |
| 6 | Race conditions | lost updates on `count++` |
| 7 | Mutex | `pthread_mutex_t`, cost, batching |
| 8 | Atomics | `stdatomic.h` |
| 9 | Condition variables | bounded queue + thread pool |
| 10 | Deadlock | lock ordering |
| 11 | Detached threads | `pthread_detach` |
| 12 | Error handling | pthread functions don't set `errno` |
| 13 | Debugging | ThreadSanitizer |

A pthread is a 1:1 kernel thread (scheduled by the OS). Go's goroutines are M:N: many goroutines multiplexed over a few OS threads by the runtime's GMP scheduler. Same "spawn and run concurrently" feel, very different cost model: a pthread costs ~8 MB of virtual stack by default, a goroutine starts at ~2 KB.

# 1. Creating threads

```c
int pthread_create(pthread_t *thread, const pthread_attr_t *attr,
                   void *(*start)(void *), void *arg);
```

The thread function must be `void* f(void*)`.

```c
#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

void* AliceTurn(void* args) {
   for (int i = 0; i <= 10; i++) {
      sleep(1);
      printf("Alice = %d\n", i);
   }
   return NULL;
}

void* BobTurn(void* args) {
   for (int i = 0; i <= 4; i++) {
      sleep(1);
      printf("Bob = %d\n", i);
   }
   return NULL;
}

int main() {
   pthread_t alice, bob;
   pthread_create(&alice, NULL, AliceTurn, NULL);
   pthread_create(&bob, NULL, BobTurn, NULL);

   sleep(10); // BAD, see next section
   return 0;
}
```

When `main` returns, the process exits and **every thread is killed**, finished or not. Same as Go: `main` exiting kills all goroutines.

# 2. Waiting for threads: `pthread_join`

The `sleep(10)` above is a race against the workers. Alice needs 11 seconds (0..10), so she gets killed before her last print. Guessing a duration is never correct; join instead.

```c
int main() {
   pthread_t alice, bob;
   pthread_create(&alice, NULL, AliceTurn, NULL);
   pthread_create(&bob, NULL, BobTurn, NULL);

   pthread_join(bob, NULL);   // blocks until Bob returns
   pthread_join(alice, NULL); // blocks until Alice returns
   return 0;
}
```

Join order doesn't matter for correctness: joining an already-finished thread returns immediately. Every joinable thread must be joined (or detached), otherwise its resources leak until process exit.

Alternative: calling `pthread_exit(NULL)` at the end of `main` lets the other threads keep running after `main` is done.

# 3. Thread IDs

```c
void* Calc(void* args) {
   printf("Calc() = %lu\n", (unsigned long) pthread_self());
   return NULL;
}

int main() {
   printf("main() = %lu\n", (unsigned long) pthread_self());

   pthread_t t[10];
   for (int i = 0; i < 10; i++) pthread_create(&t[i], NULL, Calc, NULL);
   for (int i = 0; i < 10; i++) pthread_join(t[i], NULL);

   for (int i = 0; i < 10; i++)
      printf("%d ThreadId = %lu\n", i, (unsigned long) t[i]);
   return 0;
}
```

- `pthread_t` is opaque. The `unsigned long` cast works on Linux (it's a pointer-sized value) but isn't portable. Compare IDs with `pthread_equal(a, b)`.
- This is the pthread ID, not the kernel TID that `top -H` / `gdb` shows. For that: `syscall(SYS_gettid)` (or `gettid()` on glibc 2.30+).
- Output order of the `Calc()` lines is nondeterministic. Don't rely on it.

# 4. Passing arguments

The `void* arg` is a pointer, so what it points to must outlive the thread's read of it.

```c
void* Worker(void* arg) {
   int id = *(int*) arg;
   printf("worker %d\n", id);
   return NULL;
}

int main() {
   pthread_t t[4];
   int ids[4];

   for (int i = 0; i < 4; i++) {
      ids[i] = i;                                   // one slot per thread
      pthread_create(&t[i], NULL, Worker, &ids[i]);
   }
   for (int i = 0; i < 4; i++) pthread_join(t[i], NULL);
   return 0;
}
```

**Pitfall:** `pthread_create(&t[i], NULL, Worker, &i)` is a bug. All threads share the same `i`, which keeps changing (and goes out of scope after the loop). You'll see duplicate or garbage IDs.

Cheap trick for a single integer: `(void*)(intptr_t) i` and `(int)(intptr_t) arg` on the other side. For more than one value, pass a pointer to a struct (heap-allocated or an array element, as above).

# 5. Returning values

A thread returns a `void*`, and `pthread_join` hands it back through a `void**`. Returning a pointer to a local variable is a use-after-return, so allocate on the heap (or write into a caller-owned struct).

```c
void* Sum(void* arg) {
   int n = *(int*) arg;
   int* res = malloc(sizeof(int));
   *res = 0;                          // malloc doesn't zero memory
   for (int i = 1; i <= n; i++) *res += i;
   return res;
}

int main() {
   int n = 100;
   pthread_t t;
   pthread_create(&t, NULL, Sum, &n);

   int* result;
   pthread_join(t, (void**) &result);
   printf("%d\n", *result);           // 5050
   free(result);                      // joiner owns the allocation
   return 0;
}
```

# 6. Race condition

```c
int count;

void* Calc(void* args) {
   for (int i = 1; i <= 10000000; i++) count++;
   return NULL;
}

int main() {
   count = 0;
   pthread_t t;
   pthread_create(&t, NULL, Calc, NULL);
   Calc(NULL);                 // main works too
   pthread_join(t, NULL);
   printf("%d\n", count);      // expected 20000000, got something smaller
   return 0;
}
```

`count++` is three steps: load, add, store. Two threads can both load `5`, both store `6`, and one increment is lost. The result differs on every run. (With `-O2` the compiler may also keep `count` in a register for the whole loop, which makes things worse. It's undefined behavior either way.)

# 7. Mutex

```c
pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
int count;

void* Calc(void* args) {
   for (int i = 1; i <= 10000000; i++) {
      pthread_mutex_lock(&lock);
      count++;
      pthread_mutex_unlock(&lock);
   }
   return NULL;
}
```

Now it always prints `20000000`. Rules:

- Everything that touches the shared variable must take the same lock. One unguarded access breaks it.
- Keep critical sections tiny, and never block (I/O, `sleep`) while holding a lock.
- Always unlock on every exit path (early `return`, `break`, error branches).
- Dynamic init: `pthread_mutex_init(&m, NULL)` paired with `pthread_mutex_destroy(&m)`.

**Cost:** locking 20M times is slow because the threads contend. Reduce lock traffic by accumulating locally and merging once:

```c
void* Calc(void* args) {
   int local = 0;
   for (int i = 1; i <= 10000000; i++) local++;

   pthread_mutex_lock(&lock);
   count += local;
   pthread_mutex_unlock(&lock);
   return NULL;
}
```

# 8. Atomics

For a plain counter, a lock is overkill. C11 atomics use a single hardware instruction (`lock xadd` on x86):

```c
#include <stdatomic.h>

atomic_int count;

void* Calc(void* args) {
   for (int i = 1; i <= 10000000; i++)
      atomic_fetch_add(&count, 1);
   return NULL;
}
```

Atomics only protect a single variable. Invariants that span several variables (e.g. a queue's `size` + `readAt` + `insertAt`) still need a mutex.

# 9. Condition variables: bounded queue + thread pool

A mutex answers "who may touch the data". A condition variable answers "wait until the data is in the state I need". The producer pushes tasks into a fixed-size ring buffer. Workers pop and execute them. Full queue blocks the producer, empty queue blocks the workers.

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
static pthread_cond_t  notEmpty = PTHREAD_COND_INITIALIZER; // workers wait on this
static pthread_cond_t  notFull  = PTHREAD_COND_INITIALIZER; // producer waits on this

void push(Task task) {
   pthread_mutex_lock(&lock);
   while (size == QUEUE_CAP)                  // full: wait for a worker to pop
      pthread_cond_wait(&notFull, &lock);

   queue[insertAt] = task;
   insertAt = (insertAt + 1) % QUEUE_CAP;
   size++;

   pthread_cond_signal(&notEmpty);            // wake one worker
   pthread_mutex_unlock(&lock);
}

// returns 1 and fills *t, or 0 when the queue is drained and closed
int pop(Task* t) {
   pthread_mutex_lock(&lock);
   while (size == 0 && !finish)               // empty: wait for push or shutdown
      pthread_cond_wait(&notEmpty, &lock);

   if (size == 0 && finish) {                 // nothing left, pool is closing
      pthread_mutex_unlock(&lock);
      return 0;
   }

   *t = queue[readAt];
   readAt = (readAt + 1) % QUEUE_CAP;
   size--;

   pthread_cond_signal(&notFull);             // wake the producer
   pthread_mutex_unlock(&lock);
   return 1;
}

void executeTask(Task task) {
   printf("ThreadId = %lu  ->  %d + %d = %d\n",
          (unsigned long) pthread_self(), task.a, task.b, task.a + task.b);
}

void* startThread(void* args) {
   Task t;
   while (pop(&t)) executeTask(t);
   return NULL;
}

int main() {
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

- **`while`, not `if`, around `pthread_cond_wait`.** Wakeups can be spurious, and another thread may grab the item between the signal and your wakeup. Always re-check the predicate.
- `pthread_cond_wait` atomically releases the mutex and sleeps, then re-acquires the mutex before returning. That's why you must hold the lock when calling it.
- **Shutdown uses `broadcast`**, not `signal`. Every idle worker must wake and see `finish`. A single `signal` would leave the rest blocked forever and `pthread_join` would hang.
- Workers drain the queue before exiting: `finish` only stops them when `size == 0`.
- Use `#define` (or `enum`) for array sizes in C. `const int N = 10; int a[N];` at file scope does **not** compile in C, because a `const` variable is not a constant expression (unlike C++).

Naming tip: name condition variables after the state you're waiting *for* (`notEmpty`, `notFull`), not the state you're leaving.

# 10. Deadlock

```c
// Thread 1                   // Thread 2
lock(&A);                     lock(&B);
lock(&B);  // waits for T2    lock(&A);  // waits for T1
```

Each thread holds one lock and waits forever for the other. Fix: **every thread acquires locks in the same global order** (always A then B). Other rules: never lock the same non-recursive mutex twice in one thread, and use `pthread_mutex_trylock` when you can back off instead of blocking.

# 11. Detached threads

If you don't need the result and don't want to join, detach:

```c
pthread_t t;
pthread_create(&t, NULL, Worker, &id);
pthread_detach(t);   // resources are freed automatically when it exits
```

A detached thread can't be joined. `main` must still stay alive long enough (or use `pthread_exit`) for it to finish, so this is for long-lived background workers, not a replacement for join.

# 12. Error handling

pthread functions return the error number directly (`0` on success) and **do not set `errno`**. So `perror()` after a failed `pthread_create` prints a misleading message.

```c
int rc = pthread_create(&t, NULL, Worker, NULL);
if (rc != 0) {
   fprintf(stderr, "pthread_create: %s\n", strerror(rc));
   exit(EXIT_FAILURE);
}
```

Same for `pthread_join`, `pthread_mutex_lock`, etc. A common failure is `EAGAIN` (thread or memory limit reached).

# 13. Debugging

```bash
gcc -g -fsanitize=thread -pthread main.c -o main && ./main
```

ThreadSanitizer reports data races with both stack traces. Run the racy `count++` example from section 6 through it. Helgrind (`valgrind --tool=helgrind`) is the slower alternative and also detects lock-order inversions. In `gdb`, `info threads` and `thread apply all bt` show where every thread is stuck, which is the first thing to run on a hang.

# Cheat sheet

| Need | API |
| --- | --- |
| Start a thread | `pthread_create` |
| Wait / get result | `pthread_join` |
| Fire and forget | `pthread_detach` |
| Protect shared data | `pthread_mutex_lock` / `unlock` |
| Single counter or flag | `atomic_*` from `<stdatomic.h>` |
| Wait for a condition | `pthread_cond_wait` in a `while` loop |
| Wake one / all | `pthread_cond_signal` / `pthread_cond_broadcast` |
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