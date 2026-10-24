---
title: Socket Programming in C
date: 2026-10-25
tags:
  - c
  - networking
  - sockets
summary: The arguments of socket() in C and what each one means
---

## Creating a socket

```c
#include <sys/socket.h>

int fd = socket(domain, type, protocol);
```

- `domain`: the address family, which decides what kind of addresses the socket uses. Examples: `AF_INET` (IPv4), `AF_INET6` (IPv6), `AF_UNIX` (local communication on one machine).
- `type`: the kind of communication. `SOCK_STREAM` gives a reliable byte stream (TCP for IPv4/IPv6). `SOCK_DGRAM` gives separate messages without a connection (UDP).
- `protocol`: the specific protocol. Pass `0` to let the system choose the default for the given domain and type, which is almost always what you want.

On success, `socket()` returns a file descriptor (a small integer that identifies the open socket, used by later calls like `bind()` and `close()`). On failure it returns `-1` and sets `errno`.

<!-- For the header files used in socket programs, see the C notes in [[programming/README]]. -->