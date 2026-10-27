---
title: "C Memory Layout: Text, Data, BSS, Heap, Stack"
date: 2026-10-27
tags:
  - c
  - memory
---

# Overview

Every running program lives in memory, so where exactly do its code and variables go? These are my notes from trying to understand that.

When a C program runs, the operating system gives it a block of **virtual memory** (addresses the program sees, which the OS maps to real RAM). That block is divided into regions called **segments**. Each segment holds a different kind of thing.

> The C standard does not define these segments. They are how Linux and macOS implement it, so names and details vary by platform.

| # | Segment | Holds |
| --- | --- | --- |
| 1 | [TEXT](#text) | Machine code (your functions) |
| 2 | [RODATA](#rodata) | Read-only data (string literals, `const` globals) |
| 3 | [DATA](#data) | Initialized global and static variables |
| 4 | [BSS](#bss) | Uninitialized global and static variables |
| 5 | [HEAP](#heap) | Memory from `malloc` / `calloc` |
| 6 | [STACK](#stack) | Local variables and function call info |

---

# Layout

```text
High address
+------------------+
| Stack            |  grows down
+------------------+
|   (free space)   |
+------------------+
| Heap             |  grows up
+------------------+
| BSS              |
+------------------+
| Data             |
+------------------+
| Rodata           |
+------------------+
| Text             |
+------------------+
Low address
```

TEXT, RODATA, DATA and BSS have a fixed size once the program is loaded. The heap and stack change size while the program runs.

---

# TEXT

The compiled machine instructions. Every function you write lives here. It is normally read-only and executable, so a bug cannot overwrite your code.

# RODATA

Read-only data: string literals and `const` globals. Writing to it crashes the program (segmentation fault).

```c
char *s = "hi";   // "hi" is in rodata; s itself is a variable (data, if global)
s[0] = 'H';       // crash: writing to read-only memory
```

# DATA

Global and `static` variables that have a non-zero initial value. The initial values are stored inside the executable file.

```c
int a = 5;
static int count = 10;
```

# BSS

Global and `static` variables that are uninitialized or initialized to zero. The name is a historical abbreviation, "block started by symbol". The executable only records how big BSS is, not its contents, so it costs no file size. The OS fills it with zeros at startup.

```c
int b;
int c = 0;   // also BSS
```

# HEAP

Memory you request manually with `malloc` and release with `free`. It lasts until you free it, not until the function returns. Forgetting to free causes a memory leak. Using memory after freeing is a bug called use-after-free.

Small allocations usually come from the heap. Large ones (roughly 128 KB and up on Linux) are often served by `mmap` and appear as a separate region.

# STACK

Local variables, function arguments, and return addresses. Each function call pushes a **stack frame** (a block for that call's data), and the frame is removed when the function returns. It is fast and automatic, but small, and very deep recursion causes a stack overflow.

An uninitialized local such as `int c;` holds garbage, not zero.

---

# Example

```c
int a = 5;              // data
int b;                  // bss
const char *s = "hi";   // "hi" is in rodata, s itself is in data
void foo(void) {}       // text

int main(void) {
    int c;                  // stack
    int *p = malloc(4);     // p is on the stack, the 4 bytes are on the heap
    free(p);
    return 0;
}
```

# Printing real addresses

Compare the printed addresses to see that each kind of variable lands in a different range.

```c
#include <stdio.h>
#include <stdlib.h>

int a = 5;               // data
int b;                   // bss
const int r = 7;         // rodata
void foo(void) {}        // text

int main(void) {
    int c;               // stack
    int *p = malloc(4);  // heap

    printf("text   %p\n", (void *)foo);
    printf("rodata %p\n", (void *)&r);
    printf("data   %p\n", (void *)&a);
    printf("bss    %p\n", (void *)&b);
    printf("heap   %p\n", (void *)p);
    printf("stack  %p\n", (void *)&c);

    getchar();           // pause so you can inspect the live process
    free(p);
    return 0;
}
```

Converting a function pointer to `void *` is not strictly standard C, but it works on Linux and macOS.

Addresses change between runs because of **ASLR** (address space layout randomization, a security feature). Expect the order text < rodata < data < bss < heap < stack, though exact gaps vary.

---

# Inspecting memory on Linux

Compile with debug info: `gcc -g -o a.out main.c`

```bash
readelf -S ./a.out     # sections: .text, .data, .bss, .rodata
size ./a.out           # summary of text/data/bss sizes
nm ./a.out             # symbols, with a letter for each segment
```

`nm` letters: `T` text, `D` data, `B` bss, `R` rodata. Lowercase means file-local (`static`).

Live process, using `getchar()` to pause it:

```bash
pmap -x <pid>
cat /proc/<pid>/maps
```

In gdb:

```text
(gdb) break main
(gdb) run
(gdb) info proc mappings      # all address ranges
(gdb) info symbol main        # main in section .text
(gdb) info symbol &a          # a in section .data
(gdb) print &b                # address of a bss variable
(gdb) print $sp               # stack pointer, inside [stack]
(gdb) print p                 # pointer from malloc, inside [heap]
```

Compare each address with the ranges from `info proc mappings` to see which region it falls in.

# Inspecting memory on macOS

`/proc`, `readelf`, `pmap` and `info proc mappings` are Linux-only. macOS uses the Mach-O executable format and lldb in place of gdb.

```bash
size -m ./a.out        # segments and sections with sizes
otool -l ./a.out       # full Mach-O layout (like readelf)
nm ./a.out             # T, D, B letters work similarly
vmmap <pid>            # memory map of a running process
```

In lldb:

```text
(lldb) breakpoint set -n main
(lldb) run
(lldb) image lookup -a <address>   # which section an address is in
(lldb) memory region <address>     # which region, with permissions
(lldb) image dump sections         # list of sections
```

macOS names differ: segments are `__TEXT` and `__DATA`, string literals are in `__TEXT,__cstring`, and heap memory shows up as `MALLOC_*` regions.

---

# Summary

- Text and rodata are read-only. Data and BSS are global state.
- Heap is manual and long-lived. Stack is automatic and tied to function calls.
- Tools: `nm` and `size` for the file, `/proc/<pid>/maps` or `vmmap` for a live process.