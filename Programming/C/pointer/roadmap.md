---
title: Pointers in C Roadmap
date: 2026-10-25
tags:
  - c
  - pointers
  - roadmap
summary: A ten-stage study roadmap for learning pointers in C, from the memory model to a storage engine
---

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