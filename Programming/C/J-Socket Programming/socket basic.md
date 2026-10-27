---
title: "Socket Programming in C"
date: 2026-10-25
tags:
  - c
  - networking
  - sockets
summary: How sockets work in C, from the arguments of socket() to a working TCP server and client
---

# Overview

A **socket** is an endpoint for communication. You create one, connect it to something, and then read and write bytes, much like a file. It can connect two programs on one machine, or two machines across the internet.

In C, a socket is a **file descriptor** (a small integer that the OS uses to identify an open file, pipe, or socket). Because of that, `read`, `write` and `close` work on sockets too, and `send` and `recv` are versions with extra options.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [Pipes](#pipes) | A simpler one-way channel, as a warm-up |
| 2 | [Creating a socket](#creating-a-socket) | `socket()` arguments |
| 3 | [Client and server flow](#client-and-server-flow) | Which call comes when |
| 4 | [A minimal server](#a-minimal-server) | Smallest working version |
| 5 | [A client](#a-client) | `connect()` |
| 6 | [Byte order and addresses](#byte-order-and-addresses) | `htons`, `sockaddr`, `inet_pton` |
| 7 | [A better server](#a-better-server) | Errors, timeouts, loops |
| 8 | [Socket options](#socket-options) | `SO_REUSEADDR`, `SO_RCVTIMEO` |
| 9 | [Partial reads and writes](#partial-reads-and-writes) | TCP has no message boundaries |
| 10 | [Handling many clients](#handling-many-clients) | Why one-at-a-time is limiting |
| 11 | [Testing and debugging](#testing-and-debugging) | `nc`, `lsof`, `ss` |
| 12 | [Common errors](#common-errors) | Messages and causes |
| 13 | [Platform notes](#platform-notes) | Linux vs macOS |
| 14 | [Exercises](#exercises) | Practice |

---

# Pipes

A pipe is the simplest way for two parts of a program to talk. It is **one-directional** (data goes in one end and comes out the other), and it works between a parent and child process (after `fork()`), or between threads of one process.

```text
write(fd[1])  --->  [ kernel buffer ]  --->  read(fd[0])
```

```c
#include <stdio.h>
#include <string.h>
#include <unistd.h>

int main(void) {
    int pipe_fd[2];
    if (pipe(pipe_fd) == -1) {
        perror("pipe");
        return 1;
    }

    // fd[0] is the read end, fd[1] is the write end.
    printf("Read end,  fd[0] = %d\n", pipe_fd[0]);
    printf("Write end, fd[1] = %d\n", pipe_fd[1]);

    const char *msg = "Hello Labib\n";
    ssize_t n = write(pipe_fd[1], msg, strlen(msg));
    printf("%zd bytes written to fd[1]\n", n);

    char buf[16];
    n = read(pipe_fd[0], buf, sizeof buf - 1);   // leave room for '\0'
    if (n < 0) {
        perror("read");
        return 1;
    }
    buf[n] = '\0';                               // read() does not add a terminator
    printf("%zd bytes read from fd[0]\n", n);
    printf("Data = %s", buf);

    close(pipe_fd[0]);
    close(pipe_fd[1]);
    return 0;
}
```

Points to remember:
- `read` returns how many bytes it got and does **not** add `'\0'`, so printing the buffer with `%s` without terminating it reads garbage.
- Reading from an empty pipe **blocks** (waits). Writing to a full pipe also blocks (the buffer is typically 64 KB on Linux).
- Closing every write end makes `read` return `0` (end of file). That is how the reader knows the writer is done.
- Writing when all read ends are closed raises `SIGPIPE` and kills the process by default. Sockets have the same problem (see [Common errors](#common-errors)).

## Pipe with fork

Each process should close the end it doesn't use:

```c
#include <stdio.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main(void) {
    int fd[2];
    if (pipe(fd) == -1) { perror("pipe"); return 1; }

    pid_t pid = fork();
    if (pid == 0) {                          // child: reader
        close(fd[1]);                        // not writing
        char buf[64];
        ssize_t n = read(fd[0], buf, sizeof buf - 1);
        if (n > 0) { buf[n] = '\0'; printf("child got: %s\n", buf); }
        close(fd[0]);
        _exit(0);
    }

    close(fd[0]);                            // parent: writer, not reading
    const char *msg = "hello from parent";
    write(fd[1], msg, strlen(msg));
    close(fd[1]);
    wait(NULL);                              // wait for the child to finish
    return 0;
}
```

A pipe can't go both ways and can't cross machines. Sockets fix both limits.

---

# Creating a socket

```c
#include <sys/socket.h>

int fd = socket(domain, type, protocol);
```

- `domain`: the address family, which decides what kind of addresses the socket uses. Examples: `AF_INET` (IPv4), `AF_INET6` (IPv6), `AF_UNIX` (local communication on one machine).
- `type`: the kind of communication. `SOCK_STREAM` gives a reliable byte stream (TCP for IPv4/IPv6). `SOCK_DGRAM` gives separate messages without a connection (UDP).
- `protocol`: the specific protocol. Pass `0` to let the system choose the default for the given domain and type, which is almost always what you want.

On success, `socket()` returns a file descriptor. On failure it returns `-1` and sets `errno` (a global error code that `perror` turns into a message).

```c
int fd = socket(AF_INET, SOCK_STREAM, 0);   // an IPv4 TCP socket
if (fd < 0) {
    perror("socket");
    return 1;
}
```

The first descriptor you get is usually `3`, because `0`, `1` and `2` are stdin, stdout and stderr. Always `close(fd)` when done, or you leak descriptors.

---

# Client and server flow

```text
Server                              Client
------                              ------
socket()                            socket()
bind()      attach to an address
listen()    start accepting
accept()  <------ handshake ------  connect()
recv()    <------- data ----------  send()
send()    -------- data -------->   recv()
close()                             close()
```

| Call | Who | Does |
| --- | --- | --- |
| `socket` | both | create the endpoint |
| `bind` | server | attach it to an IP address and port |
| `listen` | server | mark it as accepting connections |
| `accept` | server | wait for a client, return a **new** socket for that client |
| `connect` | client | connect to a server's address |
| `send` / `recv` | both | transfer data |
| `close` | both | release the descriptor |

An important detail: the server has **two kinds** of socket. The *listening socket* (`fd`) only waits for new connections. Each `accept` returns a separate *connection socket* (`c_fd`) for talking to one client.

---

# A minimal server

This accepts one client, reads once, replies, and exits.

```c
#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main(void) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) { perror("socket"); return 1; }

    struct sockaddr_in addr;
    memset(&addr, 0, sizeof addr);
    addr.sin_family = AF_INET;                  // IPv4
    addr.sin_port = htons(8080);                // port, converted to network byte order
    addr.sin_addr.s_addr = htonl(INADDR_ANY);   // all local interfaces (0.0.0.0)

    if (bind(fd, (struct sockaddr *)&addr, sizeof addr) < 0) {
        perror("bind");
        close(fd);
        return 1;
    }
    if (listen(fd, 10) < 0) {
        perror("listen");
        close(fd);
        return 1;
    }

    int c_fd = accept(fd, NULL, NULL);          // blocks until a client connects
    if (c_fd < 0) {
        perror("accept");
        close(fd);
        return 1;
    }
    printf("client = %d\n", c_fd);

    char request[1024];
    ssize_t n = recv(c_fd, request, sizeof request - 1, 0);
    if (n > 0) {
        request[n] = '\0';
        printf("received: %s\n", request);

        const char reply[] = "HELLO CLIENT";
        send(c_fd, reply, strlen(reply), 0);    // strlen: don't send the '\0'
    }

    close(c_fd);
    close(fd);
    return 0;
}
```

Notes:
- `INADDR_ANY` binds to every network interface, so other machines on your network can reach the server. To keep it local only, bind to `127.0.0.1` (see the better server below).
- `listen(fd, 10)`: the second argument is the **backlog**, the number of connections the kernel will hold in a queue while they wait for you to call `accept`. Without a queue, a client arriving while you are busy would be refused. The OS may cap the value (on Linux it is limited by `/proc/sys/net/core/somaxconn`).
- Ports below 1024 usually need root. Pick a high port like `8080`.

---

# A client

```c
#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main(void) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) { perror("socket"); return 1; }

    struct sockaddr_in server;
    memset(&server, 0, sizeof server);
    server.sin_family = AF_INET;
    server.sin_port = htons(8080);
    if (inet_pton(AF_INET, "127.0.0.1", &server.sin_addr) != 1) {
        fprintf(stderr, "bad address\n");
        close(fd);
        return 1;
    }

    if (connect(fd, (struct sockaddr *)&server, sizeof server) < 0) {
        perror("connect");                      // "Connection refused" if nothing listens
        close(fd);
        return 1;
    }

    const char *msg = "hello server";
    send(fd, msg, strlen(msg), 0);

    char buf[1024];
    ssize_t n = recv(fd, buf, sizeof buf - 1, 0);
    if (n > 0) {
        buf[n] = '\0';
        printf("server said: %s\n", buf);
    }

    close(fd);
    return 0;
}
```

The client doesn't call `bind`. The OS picks a free local port (an **ephemeral port**) automatically.

---

# Byte order and addresses

## Byte order

Computers store multi-byte numbers in different orders. Most CPUs today (x86, ARM) are **little-endian** (least significant byte first), but the network always uses **big-endian**, called *network byte order*. Convert whenever you put a number into an address struct:

| Function | Converts |
| --- | --- |
| `htons(x)` | host to network, 16-bit (ports) |
| `htonl(x)` | host to network, 32-bit (IPv4 addresses) |
| `ntohs(x)` | network to host, 16-bit |
| `ntohl(x)` | network to host, 32-bit |

Easy way to remember: **h**ost, **n**etwork, **s**hort, **l**ong. Use `htons` when you *write* a port, and `ntohs` when you *read* one (for example, to print it).

## Address structs

The socket functions take a generic `struct sockaddr *`, but you fill in a more specific struct and cast it:

| Struct | Used for |
| --- | --- |
| `struct sockaddr` | generic type in function signatures |
| `struct sockaddr_in` | IPv4: family, port, address |
| `struct sockaddr_in6` | IPv6 |
| `struct sockaddr_un` | Unix domain sockets (a file path) |
| `struct sockaddr_storage` | big enough for **any** of the above |

Use `sockaddr_storage` when you don't know the family in advance, such as the client address from `accept` on a server that may handle both IPv4 and IPv6. Check its `ss_family` field, then cast.

For `accept`, the length argument works both ways: you pass in the size of your buffer, and the kernel writes back the real size.

```c
struct sockaddr_storage addr;
socklen_t len = sizeof addr;        // in: buffer size
int c_fd = accept(fd, (struct sockaddr *)&addr, &len);   // out: actual size
```

## Text addresses

```c
struct in_addr a;
inet_pton(AF_INET, "127.0.0.1", &a);        // text -> binary. Returns 1 on success, 0 if not valid, -1 on error

char ip[INET_ADDRSTRLEN];
inet_ntop(AF_INET, &a, ip, sizeof ip);      // binary -> text
```

The `p` and `n` stand for "presentation" (text) and "network" (binary). Prefer these over the old `inet_addr` and `inet_ntoa`, which are not safe for IPv6 or error handling.

---

# A better server

This version handles many clients (one after another), sets a receive timeout, prints each client's address, and checks every error. It binds to localhost only.

```c
#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>

#define PORT 8080
#define BACKLOG 10

typedef struct {
    int fd;
    struct sockaddr_storage addr;
} client_t;

static void print_client(const client_t *c) {
    if (c->addr.ss_family == AF_INET) {
        const struct sockaddr_in *a = (const struct sockaddr_in *)&c->addr;
        char ip[INET_ADDRSTRLEN];
        inet_ntop(AF_INET, &a->sin_addr, ip, sizeof ip);
        printf("[CLIENT] %s:%d\n", ip, ntohs(a->sin_port));
    }
}

// send() may write fewer bytes than asked, so loop until everything is sent
static int send_all(int fd, const char *buf, size_t len) {
    size_t sent = 0;
    while (sent < len) {
        ssize_t n = send(fd, buf + sent, len - sent, 0);
        if (n < 0) {
            if (errno == EINTR) continue;       // interrupted by a signal, retry
            return -1;
        }
        sent += (size_t)n;
    }
    return 0;
}

static void handle_conn(client_t client) {
    // recv() gives up after 5 seconds of silence instead of blocking forever
    struct timeval tv = {.tv_sec = 5, .tv_usec = 0};
    if (setsockopt(client.fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof tv) < 0) {
        perror("setsockopt SO_RCVTIMEO");
        close(client.fd);
        return;
    }

    print_client(&client);

    char request[1024];
    const char reply[] = "HELLO CLIENT";

    for (;;) {
        ssize_t n = recv(client.fd, request, sizeof request - 1, 0);

        if (n == 0) {                           // client closed the connection normally
            printf("client disconnected\n");
            break;
        }
        if (n < 0) {
            int err = errno;                    // save it: printf could change errno
            if (err == EINTR) continue;
            if (err == EAGAIN || err == EWOULDBLOCK) {
                printf("timed out\n");
            } else {
                fprintf(stderr, "recv: %s\n", strerror(err));
            }
            break;
        }

        request[n] = '\0';
        printf("received %zd bytes\n", n);

        if (send_all(client.fd, reply, strlen(reply)) < 0) {
            perror("send");
            break;
        }
    }

    close(client.fd);
}

int main(void) {
    signal(SIGPIPE, SIG_IGN);       // writing to a closed connection returns an error, not a kill

    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) { perror("socket"); return 1; }

    int opt = 1;
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof opt) < 0) {
        perror("setsockopt SO_REUSEADDR");
        close(fd);
        return 1;
    }

    struct sockaddr_in server_addr;
    memset(&server_addr, 0, sizeof server_addr);
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    if (inet_pton(AF_INET, "127.0.0.1", &server_addr.sin_addr) != 1) {
        fprintf(stderr, "bad address\n");
        close(fd);
        return 1;
    }

    if (bind(fd, (struct sockaddr *)&server_addr, sizeof server_addr) < 0) {
        perror("bind");
        close(fd);
        return 1;
    }
    if (listen(fd, BACKLOG) < 0) {
        perror("listen");
        close(fd);
        return 1;
    }
    printf("listening on 127.0.0.1:%d\n", PORT);

    for (;;) {
        client_t client;
        socklen_t len = sizeof client.addr;
        int c_fd = accept(fd, (struct sockaddr *)&client.addr, &len);
        if (c_fd < 0) {
            perror("accept");
            continue;
        }
        client.fd = c_fd;
        handle_conn(client);
    }

    // unreachable until you add a signal handler for Ctrl+C
    close(fd);
    return 0;
}
```

Details worth understanding:

- `recv` has three kinds of result: **`> 0`** is the number of bytes read, **`0`** means the peer closed the connection cleanly, and **`-1`** means an error (check `errno`). Never use the return value as an array index before checking it.
- Reserve one byte (`sizeof request - 1`) so there is room to add `'\0'`.
- `errno` can be overwritten by any library call, including `printf`. Copy it right after the failing call.
- `SO_RCVTIMEO` makes a timed-out `recv` return `-1` with `EAGAIN` or `EWOULDBLOCK` (both mean "try again later"; they may be the same number).
- Passing a `client_t` by value copies the struct. For bigger structs, pass a pointer.

---

# Socket options

`setsockopt(fd, level, option, &value, sizeof value)` changes how a socket behaves. `level` is `SOL_SOCKET` for generic options.

| Option | What it does |
| --- | --- |
| `SO_REUSEADDR` | allow `bind` to reuse a port that is still in `TIME_WAIT` |
| `SO_RCVTIMEO` | `recv` fails after this much silence (a `struct timeval`) |
| `SO_SNDTIMEO` | same for `send` |
| `SO_KEEPALIVE` | periodically probe idle connections to detect dead peers |

## Why SO_REUSEADDR

When a TCP connection closes, the side that closed first keeps the connection in a `TIME_WAIT` state for about a minute (to catch stray late packets). While that lasts, restarting your server can fail with `bind: Address already in use`. Setting `SO_REUSEADDR` **before `bind`** avoids this. You will want it on almost every server you write.

---

# Partial reads and writes

TCP is a **byte stream**, not a stream of messages. It does not remember where one `send` ended and the next began:

- Two `send` calls of 5 bytes each may arrive as one `recv` of 10 bytes.
- One `send` of 1000 bytes may arrive as several smaller `recv` calls.
- `send` itself may accept fewer bytes than you asked for (that is why `send_all` loops).

So the program must decide where a message ends. The two common methods:

1. **Delimiter**: messages end with a marker, such as `\r\n` in HTTP headers or `\n` in a line-based protocol. Keep reading into a buffer until you find the marker.
2. **Length prefix**: send the size first (for example 4 bytes in network order), then that many bytes. Read exactly that many.

This is the most common beginner mistake in socket code: assuming one `recv` equals one message.

---

# Handling many clients

The server above serves **one client at a time**. While `handle_conn` runs, a second client can connect (the kernel completes the handshake and holds it in the backlog queue), but gets no reply until the first one leaves. That is why the 5-second timeout matters: it stops one idle client from blocking everyone.

Ways to serve many clients at once:

| Approach | Idea | Notes |
| --- | --- | --- |
| `fork()` per client | a child process for each connection | simple and isolated, but heavy; reap children to avoid zombies |
| Thread per client | a thread for each connection | lighter than processes, but needs care with shared data |
| I/O multiplexing | one thread watches many sockets with `poll`, `epoll` or `kqueue` | scales to thousands of clients |

`select` and `poll` are portable. `epoll` is Linux-only and `kqueue` is macOS and BSD. This is the same mechanism that Go's runtime uses under its network code.

---

# Testing and debugging

Run the server in one terminal, then in another:

```bash
nc 127.0.0.1 8080           # type text, press Enter, see the reply
```

Check who is listening:

```bash
lsof -i :8080               # Linux and macOS
ss -ltnp                    # Linux only: listening TCP sockets with process names
```

Compile with warnings and the sanitizer:

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined server.c -o server
```

To see the actual bytes on the wire, use `tcpdump -i lo0 port 8080` on macOS (the loopback interface is `lo` on Linux), or Wireshark.

---

# Common errors

| Message or symptom | Cause | Fix |
| --- | --- | --- |
| `bind: Address already in use` | port taken, or old connection in `TIME_WAIT` | stop the other program; set `SO_REUSEADDR` |
| `bind: Permission denied` | port below 1024 without root | use a higher port |
| `connect: Connection refused` | nothing is listening on that address and port | start the server first |
| Program dies silently on `send` | `SIGPIPE`: peer already closed | `signal(SIGPIPE, SIG_IGN)` and check the return value |
| `recv` returns 0 | peer closed normally | close your side and stop reading |
| Garbage or missing text | buffer not terminated, or partial read | add `'\0'` after `recv`; see [Partial reads and writes](#partial-reads-and-writes) |
| Server unreachable from another machine | bound to `127.0.0.1` | bind to `INADDR_ANY` if you really want that |
| Descriptors run out (`Too many open files`) | forgot `close` | close on every path, including error paths |

---

# Platform notes

- `MSG_NOSIGNAL` exists on Linux but not on macOS. Ignoring the signal with `signal(SIGPIPE, SIG_IGN)` works on both, so the examples use that. (macOS also has the `SO_NOSIGPIPE` socket option.)
- `ss` and `epoll` are Linux-only. On macOS use `lsof` or `netstat -an`, and `kqueue`.
- The loopback interface for packet capture is `lo` on Linux and `lo0` on macOS.
- Windows uses Winsock, which has a different setup (`WSAStartup`) and `closesocket`. These notes cover POSIX systems only.

---

# Exercises

## Exercise 01: echo server

Send back exactly the bytes you received. Test with `nc`, and with a message larger than your buffer to see partial reads.

## Exercise 02: line reader

Write `read_line(int fd, char *buf, size_t size)` that keeps calling `recv` until it finds `'\n'`, handling partial reads and a line that is too long.

## Exercise 03: mini HTTP server

Reply to any request with a valid response, and open it in a browser or with `curl`:

```text
HTTP/1.1 200 OK\r\n
Content-Type: text/plain\r\n
Content-Length: 5\r\n
\r\n
hello
```

`Content-Length` must match the body size exactly. Then add a route table of function pointers, using the ideas from your pointer notes.

## Exercise 04: concurrent server

Serve many clients at once, first with `fork()` (reap children with `waitpid` in a `SIGCHLD` handler), then with threads.

## Exercise 05: chat server

Use `poll()` to watch the listening socket and every client in one thread, and broadcast each message to the others.

## Exercise 06: UDP

Build the same echo using `SOCK_DGRAM` with `sendto` and `recvfrom`. Note what changes: no `listen`, no `accept`, and each datagram is a separate message.

---

### Header Files for Socket Programming

- `<unistd.h>`: `close()`, `read()`, `write()`
- `<sys/socket.h>`: `socket()`, `bind()`, `listen()`, `accept()`
- `<netinet/in.h>`: `sockaddr_in`, `htons()`, `ntohs()`
- `<arpa/inet.h>`: `inet_pton()`, `inet_ntop()`
- `<errno.h>`: the `errno` variable for error codes

# Summary

- A socket is a file descriptor for network communication. `socket(domain, type, protocol)` creates one.
- Server: `socket`, `bind`, `listen`, then `accept` in a loop. Client: `socket`, `connect`.
- `accept` returns a new socket for each client, so the listening socket keeps listening.
- Use `htons` and `ntohs` for ports. Check every return value.
- `recv` returns `> 0` (data), `0` (closed), or `-1` (error).
- TCP is a byte stream, so you must frame messages and handle partial reads and writes.
- Set `SO_REUSEADDR` before `bind`, ignore `SIGPIPE`, and close every descriptor.