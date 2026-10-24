# C Programming Language - Mother Language
Some of the parts are in `Competitive-Programming` and `DSA` section.
---
- **01-Compiler-Phase :** Compilation phase of a `C program`. Just read the `readme.md` file.
- **02-Function**
- **03-String :** String operations using library functions. Tokenization.
- **04-Pointer :** void pointer, function pointer.
- **05-Struct :** Struct size.
- **06-Enum-And-Union**
- **07-Custom-Header**
- **08-Memory-Management :** Stack frames, object pool, garbage collection.
- **09-Linux-System-Calls :** File syscalls.
- **10-Socket-Programming :** Socket creation, binding. (Need to add roadmap).
- **11-Thread-POSIX :** Multithread programming using pthread api, thread create & join, concurrency, syncronization using mutex, conditional variables and roadmap.
---
# Notes
- memmove : moves memory **(always question how many bytes are moving).** **source and destination must be pointer**.
```c
#include<string.h>
memmove(destination, source, size_of_memory);
memmove(&arr[idx + 1], &arr[idx], sizeof(int) * 3);   // move 3 int
```

- bool : must include `<stdbool.h>` header.
- socket header: `<sys/socket.h>` and `<netinet/in.h>`.

Key Header Files for Socket Programming:
    <unistd.h> - close(), read(), write()
    <sys/socket.h> - socket(), bind(), listen(), accept()
    <netinet/in.h> - sockaddr_in, htons(), ntohs()
    <arpa/inet.h> - inet_pton(), inet_ntop()
    <errno.h> - errno variable for error codes