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
- Capstone Project Ideas
- Compilation Phases
- Memory Management
- Multithreading

## Dev Containers - VS Code Extension

## Resources

- [Stack Frame and GDB](https://www.youtube.com/watch?v=7dMTCdFM2ss)
- [Garbage Collector](https://journal.stuffwithstuff.com/2013/12/08/babys-first-garbage-collector/)
- [Memory Management in C](https://www.youtube.com/watch?v=rJrd2QMVbGM&t=30s)
- [Socket to Me: Where Do Sockets Live in Go? - Gabbi Fisher](https://www.youtube.com/watch?v=pGR3r0UhoS8)
- [Asynchronous Networking Patterns - Filippo Valsorda](https://www.youtube.com/watch?v=afSiVelXDTQ)

## C Notes

- `memmove`: copies bytes from one memory region to another, and it is safe when the two regions overlap. **Always check how many bytes you are moving**, and make sure the destination has enough room. Both the source and the destination must be pointers.

```c
#include <string.h>

memmove(destination, source, size_of_memory);
memmove(&arr[idx + 1], &arr[idx], sizeof(int) * 3);   // shift 3 ints one position to the right
```

- `bool`: before C23, you must include the `<stdbool.h>` header. From C23 onward, `bool` is a built-in keyword.

### Header Files for Socket Programming

- `<unistd.h>`: `close()`, `read()`, `write()`
- `<sys/socket.h>`: `socket()`, `bind()`, `listen()`, `accept()`
- `<netinet/in.h>`: `sockaddr_in`, `htons()`, `ntohs()`
- `<arpa/inet.h>`: `inet_pton()`, `inet_ntop()`
- `<errno.h>`: the `errno` variable for error codes


## Alias

```c
typedef uint8_t int8;
typedef uint8_t byte;
typedef uint16_t int16;
typedef uint32_t int32;
typedef uint64_t int64;
```