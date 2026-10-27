---
title: "Custom Header in C Programming"
date: 2026-10-27
tags:
  - c
  - header
  - compilation
---

# Overview

Sometimes you want to split code into several files, or write your own small library. In C this is done with two kinds of files:

- A **source file** (`.c`) holds the implementation, the actual function bodies.
- A **header file** (`.h`) holds the declarations (signatures) that other files need in order to call those functions.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [Example](#example) | `mymath.c`, `mymath.h`, `main.c` |
| 2 | [How to run](#how-to-run) | Compile and link |
| 3 | [Declaration vs definition](#declaration-vs-definition) | What goes where |
| 4 | [Include guards](#include-guards) | Why `#ifndef` is needed |
| 5 | [Angle brackets vs quotes](#angle-brackets-vs-quotes) | `<>` vs `""` |
| 6 | [Separate compilation](#separate-compilation) | `.o` files, linking, Makefile |
| 7 | [Sharing variables](#sharing-variables) | `extern`, `static` |
| 8 | [Common errors](#common-errors) | Messages and their causes |
| 9 | [Exercises](#exercises) | Practice |

---

# Example

Don't name your files `math.c` and `math.h`. `math.h` is a standard library header, so a file with that name can clash with it (more in [Common errors](#common-errors)). Use a different name such as `mymath`.

Project layout:

```text
project/
├── main.c
├── mymath.c
└── mymath.h
```

## mymath.h (the interface)

Only declarations go here. This is what other files see.

```c
#ifndef MYMATH_H
#define MYMATH_H

int Add(int a, int b);

#endif
```

## mymath.c (the implementation)

```c
#include "mymath.h"

int Add(int a, int b) {
    return a + b;
}
```

Including its own header in the `.c` file is good practice: the compiler then checks that the declaration and the definition match.

## main.c (the user)

```c
#include <stdio.h>
#include "mymath.h"

int main(void) {
    printf("%d\n", Add(2, 3));   // 5
    return 0;
}
```

---

# How to run

```bash
gcc -Wall -Wextra mymath.c main.c -o app && ./app
```

- `mymath.c main.c` lists every `.c` file. Header files are not listed, because `#include` pulls them in.
- `-o app` names the output. Without it, the program is called `a.out`.
- `-Wall -Wextra` turns on useful warnings.

---

# Declaration vs definition

| | Declaration | Definition |
| --- | --- | --- |
| Meaning | "this function exists, here is its signature" | "here is the actual code" |
| Example | `int Add(int a, int b);` | `int Add(int a, int b) { return a + b; }` |
| Goes in | `.h` file | `.c` file |
| How many times | as many as you like | exactly once in the whole program |

The compiler needs a declaration before you call a function, so it knows the argument and return types. The linker later finds the one definition.

Rule of thumb: **don't put function definitions or variable definitions in a header.** If two `.c` files include that header, each gets its own copy and the linker reports a "multiple definition" error.

Headers normally contain:
- function declarations
- `struct`, `enum` and `typedef` definitions
- `#define` constants and macros
- `extern` variable declarations

---

# Include guards

`#include` is plain copy-paste done by the **preprocessor** (the step that runs before the compiler). If the same header gets included twice in one file, the compiler sees its contents twice, and repeated `struct` or `typedef` definitions cause errors.

This happens easily:

```text
main.c includes a.h and b.h
b.h also includes a.h   ->   a.h is pasted twice into main.c
```

An **include guard** stops it:

```c
#ifndef MYMATH_H     // if MYMATH_H is not defined yet...
#define MYMATH_H     // ...define it, and keep the contents below

/* declarations */

#endif               // end of the guarded part
```

The first time, the macro is undefined, so the contents are kept. The second time it is defined, so the whole body is skipped.

- Name the macro after the file, in uppercase (`MYMATH_H`). It must be unique across the project.
- A common alternative is `#pragma once` on the first line. It is shorter and works in GCC, Clang and MSVC, but it is not part of the C standard.

---

# Angle brackets vs quotes

```c
#include <stdio.h>      // system headers: searched in the compiler's standard directories
#include "mymath.h"     // your headers: searched in the current file's folder first, then the standard ones
```

Use `<>` for the standard library and installed libraries, and `""` for your own files. To tell the compiler about extra header folders, use `-I`:

```bash
gcc -Iinclude src/main.c src/mymath.c -o app
```

A common layout is headers in `include/` and sources in `src/`.

Include order in a `.c` file (a common convention): its own header first, then standard headers, then other project headers. Putting your own header first proves it doesn't secretly depend on something included earlier.

---

# Separate compilation

Building a C program has two main stages:

1. **Compile**: each `.c` file becomes an **object file** (`.o`) containing machine code, with calls to other files left as unresolved references.
2. **Link**: the **linker** joins all `.o` files (and libraries) into one executable, connecting each call to its definition.

```bash
gcc -c mymath.c        # produces mymath.o
gcc -c main.c          # produces main.o
gcc mymath.o main.o -o app   # link
```

Why bother? When you change one file, you only recompile that file and re-link. For big projects this saves a lot of time.

## Makefile

`make` automates this. It rebuilds only what changed.

```make
CC = gcc
CFLAGS = -Wall -Wextra -g

app: main.o mymath.o
	$(CC) main.o mymath.o -o app

main.o: main.c mymath.h
	$(CC) $(CFLAGS) -c main.c

mymath.o: mymath.c mymath.h
	$(CC) $(CFLAGS) -c mymath.c

clean:
	rm -f *.o app
```

Notes:
- Command lines must start with a **tab**, not spaces.
- Each rule is `target: dependencies`. Listing `mymath.h` means a header change triggers a rebuild.

## Static library (optional)

To package your code as a reusable library:

```bash
gcc -c mymath.c
ar rcs libmymath.a mymath.o
gcc main.c -L. -lmymath -o app    # -L. = look in the current folder, -lmymath = link libmymath.a
```

---

# Sharing variables

Functions are visible to other files by default. For variables and for hiding things, use these keywords:

```c
/* config.h */
extern int counter;       // declaration only: "this exists somewhere else"
```

```c
/* config.c */
int counter = 0;          // the one definition
```

`static` at file level makes a name **private to that file**:

```c
/* mymath.c */
static int helper(int x) {   // not visible to other files, can't clash with their names
    return x * 2;
}
```

Good habit: anything not declared in the header should be `static`. The header is then your public interface, and everything else is internal.

C has no namespaces, so two libraries that both define `Add` will clash. Libraries avoid this with a prefix, like `mymath_add`.

---

# Common errors

| Error message | Cause | Fix |
| --- | --- | --- |
| `undefined reference to 'Add'` | declared in a header, but the `.c` file with the definition was not compiled or linked | add `mymath.c` to the command |
| `multiple definition of 'Add'` | a function or variable is defined in a header included by several files, or the same `.c` was passed twice | move the definition to a `.c` file |
| `implicit declaration of function 'Add'` | called before any declaration, usually a missing `#include` | include the header |
| `redefinition of 'struct ...'` | missing include guard | add `#ifndef` guard |
| `fatal error: mymath.h: No such file` | wrong path | check the folder or add `-I` |

## Why not name it math.h?

If your project folder has its own `math.h`, then `#include "math.h"` picks yours. But any code that does `#include <math.h>` may also find yours if the folder is on the include path (`-I.`), and then standard functions like `sqrt` disappear. It also confuses readers. Pick a unique name.

---

# Exercises

## Exercise 01: extend mymath

Add `Sub`, `Mul` and `Div` to your library. Think about what `Div` should do when the divisor is zero: return an error code, or take an output pointer like `int Div(int a, int b, int *result)`?

## Exercise 02: string library

Move the string builder from the strings note into `mystring.c` and `mystring.h`. Keep `string_t` in the header and the function bodies in the `.c` file.

## Exercise 03: hide the internals

Make the struct **opaque** (its fields are hidden from users). The header declares only `typedef struct Stack Stack;` plus functions like `Stack *stack_new(void)`, `stack_push`, and `stack_free`. The `struct Stack { ... }` definition lives only in the `.c` file. Users can only use the pointer, which protects your internals.

## Exercise 04: Makefile

Write a Makefile for a project with `src/`, `include/`, and a `build/` folder for `.o` files.

---

# Summary

- `.h` is the interface (declarations), `.c` is the implementation (definitions).
- Compile every `.c` file, and the linker joins them.
- Always use include guards in headers.
- Never define functions or variables in a header.
- Use `static` for anything private to a file.
- Don't reuse standard header names like `math.h`.