---
title: C Compilation Stages with GCC
date: 2026-10-25
tags:
  - c
  - gcc
  - compilation
summary: The four stages GCC runs to turn C source into an executable, and the command for each
---

## Keep all intermediate files

```bash
gcc -save-temps test.c
```

This builds the program as usual, but keeps the files from every stage (`test.i`, `test.s`, `test.o`) in the current folder instead of deleting them.

## The stages

1. **Source code:** `test.c`
2. **Preprocessing:** `gcc -E test.c -o test.i`
   Handles `#include`, `#define` and other `#` lines. The output is still C code, with the headers pasted in.
3. **Compilation:** `gcc -S test.i -o test.s`
   Turns C code into assembly code (human-readable machine instructions).
4. **Assembly:** `gcc -c test.s -o test.o`
   Turns assembly into an object file (`test.o`): machine code that is not yet runnable, because calls to other files and libraries are not connected.
5. **Linking:** `gcc test.o -o test`
   Joins the object files and libraries (such as the C standard library) into the final executable `test`. With several files, list all the `.o` files: `gcc a.o b.o -o program`.