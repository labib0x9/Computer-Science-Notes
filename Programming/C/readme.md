## Overview

You might wonder: why did I choose C instead of a simpler language like Python?

I have used both, but I started coding in C. My reasons:

1. I wanted something less abstracted.
2. I started learning Go, and when I got into concurrency, I mixed up threads and goroutines.
3. I was afraid of pointers and memory management, and learning C was the best way to push myself out of my comfort zone.
4. eBPF programs are usually written in a restricted subset of C.

I did not cover every basic topic. Fundamentals such as variable declarations, loops and conditions are missing.
Also these C codes are tested on Mac and Linux using the docker image.

- Functions
- Structs, Enums, Unions
- Pointers
- Strings
- Custom Headers
- Compilation Phases
- Memory Management
- Multithreading
- Sockets
- Capstone Project

## Dev Containers - VS Code Extension

## C Notes

- `memmove`: copies bytes from one memory region to another, and it is safe when the two regions overlap. **Always check how many bytes you are moving**, and make sure the destination has enough room. Both the source and the destination must be pointers.

- `bool`: before C23, you must include the `<stdbool.h>` header. From C23 onward, `bool` is a built-in keyword.


## Alias

```c
typedef uint8_t int8;
typedef uint8_t byte;
typedef uint16_t int16;
typedef uint32_t int32;
typedef uint64_t int64;
```