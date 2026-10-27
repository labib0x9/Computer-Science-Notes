---
title: "Pointers in C Programming"
date: 2026-10-27
tags:
  - c
  - memory
  - pointers
---

# Overview

Every variable lives in memory, so how do we work with that memory directly? With pointers. They are one of the most important topics in C, and they explain a lot about how memory behaves.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [What is a pointer](#what-is-a-pointer) | Addresses, `&`, `*`, `NULL` |
| 2 | [Pass by value](#pass-by-value) | Why functions get copies, and how pointers fix it |
| 3 | [Arrays, strings and `const`](#arrays-strings-and-const) | Pointer arithmetic, `char*` vs `char[]` |
| 4 | [Swap](#swap) | Pointer swap, generic swap |
| 5 | [Pointer to pointer](#pointer-to-pointer) | `char **`, out-parameters |
| 6 | [void pointer](#void-pointer) | Generic pointers and casting |
| 7 | [Function pointers](#function-pointers) | Callbacks, dispatch tables |
| 8 | [Dangling pointers and escape to heap](#dangling-pointers-and-escape-to-heap) | Returning pointers safely |
| 9 | [Common mistakes](#common-mistakes) | Bugs and the tools that catch them |
| 10 | [Exercises](#exercises) | Practice |

---

# What is a pointer

A pointer is a variable that stores a **memory address**, meaning the location of another value.

```c
int x = 20;
int *p = &x;                 // p holds the address of x

printf("%d\n", x);           // 20
printf("%p\n", (void *)p);   // the address of x
printf("%d\n", *p);          // 20, reading through the pointer
*p = 30;                     // writing through the pointer
printf("%d\n", x);           // 30, x changed
```

Picture of memory (the addresses are made up):

```text
Address   Value     Name
0x1000    30        x
0x1008    0x1000    p      <- p points to x
```

Two operators and one declaration symbol:

| Symbol | Name | Meaning |
| --- | --- | --- |
| `&x` | address-of | the address of `x` |
| `*p` | dereference | the value stored at the address `p` holds |
| `int *p` | declaration | `*` here means "p is a pointer to int" |

The `*` means two different things depending on where it appears: "is a pointer" in a declaration, and "follow the pointer" in an expression.

```c
char *s;
int *a;
void *b;
```

Gotcha: the `*` belongs to the variable, not the type.

```c
int* a, b;     // a is a pointer, b is a plain int
int *a, *b;    // both pointers, so write it this way
```

Facts to remember:
- On a 64-bit system every pointer is 8 bytes, whatever it points to (`sizeof(int *)` and `sizeof(char *)` are both 8).
- The type tells the compiler how many bytes to read when you dereference, and how far `p + 1` moves.
- Print addresses with `%p`, and cast to `void *`.

## NULL

`NULL` is a pointer value that means "points to nothing". Dereferencing it is **undefined behavior** (the C standard does not say what happens, and in practice the program usually crashes with a segmentation fault).

```c
int *p = NULL;
if (p != NULL) {
    printf("%d\n", *p);
}
```

Always initialize pointers, either to a real address or to `NULL`. An uninitialized pointer holds garbage and points to a random place.

---

# Pass by value

C passes arguments **by value**: a function receives a *copy* of each argument.

```c
void set(int x) {   // x is a copy
    x = 10;
}

int x = 20;
set(x);
printf("%d\n", x);  // 20, unchanged
```

To change the original, pass its address. The function still gets a copy, but it is a copy of the *address*, so it reaches the same variable.

```c
void set(int *x) {  // copy of the address of x
    *x = 10;        // dereference: write to the original
}

int x = 20;
set(&x);
printf("%d\n", x);  // 10
```

## Why not just copy everything?

Copying large data is slow and wastes memory. Pointers let you pass a small address (8 bytes) instead of the data.

Note: in C, arrays are never copied when passed to a function (see the next section). The copying problem shows up with large **structs**:

```c
typedef struct {
    int data[1000000];   // about 4 MB
} big_t;

void by_value(big_t b);     // copies 4 MB on every call, and may overflow the stack
void by_pointer(big_t *b);  // copies 8 bytes
```

The stack is small (typically around 8 MB), so a few large copies can cause a stack overflow.

---

# Arrays, strings and const

## Arrays and pointer arithmetic

An array name becomes a pointer to its first element in most expressions. This is called **decay**.

```c
int arr[3] = {10, 20, 30};
int *p = arr;               // same as &arr[0]

printf("%d\n", *p);         // 10
printf("%d\n", *(p + 1));   // 20
printf("%d\n", p[2]);       // 30, p[i] is the same as *(p + i)
```

`p + 1` moves by `sizeof(int)` bytes, not by one byte. Going outside the array is undefined behavior.

Because of decay, a function does not know the array's length:

```c
void print_size(int arr[]) {              // arr is really "int *arr" here
    printf("%zu\n", sizeof(arr));         // 8, not the array size
}
```

So pass the length separately:

```c
int sum(const int *arr, size_t n) {
    int total = 0;
    for (size_t i = 0; i < n; i++) total += arr[i];
    return total;
}
```

## Strings

A string is just a `char` array ending with `'\0'`, and a `char *` often points at one.

```c
char *s1 = "hi";     // points to read-only memory (rodata)
char s2[] = "hi";    // a writable copy; if local, it lives on the stack

s2[0] = 'H';         // ok
s1[0] = 'H';         // crash: writing to read-only memory
```

Prefer `const char *s1 = "hi";` so the compiler stops you from writing to a literal.

## const and pointers

Read the declaration right to left:

```c
const int *p;        // pointer to const int: cannot change *p, can change p
int *const q = &x;   // const pointer to int: can change *q, cannot change q
```

A parameter like `const int *arr` promises the function will not modify the data.

---

# Swap

To swap two integers we need their addresses.

## Common mistake

```c
void swap(int *a, int *b) {
    int *temp = a;   // BUG
    a = b;
    b = a;
}
```

This does nothing to the caller's variables:
- `a` and `b` are local copies of the addresses, so reassigning them changes only the copies.
- `b = a` runs after `a = b`, so both end up with the same address.
- `temp` should hold an `int`, not an address.

## Correct version

```c
void swap(int *a, int *b) {
    int temp = *a;
    *a = *b;
    *b = temp;
}

int a = 10, b = 20;
swap(&a, &b);
printf("a = %d, b = %d\n", a, b);   // a = 20, b = 10
```

## Swapping any type

Now swap two `struct` values. The function does not know the type, only the **size in bytes**, so it uses `void *` and copies raw bytes.

```c
#include <stdlib.h>
#include <string.h>

typedef struct {
    char name[10];
    int age;
} node_t;

void swap(void *a, void *b, size_t size) {
    if (a == NULL || b == NULL || size == 0) return;

    void *tmp = malloc(size);
    if (tmp == NULL) return;

    memmove(tmp, a, size);   // tmp = a
    memmove(a, b, size);     // a = b
    memmove(b, tmp, size);   // b = tmp

    free(tmp);
}

int main(void) {
    node_t u1 = {.name = "labib", .age = 24};
    node_t u2 = {.name = "usha", .age = 21};

    swap(&u1, &u2, sizeof(node_t));
    printf("u1 = %s %d\n", u1.name, u1.age);   // u1 = usha 21
    printf("u2 = %s %d\n", u2.name, u2.age);   // u2 = labib 24
}
```

`memmove` is used because it is safe even if the two memory ranges overlap. `memcpy` is not.

## What is the output?

```c
char *a = "afafadf";
char *b = "xyz";
swap(&a, &b, sizeof(char *));
printf("a = %s, b = %s\n", a, b);
```

This prints `a = xyz, b = afafadf`. We passed `&a` and `&b` (the addresses of the *pointers*) and swapped 8 bytes each. The strings themselves never move, only the pointers to them.

```c
int int_num = 5;
float float_num = 10.45;
swap(&int_num, &float_num, sizeof(int));
printf("int_num = %d, float_num = %f\n", int_num, float_num);
```

This compiles and runs, but the result is nonsense. `swap` moves raw bytes and knows nothing about types, so each variable now holds the other's *bit pattern*. With IEEE 754 floats, `int_num` becomes a large number (about 1093088051, so run it to check on your machine) and `float_num` prints `0.000000`, because the bits of integer `5` are a tiny float. It only "works" because `sizeof(int) == sizeof(float)` on most platforms.

Lesson: `void *` removes type safety. The compiler cannot catch this, so you must.

---

# Pointer to pointer

A pointer can point to another pointer. The type is written with two stars, like `char **`. You need this when a function must change *where a pointer points* (this is why the `char *` swap above passed `&a`).

```c
void make_buffer(char **out) {
    *out = malloc(16);           // changes the caller's pointer
}

char *buf = NULL;
make_buffer(&buf);
if (buf != NULL) {
    strcpy(buf, "hello");
    free(buf);
}
```

This pattern is called an **out-parameter**.

---

# void pointer

`void *` is a generic pointer that can hold the address of any object. It carries no type information.

Rules:
- Any object pointer converts to `void *` and back without a cast in C.
- You **cannot dereference** a `void *`, because the compiler does not know how many bytes to read. Convert it to a real type first.
- Standard C does not allow pointer arithmetic on `void *`.
- You must convert back to the *correct* original type. Nothing checks this.

Because the type is lost, you usually pass a tag that says what it is:

```c
typedef enum {
    TYPE_INT,
    TYPE_STRING,
} type_t;

void print(void *a, type_t type) {
    switch (type) {
    case TYPE_INT: {
        int *x = (int *)a;
        printf("%d\n", *x);
        break;
    }
    case TYPE_STRING: {
        char *s = (char *)a;
        printf("%s\n", s);
        break;
    }
    default:
        printf("unknown\n");
    }
}

int x = 10;
char *y = "hello";

print(&x, TYPE_INT);      // 10
print(y, TYPE_STRING);    // hello
```

---

# Function pointers

Functions live in the text segment, so they have addresses too, and you can store and pass them.

General shape:

```text
return_type (*name)(parameter_types);
```

The parentheses matter:

```c
int *f(int);      // a function returning int *
int (*f)(int);    // a pointer to a function taking int and returning int
```

Basic use:

```c
int minus(int a, int b) { return a - b; }
int mult(int a, int b)  { return a * b; }

int main(void) {
    int (*exec)(int, int);

    exec = minus;
    printf("%d\n", exec(20, 10));   // 10

    exec = mult;
    printf("%d\n", exec(10, 5));    // 50
}
```

A `typedef` makes the type readable:

```c
typedef void *(*func)(void *args);   // pointer to a function taking void*, returning void*

void *add(void *args) {
    int *A = (int *)args;
    int *res = malloc(sizeof(int));
    if (res == NULL) return NULL;
    *res = A[0] + A[1];
    return res;
}

int main(void) {
    int args[2] = {10, 20};
    func f = add;
    int *sum = f(args);

    if (sum != NULL) {
        printf("%d\n", *sum);   // 30
        free(sum);
    }
}
```

## Callbacks

A **callback** is a function you pass to another function so it can call it later. The standard library's `qsort` works this way, and it combines `void *` with function pointers:

```c
int cmp_int(const void *a, const void *b) {
    int x = *(const int *)a;
    int y = *(const int *)b;
    return (x > y) - (x < y);   // negative, zero, or positive
}

int arr[] = {5, 2, 9, 1};
qsort(arr, 4, sizeof(arr[0]), cmp_int);
```

## Dispatch table (route handlers)

A table that maps a key to a function is common in servers.

```c
#include <stdio.h>
#include <string.h>

typedef int (*handler_func)(int client_id, const char *body);

typedef struct {
    const char *method;
    const char *path;
    handler_func handler;
} route_t;

#define MAX_ROUTES 16
static route_t routes[MAX_ROUTES];
static int route_count = 0;

int register_route(const char *method, const char *path, handler_func fn) {
    if (route_count >= MAX_ROUTES) return -1;
    routes[route_count++] = (route_t){method, path, fn};
    return 0;
}

int dispatch(const char *method, const char *path, int client_id, const char *body) {
    for (int i = 0; i < route_count; i++) {
        if (strcmp(routes[i].method, method) == 0 &&
            strcmp(routes[i].path, path) == 0) {
            return routes[i].handler(client_id, body);
        }
    }
    printf("404 not found\n");
    return -1;
}

int handle_home(int client_id, const char *body)  { printf("Hello\n"); return 0; }
int handle_login(int client_id, const char *body) { printf("LOGIN\n"); return 0; }

int main(void) {
    register_route("GET", "/", handle_home);
    register_route("GET", "/login", handle_login);

    dispatch("GET", "/login", 1, "");   // LOGIN
    dispatch("GET", "/nope", 1, "");    // 404 not found
}
```

The route table must be stored somewhere that outlives `register_route`. A local `Route x` would vanish when the function returns. Storing string literals is safe because they live in read-only memory for the whole run.

---

# Dangling pointers and escape to heap

## The problem

```c
int *test(void) {
    int x = 10;
    return &x;      // address of a local variable
}

int *temp = test();
printf("%d\n", *temp);   // undefined behavior
```

`x` lives in `test`'s stack frame. When the function returns, that frame is gone and the memory is reused, so `temp` is a **dangling pointer** (it points to memory that is no longer valid). This usually compiles with only a warning (`gcc -Wall` says "function returns address of local variable"), so it is a bug, not a compile error.

Ways to fix it:
1. Let the caller own the variable and pass its address in (out-parameter).
2. Allocate on the heap with `malloc`; the caller frees it.
3. Use a `static` variable. This works, but it is shared and not thread-safe.

## Escape to heap

In Go, the compiler decides automatically when a local variable must "escape" to the heap, because a pointer to it outlives the function (escape analysis). In C, you make that decision yourself with `malloc`.

The simplest rule: the function returns a `malloc`'d pointer and the caller calls `free`.

```c
int *new_int(int value) {
    int *x = malloc(sizeof *x);
    if (x == NULL) return NULL;
    *x = value;
    return x;
}

int *p = new_int(60);
if (p != NULL) {
    printf("%d\n", *p);   // 60
    free(p);
}
```

A variation is to keep a registry of everything allocated so it can all be freed in one place. This is a simplified version of an arena allocator.

```c
#define MAX_REFS 10
static int *heap_refs[MAX_REFS];
static int ref_count = 0;

int *new_int_tracked(int value) {
    if (ref_count >= MAX_REFS) return NULL;
    int *x = malloc(sizeof *x);
    if (x == NULL) return NULL;
    *x = value;
    heap_refs[ref_count++] = x;
    return x;
}

void free_all(void) {
    for (int i = 0; i < ref_count; i++) {
        free(heap_refs[i]);
        heap_refs[i] = NULL;
    }
    ref_count = 0;
}

int main(void) {
    int *tmp = new_int_tracked(60);
    if (tmp != NULL) printf("Got (escape to heap) = %d\n", *tmp);
    free_all();
}
```

Note that the array size must be a constant expression. `const int size = 10;` is not one in C, so use `#define` or an `enum`.

---

# Common mistakes

| Bug | What happens | How to avoid |
| --- | --- | --- |
| Uninitialized pointer | Points to a random address | Initialize to a value or `NULL` |
| NULL dereference | Segmentation fault | Check before use |
| Dangling pointer | Points to freed or expired memory | Never return `&local` |
| Use after free | Reads or writes freed memory | Set the pointer to `NULL` after `free` |
| Double free | Corrupts the heap, often crashes | One owner, free once |
| Memory leak | Memory never returned | Every `malloc` needs a `free` |
| Out of bounds | Overwrites neighboring memory | Track lengths, check indexes |

Tools that catch these:

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined main.c -o a.out
```

AddressSanitizer (`-fsanitize=address`) reports use-after-free, double free, and out-of-bounds access at runtime. It works on Linux and macOS. `valgrind` is another option on Linux.

---

# Exercises

## Exercise 01: zero value

You have an object that holds either an `int`, a `float`, or a `bool`. Write a function that sets it to its **zero value**: `0` for int, `0.0` for float, `false` for bool.

```c
#include <stdbool.h>

typedef enum {
    KIND_INT,
    KIND_FLOAT,
    KIND_BOOL,
} object_kind_t;

typedef struct { int val; }   obj_int_t;
typedef struct { float val; } obj_float_t;
typedef struct { bool val; }  obj_bool_t;

void set_zero(void *obj, object_kind_t kind);   // you implement this
```

Hints:
- Use `kind` to decide which struct type `obj` really is, then cast.
- Test with a value that starts non-zero and print it after `set_zero`.
- Bonus: replace the `switch` with an array of function pointers indexed by `kind`.

## Exercise 02: reverse in place

Write `void reverse(int *arr, size_t n)` using two pointers, one at the start and one at the end, moving toward each other.

## Exercise 03: string length

Implement `size_t my_strlen(const char *s)` using only pointer arithmetic, with no indexing.

## Exercise 04: sort by function pointer

Write a `qsort` comparator that sorts an array of `node_t` by `age`, then by `name`.

---

# Summary

- A pointer stores an address. `&` gets an address, `*` follows one.
- C is pass by value, so pass a pointer when a function must change the original or avoid a big copy.
- Arrays decay to pointers and lose their length, so pass the length separately.
- `void *` is generic but not type-safe. You must cast back to the right type yourself.
- Function pointers let you pass behavior around: callbacks and dispatch tables.
- Never return the address of a local variable. Use the heap and free it.