---
title: "Memory Management in C"
date: 2026-10-25
tags:
  - c
  - memory
---

# Overview

In C, memory is managed manually. The compiler handles the stack for you, but you allocate and free heap memory yourself, and the language will not stop you from getting it wrong.

| Stack | Heap |
| --- | --- |
| Automatic: created and removed with each function call | Manual: `malloc` to get it, `free` to give it back |
| Fast, small (about 8 MB by default) | Slower, limited only by available memory |
| Dies when the function returns | Lives until you free it |
| Failure: stack overflow | Failure: leaks, use-after-free, double free |

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [Static variables](#static-variables) | Variables that survive between calls |
| 2 | [Stack overflow](#stack-overflow) | Two ways to blow the stack |
| 3 | [Stack frames](#stack-frames) | Frame addresses, return addresses, frame size |
| 4 | [Heap allocation](#heap-allocation) | `malloc`, `calloc`, `realloc`, `free` |
| 5 | [Heap bugs](#heap-bugs) | Leaks, use-after-free, double free, tools |
| 6 | [Memory functions](#memory-functions) | `memset`, `memcpy`, `memmove`, `memcmp` |
| 7 | [Object pool](#object-pool) | Reusing fixed-size objects |
| 8 | [Free list](#free-list) | A pool built on a linked list |
| 9 | [Exercise](#exercise) | Build a block allocator |

---

# Static variables

A `static` local variable is created once, keeps its value between calls, and is visible only inside its function.

```c
#include <stdio.h>

void test(void) {
    static int a = 0;   // initialized once, before main runs
    a += 10;
    printf("a = %d\n", a);
}

int main(void) {
    test();   // a = 10
    test();   // a = 20
    test();   // a = 30
    return 0;
}
```

- It does **not** live on the stack. A `static` initialized to a non-zero value goes in the **data** segment, and one initialized to zero (like above) usually goes in **BSS**.
- The `= 0` is not executed on each call. The value is set up once, when the program loads.
- Scope and lifetime are different things. The *scope* is the function, but the *lifetime* is the whole program.
- Because it outlives the call, returning its address is valid, unlike returning the address of a normal local:

```c
int *counter(void) {
    static int n = 0;
    n++;
    return &n;          // valid: n is not on the stack
}
```

The cost is shared state: all callers share one copy, and it is not thread-safe.

---

# Stack overflow

The stack has a fixed maximum size. On Linux the default is usually 8 MB (check with `ulimit -s`, which prints KB), and macOS is similar. Using more than that crashes the program with a segmentation fault.

Run these two examples **one at a time**, because the first one crashes before the second runs.

## Overflow 01: a huge local array

```c
#include <stdio.h>

int main(void) {
    const int size = 1024 * 1024 * 10;   // 10 MB, bigger than an 8 MB stack
    char buf[size];                      // allocated on the stack
    buf[0] = 'A';
    buf[1] = 'C';
    buf[size - 1] = '\0';
    printf("%c\n", buf[0]);
    return 0;
}
```

Why it crashes: the array is bigger than the whole stack, so touching it reaches memory the OS never gave to the stack. That access causes the segmentation fault.

Note that `const int size` is not a compile-time constant in C, so `buf` is a **variable length array (VLA)**, sized at runtime. VLAs on the stack are risky for exactly this reason. For big buffers use `malloc`, or make the array `static`.

## Overflow 02: infinite recursion

```c
#include <stdio.h>

void print(int x) {
    printf("%p = %d\n", __builtin_frame_address(0), x);
    print(x + 1);    // no base case, so it never stops
}

int main(void) {
    print(0);
    return 0;
}
```

Every call pushes a new stack frame, and nothing is ever popped. Eventually the stack runs out. Watch the printed addresses decrease, since the stack grows toward lower addresses.

Compile with `-O0`. With optimization (`-O2`) the compiler may turn this call into a loop (**tail-call optimization**, reusing the same frame), and then the program runs forever without crashing.

---

# Stack frames

Each function call gets its own **stack frame**, a block of stack memory for that call's local variables and bookkeeping. A typical x86-64 layout (it varies by compiler and CPU):

```text
higher addresses
+------------------------+
| caller's frame         |
+------------------------+
| return address         |  where to resume in the caller (an address in text)
+------------------------+
| saved frame pointer    |  <- __builtin_frame_address(0) usually points here
+------------------------+
| local variables        |
+------------------------+
lower addresses            <- stack pointer
```

`__builtin_frame_address` and `__builtin_return_address` are **GCC/Clang extensions**, not standard C. The argument is the level:
- `0` is the current function.
- `1` is the caller, `2` is the caller's caller.
- Levels above 0 are unreliable and can crash, so stay with `0`.

## Frame and return addresses

```c
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

void call_me(void *caller_frame) {
    void *mine = __builtin_frame_address(0);
    printf("call_me frame  = %p\n", mine);                 // lower than main's
    printf("call_me return = %p\n", __builtin_return_address(0));   // inside main, in the text segment
    printf("distance       = %" PRIuPTR " bytes\n",
           (uintptr_t)caller_frame - (uintptr_t)mine);
}

int main(void) {
    void *frame = __builtin_frame_address(0);
    printf("main frame     = %p\n", frame);
    call_me(frame);
    return 0;
}
```

Compile with `gcc -O0 -fno-omit-frame-pointer`. `uintptr_t` is an integer type big enough to hold a pointer, which is what you want for address arithmetic. Print it with `PRIuPTR`, not `%lu`, so the format matches on every platform.

To prove the return address is in text, ask gdb: `(gdb) info symbol <address>` should answer `main + <offset>`.

## Measuring a frame

```c
#include <stdio.h>
#include <stdint.h>

static void stack_frame(void) {
    static uintptr_t prev = 0;
    uintptr_t current = (uintptr_t)__builtin_frame_address(0);
    if (prev != 0) {
        printf("diff = %ld\n", (long)(intptr_t)(prev - current));
    }
    prev = current;
}

void A(void) {
    int x = 10;
    (void)x;
    printf("A() = %p\n", __builtin_frame_address(0));
    stack_frame();       // called from deeper: main -> D -> A -> stack_frame
}

void D(void) {
    A();
    printf("D() = %p\n", __builtin_frame_address(0));
    stack_frame();       // called from shallower: main -> D -> stack_frame
}

int main(void) {
    D();
    return 0;
}
```

The first `stack_frame` call happens deeper (inside `A`), so its frame is at a lower address. The second call comes straight from `D`, at a higher address. So `prev - current` is negative, and its size is roughly the size of `A`'s frame. That is the proof that the stack grows downward.

Exact numbers depend on the compiler, flags and alignment. Inlining can also remove the calls, so keep `-O0`.

---

# Heap allocation

All of these are in `<stdlib.h>`. They return memory from the heap, or `NULL` on failure.

## malloc

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int *ptr = malloc(sizeof *ptr);      // sizeof *ptr stays correct if the type changes
if (ptr == NULL) {
    perror("malloc");                // handle failure, don't just ignore it
    return 1;
}

printf("%d\n", *ptr);                // wrong: contents are indeterminate (garbage)
memset(ptr, 0, sizeof *ptr);         // every byte set to 0x00, so the int is 0
printf("Address: %p, Value: %d\n", (void *)ptr, *ptr);

free(ptr);
ptr = NULL;
```

- `malloc` does not clear the memory. Reading it before writing is a bug.
- In C you don't cast the result (`void *` converts automatically). In C++ you must.

## calloc

```c
int *cptr = calloc(3, sizeof *cptr);   // 3 ints, all zero
if (cptr == NULL) { /* handle failure */ }

for (int i = 0; i < 3; i++) {
    printf("%d\n", cptr[i]);           // 0 0 0
}
free(cptr);
```

`calloc(count, size)` is like `malloc` plus zeroing. It also checks that `count * size` doesn't overflow, which `malloc(count * size)` does not.

## realloc

`realloc` changes the size of an existing block. It either grows the block in place if there is room, or allocates a new block, copies the old contents, and frees the old block. After a successful call, treat the **old pointer as invalid**.

```c
int n = 5;
int *arr = malloc(n * sizeof *arr);
if (arr == NULL) return 1;

for (int i = 0; i < n; i++) {
    arr[i] = i * i;
}

printf("arr[3]      = %d\n", arr[3]);
printf("*(arr + 3)  = %d\n", *(arr + 3));       // same element
printf("before: %p\n", (void *)arr);

int m = n * 2;
int *tmp = realloc(arr, m * sizeof *arr);
if (tmp == NULL) {
    free(arr);          // realloc failed, but the OLD block is still valid, so free it
    return 1;
}
arr = tmp;              // success: the old pointer must not be used again
                        // and you must not free it, realloc already did if it moved

for (int i = n; i < m; i++) {
    arr[i] = i * i;     // the new half is uninitialized, fill it yourself
}

printf("after:  %p\n", (void *)arr);             // may or may not be the same address
free(arr);              // free once, through the final pointer
```

Always use a temporary like `tmp`. Writing `arr = realloc(arr, ...)` loses the only pointer to the old block if `realloc` fails, which is a leak.

Reading `arr[n]` before the resize (your original code did this) is an **out-of-bounds read**, because valid indexes are `0` to `n - 1`. It is undefined behavior: it may print garbage, crash, or silently corrupt the allocator's bookkeeping.

## free

```c
free(ptr);
ptr = NULL;     // optional, but makes a later use crash cleanly or be caught by a check
```

- `free(NULL)` is safe and does nothing, so `if (ptr) free(ptr);` is not needed.
- Only pass pointers returned by `malloc`, `calloc` or `realloc`, and only once.
- Avoid `realloc(p, 0)`. Its behavior differs between systems.

## How the heap works

`malloc` is a library function, not a system call. It manages a pool of memory that it requests from the OS, using `brk` (which grows the `[heap]` region) for small blocks and `mmap` for large ones (roughly 128 KB and up on Linux). The memory you get back is usually a bit larger than you asked for, and the allocator keeps hidden bookkeeping next to each block. This is why writing past the end of a block can corrupt things far from your bug.

---

# Heap bugs

| Bug | Example | Result |
| --- | --- | --- |
| Memory leak | `malloc` without `free` | Memory grows until the program ends |
| Use after free | `free(p); *p = 1;` | Undefined behavior, may corrupt other data |
| Double free | `free(p); free(p);` | Heap corruption, often a crash |
| Heap overflow | Writing `arr[n]` on an `n`-element array | Overwrites neighboring memory |
| Uninitialized read | Reading before writing | Garbage values |
| Freeing the wrong thing | `free` on a stack or offset pointer | Crash or corruption |

Habits that prevent most of them:
- Decide which function **owns** each allocation, and have that owner free it.
- Set pointers to `NULL` after `free`.
- Check every allocation result.

Tools:

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined main.c -o a.out
```

- **AddressSanitizer** (`-fsanitize=address`) reports overflows, use-after-free and double free at runtime, and reports leaks on Linux at exit. It works on Linux and macOS.
- `valgrind ./a.out` is a good Linux alternative (support on recent macOS is limited).
- On macOS, `leaks --atExit -- ./a.out` reports leaks.

---

# Memory functions

All in `<string.h>`. They work on raw bytes, whatever the type.

## memset

```c
int A[10];

memset(A, 0, sizeof A);    // all ints = 0
memset(A, -1, sizeof A);   // every byte 0xFF, so every int = -1
memset(A, 1, sizeof A);    // every byte 0x01, so every int = 0x01010101 = 16843009, not 1
```

`memset` fills **bytes**, not elements. It is only safe for ints when the byte pattern gives the value you want (`0` and `-1`). The value `1` does not give 1. That number is not "garbage": it is a predictable result.

## memcpy

```c
for (int i = 0; i < 10; i++) A[i] = i + 1;     // 1 2 3 ... 10

int B[5];
memcpy(B, A, sizeof B);                         // B = 1 2 3 4 5

memcpy(A + 4, A, 3 * sizeof *A);                // copy A[0..2] to A[4..6]
// A = 1 2 3 4 1 2 3 8 9 10   (the two ranges do not overlap, so this is fine)
```

`memcpy` requires that the source and destination do **not** overlap. That promise is what the `restrict` keyword in its signature means. (`restrict` is a keyword telling the compiler that two pointers never refer to the same memory, which lets it optimize.) If they overlap, the behavior is undefined.

## memmove

`memmove` handles overlap correctly, so use it whenever the ranges might overlap. (Your original string example copied `s[0..4]` to `s[5..9]`, which does not overlap, so it didn't show the difference.) A real overlap:

```c
char s[] = "helloworld";
memmove(s + 2, s, 5);       // copy "hello" two bytes to the right
printf("%s\n", s);          // hehellorld
```

With `memcpy` here, a simple forward byte copy would re-read bytes it has already overwritten and could produce `hehehehrld`. The exact result isn't guaranteed, because it is undefined.

## memcmp

```c
int cmp = memcmp(A, B, 5 * sizeof *A);   // size is in BYTES, so 5 ints = 5 * sizeof(int)
```

Returns 0 if equal, otherwise a negative or positive number (rely only on the sign). Your original `memcmp(A, B, 5)` compared only 5 bytes, not 5 ints. Don't use `memcmp` to compare structs (padding bytes may differ) or floats.

## Others

| Name | Does | Notes |
| --- | --- | --- |
| `memchr(p, c, n)` | first occurrence of a byte | standard |
| `memmem(h, hn, needle, nn)` | find a byte sequence | GNU/BSD extension, not standard |
| `memccpy(d, s, c, n)` | copy until byte `c` is found | POSIX, standard since C23 |
| `aligned_alloc(align, size)` | allocate with a given alignment | C11, `size` must be a multiple of `align` |

---

# Object pool

An object pool **reuses** objects of the same size instead of calling `malloc` and `free` each time. This avoids allocator overhead, avoids fragmentation, and makes the memory use predictable. (It is the same idea as `sync.Pool` in Go.)

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#define POOL_SIZE 10

typedef struct Pair {
    int first, second;
} pair_t;

typedef struct {
    bool allocated;
    pair_t obj;
} pool_elem_t;

static pool_elem_t object_pool[POOL_SIZE];   // static storage: starts zeroed, allocated == false
```

Use `#define` for the size. `const int POOL_SIZE = 10;` is not a constant expression in C, so it can't size a global array.

```c
// Find which pool slot a pair_t pointer belongs to, or -1 if it is not from the pool.
static long index_of(const pair_t *obj) {
    if (obj == NULL) return -1;
    uintptr_t pool_addr = (uintptr_t)object_pool;
    // obj points at the inner field; offsetof gives that field's byte offset
    // inside the struct, so subtracting it gives the start of the slot
    uintptr_t elem_addr = (uintptr_t)obj - offsetof(pool_elem_t, obj);
    if (elem_addr < pool_addr) return -1;                  // before the pool
    size_t off = elem_addr - pool_addr;
    if (off % sizeof(pool_elem_t) != 0) return -1;         // not on a slot boundary
    size_t idx = off / sizeof(pool_elem_t);
    return idx < POOL_SIZE ? (long)idx : -1;               // past the end
}

pair_t *Get(void) {
    for (int i = 0; i < POOL_SIZE; i++) {
        if (object_pool[i].allocated) continue;
        object_pool[i].allocated = true;
        memset(&object_pool[i].obj, 0, sizeof(pair_t));    // hand out a clean object
        return &object_pool[i].obj;
    }
    return NULL;                                           // pool is full
}

bool Put(pair_t *obj) {
    long idx = index_of(obj);
    if (idx < 0) return false;                             // not from this pool
    if (!object_pool[idx].allocated) return false;         // already returned (double Put)
    object_pool[idx].allocated = false;
    return true;
}

int main(void) {
    pair_t *a = Get();
    pair_t *b = Get();
    a->first = 1;
    a->second = 2;

    Put(a);
    pair_t *c = Get();
    printf("reused slot: %d\n", c == a);   // 1: the freed slot was handed out again

    Put(b);
    Put(c);
    return 0;
}
```

Your original `main` called `Get()` twice into the same variable, which lost the first pointer and never returned it. The version above keeps each pointer.

Things to know:
- The checks in `Put` matter because the caller can pass any pointer. Without them, a bad pointer would corrupt the pool.
- `Get` scans the array, so it is O(n). The free list below makes it O(1).
- Pointers are converted to `uintptr_t` before comparing, because comparing unrelated pointers with `<` is not well defined in C.
- It is not thread-safe. Concurrent callers would need a lock.

---

# Free list

Instead of scanning for a free slot, keep the free slots in a **linked list**. `Get` takes the head of the list, and `Put` pushes a slot back onto the head. Both are O(1).

```c
typedef struct PoolElem {
    pair_t obj;
    struct PoolElem *next;     // link to the next free slot
} pool_elem_t;

static pool_elem_t object_pool[POOL_SIZE];
static pool_elem_t *free_list = NULL;

void init_pool(void) {
    for (int i = 0; i + 1 < POOL_SIZE; i++) {
        object_pool[i].next = &object_pool[i + 1];
    }
    object_pool[POOL_SIZE - 1].next = NULL;
    free_list = &object_pool[0];
}

pair_t *Get(void) {
    if (free_list == NULL) return NULL;           // no free object
    pool_elem_t *elem = free_list;
    free_list = elem->next;                       // pop the head
    memset(&elem->obj, 0, sizeof elem->obj);
    return &elem->obj;
}

bool Put(pair_t *obj) {
    long idx = index_of(obj);                     // same helper as before
    if (idx < 0) return false;
    pool_elem_t *elem = &object_pool[idx];
    elem->next = free_list;                       // push on the head
    free_list = elem;
    return true;
}
```

Limitation: this version has no `allocated` flag, so a double `Put` on the same pointer inserts the slot twice and makes the list loop back on itself. Add a flag or accept that the caller must not do this.

## Saving memory with a union

Each slot above stores a `next` pointer even while the object is in use. A **union** lets both share the same memory, since a slot is either free (a link) or in use (an object), never both:

```c
typedef union Slot {
    pair_t obj;            // when allocated
    union Slot *next;      // when free
} slot_t;
```

The size of the union is the size of its largest member, so there is no per-slot overhead. This is how many real allocators work: **the free block itself stores the link**. The downside is that a freed object's contents are overwritten, so using it after `Put` is a bug.

---

# Exercise

Build a **fixed-size block allocator** (a memory pool) using a linked list.

```c
typedef struct Block {
    struct Block *next;
} block_t;
```

Suggested interface:

```c
typedef struct {
    block_t *free_list;
    /* add whatever else you need */
} pool_t;

bool  pool_init(pool_t *p, void *buffer, size_t block_size, size_t count);
void *pool_alloc(pool_t *p);
void  pool_free(pool_t *p, void *block);
```

Requirements and hints:
- The caller gives you one big `buffer`. You carve it into `count` blocks and link them all into the free list during `pool_init`.
- A free block must be able to hold a `block_t`, so `block_size` must be at least `sizeof(block_t)`. Reject smaller sizes.
- Round `block_size` up to the alignment of `max_align_t` (from `<stddef.h>`), so every block is correctly aligned.
- `pool_free` should reject pointers that are not inside the buffer or not on a block boundary, as in `index_of` above.
- Test with AddressSanitizer: allocate everything, free everything, allocate again, and check you get the same addresses back.
- Bonus: grow the pool by adding a new chunk when the free list is empty. Count how many blocks are in use. Add a lock for thread safety.

```c
#include <string.h>

memmove(destination, source, size_of_memory);
memmove(&arr[idx + 1], &arr[idx], sizeof(int) * 3);   // shift 3 ints one position to the right
```

---

# Summary

- Stack memory is automatic and small. Heap memory is manual and long-lived.
- A `static` local lives for the whole program, in data or BSS, not on the stack.
- Large locals and unbounded recursion cause stack overflow.
- Check every allocation, free exactly once, and use a temporary with `realloc`.
- `memset` fills bytes, `memcpy` forbids overlap, and `memmove` allows it.
- Pools and free lists make allocation O(1) and predictable.
- Compile with `-fsanitize=address,undefined` while learning.