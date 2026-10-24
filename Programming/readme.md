---
title: Programming Languages
date: 2026-10-24
---

Now that we know the concepts, how do we implement them? With programming, of course.

---

## Overview

I will cover only three programming languages (C, C++ and Golang), because those are the only ones I know. I will update this if I ever switch to another.

For Golang, I might use a debugger and a profiler to understand what is happening. For C and C++, I will cover only the syntax and concepts, and maybe some OOP for C++.

---

## Why C?

C influenced many later languages and it is useful for understanding concepts clearly. Also, eBPF programs are usually written in a restricted subset of C.

## Resources

- [Stack frame and gdb](https://www.youtube.com/watch?v=7dMTCdFM2ss)
- [Garbage collector](https://journal.stuffwithstuff.com/2013/12/08/babys-first-garbage-collector/)
- [Memory management in C](https://www.youtube.com/watch?v=rJrd2QMVbGM&t=30s)
- [Socket to Me: Where do Sockets Live in Go? - Gabbi Fisher](https://www.youtube.com/watch?v=pGR3r0UhoS8)
- [ Asynchronous Networking Patterns - Filippo Valsorda](https://www.youtube.com/watch?v=afSiVelXDTQ)

## C Notes

- `memmove`: copies bytes from one memory region to another, and it is safe when the two regions overlap. **Always check how many bytes you are moving.** Both the source and the destination must be pointers.

```c
#include <string.h>

memmove(destination, source, size_of_memory);
memmove(&arr[idx + 1], &arr[idx], sizeof(int) * 3);   // shift 3 ints one position to the right
```

- `bool`: you must include the `<stdbool.h>` header.

### Header files for socket programming

- `<unistd.h>`: `close()`, `read()`, `write()`
- `<sys/socket.h>`: `socket()`, `bind()`, `listen()`, `accept()`
- `<netinet/in.h>`: `sockaddr_in`, `htons()`, `ntohs()`
- `<arpa/inet.h>`: `inet_pton()`, `inet_ntop()`
- `<errno.h>`: the `errno` variable for error codes