---
title: "Capstone Project: PIKA, a Toy Key-Value Protocol and Server in C"
date: 2026-10-27
tags:
  - c
  - networking
  - sockets
  - pthreads
  - capstone
summary: A small text protocol over TCP with persistent connections, newline-framed JSON responses, and strict parsing
---

# Overview

PIKA is a small text protocol for a key-value store, served over TCP. A client opens one connection, sends many requests on it, and the server answers each one in order. The goal is to practice the whole C toolkit on one project, in three stages: **protocol, sockets, multithreading**. This note covers the protocol only.

Terms used below:
- **Persistent connection**: one TCP connection that carries many requests, instead of one connection per request.
- **Framing**: deciding where one message ends and the next begins.
- **Pipelining**: sending several requests without waiting for each response.
- **EOF** (end of file): `recv` returning `0`, which means the other side closed the connection.

| # | Section | What it covers |
| --- | --- | --- |
| 1.1 | [Wire format](#wire-format) | Request shape, response shape, methods |
| 1.2 | [Headers](#headers) | Which headers each method accepts, value rules |
| 1.3 | [Examples](#examples) | Full transcripts on one connection |
| 1.4 | [Design decisions](#design-decisions) | What was chosen and why |
| 1.5 | [Connection lifecycle](#connection-lifecycle) | Keep-alive, half-close, when each side closes |
| 1.6 | [Parse strictly](#parse-strictly) | Rules that leave one reading of a message |
| 1.7 | [Data model and semantics](#data-model-and-semantics) | Buckets, keys, value lists, paging |
| 1.8 | [Status values](#status-values) | Response fields, when the connection closes |
| 1.9 | [Limits and timeouts](#limits-and-timeouts) | Sizes and deadlines |
| 2 | [Summary](#summary) | The rules in one list |

---

# 1. The protocol (PIKA/1.0)

## Wire format

A stripped-down, HTTP-like protocol. It is inspired by HTTP, but responses are framed by a newline instead of `Content-Length` or connection close.

**Request:**

```text
METHOD <bucket-name> PIKA/1.0\r\n
name: value\r\n
\r\n
```

A request ends at the blank line (`\r\n\r\n`). There is no body, so arguments go in headers.

**Response:** one line of compact JSON followed by a single `\n`.

```text
{"status":"created","bucket":"sessions"}\n
```

The JSON is never pretty-printed. Compact output plus the string escaping from the encoder guarantees there is no raw newline inside a response, so `\n` always marks its end. A client reads one line per request.

A **bucket** is a named collection of key/value pairs, like a table. The methods:

| Method | Does | Request line | Headers |
| --- | --- | --- | --- |
| `CREATE` | create an empty bucket | `CREATE <bucket> PIKA/1.0` | none |
| `DELETE` | delete a bucket and everything in it | `DELETE <bucket> PIKA/1.0` | none |
| `GET` | read one key, or list keys with `key: *` | `GET <bucket> PIKA/1.0` | `key` (required), `limit`, `after` |
| `SET` | store a value under a key | `SET <bucket> PIKA/1.0` | `key`, `value` (required), `append` |
| `REMOVE` | remove one key | `REMOVE <bucket> PIKA/1.0` | `key` (required) |

`DELETE` and `REMOVE` differ by one letter and both reply `removed`. The response fields tell them apart (`REMOVE` has `key`), but it is an easy mix-up when testing.

## Headers

Header names are case-insensitive (`Key`, `key` and `KEY` are the same). Examples use lowercase. A method accepts only the headers listed for it, and any other header is a `bad_request`.

| Header | Used by | Rule |
| --- | --- | --- |
| `key` | `GET`, `SET`, `REMOVE` | 1 to 256 characters from `A-Z a-z 0-9 . _ -`. `GET` also accepts `*`, which can't clash because `*` is not in the key alphabet |
| `value` | `SET` | 1 to 1024 printable ASCII characters. Leading and trailing spaces are trimmed first, so a value can't start or end with a space |
| `append` | `SET` | `true` or `false`. Default `false` |
| `limit` | `GET` with `key: *` only | digits only, 1 to 1000. Default 10 |
| `after` | `GET` with `key: *` only | a key. Listing starts at the first key **greater** than it (see paging below) |

The bucket name in the request line follows the same alphabet as keys, up to 64 characters. This is why the old `sessions?append=true` form is gone: `?` is not allowed in a bucket name, and options belong in headers.

## Examples

`C>` is a line the client sends and `S>` is a line the server sends. The real bytes end each client line with `\r\n`, which is not shown. A `C>` line with nothing after it is the blank line that ends a request.

One connection, many requests:

```text
C> CREATE sessions PIKA/1.0
C>
S> {"status":"created","bucket":"sessions"}

C> SET sessions PIKA/1.0
C> key: userid_1020
C> value: sessionid_1019
C>
S> {"status":"stored","bucket":"sessions","key":"userid_1020","count":1}

C> SET sessions PIKA/1.0
C> key: userid_1020
C> value: sessionid_1020
C> append: true
C>
S> {"status":"stored","bucket":"sessions","key":"userid_1020","count":2}

C> GET sessions PIKA/1.0
C> key: userid_1020
C>
S> {"status":"found","bucket":"sessions","data":[{"key":"userid_1020","values":["sessionid_1019","sessionid_1020"]}],"more":false}

C> GET sessions PIKA/1.0
C> key: *
C> limit: 10
C>
S> {"status":"found","bucket":"sessions","data":[{"key":"userid_1020","values":["sessionid_1019","sessionid_1020"]},{"key":"userid_1021","values":["sessionid_1021"]}],"more":false}

C> REMOVE sessions PIKA/1.0
C> key: userid_1020
C>
S> {"status":"removed","bucket":"sessions","key":"userid_1020"}

C> DELETE sessions PIKA/1.0
C>
S> {"status":"removed","bucket":"sessions"}
```

(The two `GET` responses above are each one line on the wire. They are long only because every field is shown.)

Pipelining: the client sends two requests in one write, without waiting. Responses come back in the same order.

```text
C> GET sessions PIKA/1.0
C> key: userid_1020
C>
C> REMOVE sessions PIKA/1.0
C> key: userid_1020
C>
S> {"status":"found","bucket":"sessions","data":[{"key":"userid_1020","values":["sessionid_1019"]}],"more":false}
S> {"status":"removed","bucket":"sessions","key":"userid_1020"}
```

Errors that keep the connection open (`not_found`, `exists`):

```text
C> GET sessions PIKA/1.0
C> key: nobody
C>
S> {"status":"not_found","bucket":"sessions","key":"nobody"}

C> CREATE sessions PIKA/1.0
C>
S> {"status":"exists","bucket":"sessions"}
```

Errors that close the connection:

```text
C> SET sessions?append=true PIKA/1.0
C> key: userid_1020
C> value: sessionid_1021
C>
S> {"status":"bad_request","error":"invalid bucket name"}
(server closes the connection)

C> SET sessions PIKA/1.1
C> key: userid_1020
C> value: sessionid_1020
C>
S> {"status":"bad_request","error":"unsupported version"}
(server closes the connection)
```

## Design decisions

| Decision | Choice | Why |
| --- | --- | --- |
| Request body | None. Arguments go in headers | Keeps parsing simple |
| Request framing | Message ends at the blank line. Any bytes after it belong to the **next** request | Needed for persistent connections and pipelining |
| Response format | Compact JSON, one line, ends with `\n` | Self-delimiting, so no length header is needed |
| Response framing | The newline, not EOF | With a persistent connection EOF never comes between responses |
| Connection | Stays open after a response | One connection serves many requests |
| Ordering | One response per request, in request order | Clients match responses by position, so there is no request id |
| Character set | Requests contain only printable ASCII plus `\r\n` | Removes UTF-8 and control-character problems. The encoder then only escapes `"` and `\` |
| Version | Exactly `PIKA/1.0`. Anything else is `bad_request` and the server closes | Never guess what a future version means |
| Extension | New behavior means a new version token | Unknown headers are rejected, so old servers never half-understand new requests |

An alternative for response framing is a `Content-Length` header on every response, like HTTP. It allows pretty-printed JSON but needs you to build the body before the headers. Newline framing was chosen because it is simpler in C.

## Connection lifecycle

The server loop for one connection:

1. Read bytes into a buffer and look for `\r\n\r\n`.
2. Parse that one request and **remember how many bytes it used**. Bytes after it stay in the buffer for the next request.
3. Run it, write one response line, and loop to step 1.
4. Don't read the next request until the current response is fully written. A client that pipelines but never reads responses will eventually block, instead of making the server buffer without limit.

When the connection ends:

| Event | Server does |
| --- | --- |
| Client closes between requests (EOF) | closes its side |
| Client closes in the middle of a request | drops the partial bytes and closes, with no response |
| Client half-closes (`shutdown(SHUT_WR)`) after sending requests | answers every complete request, then closes |
| It sends a `bad_request` or `error` response | closes after sending it, using the lingering close from your sockets notes |
| No request in progress for 30 seconds | closes silently |
| A request is incomplete 10 seconds after its first byte | sends `bad_request` ("request timeout") and closes |
| A response can't be written for 10 seconds | closes |

Rules for clients:
- Read exactly one line per request. Don't wait for EOF.
- Match responses to requests by order.
- If the connection drops before all responses arrive, the unanswered requests may or may not have been applied. `GET` is always safe to retry, and so is a plain `SET`. A `SET` with `append: true` is **not**: retrying it can add the value twice.

## Parse strictly

The more ways one message can be read, the more ways two programs can read it differently. That gap is the root cause of request-smuggling bugs in real HTTP servers. Pick one reading and reject the rest:

- The request line has exactly three parts separated by exactly one space each.
- Line endings are strictly `\r\n`. A bare `\n` is a `bad_request`.
- Method names are case-sensitive: `GET`, not `get`.
- Header names are case-insensitive and may contain only letters and `-`.
- Split a header at the **first** colon only, so values may contain colons. Trim spaces around the value.
- No space before the colon (`key : x` is rejected), and no header line may start with a space.
- An empty header value (`key:`) is a `bad_request`.
- Duplicate headers (two `key:` lines) are a `bad_request`, not "last one wins".
- Unknown headers, and headers not allowed for that method (for example `limit` on `SET`), are a `bad_request`.
- Bucket names, keys, and values follow the character rules in the headers table.
- Any byte other than printable ASCII (`0x20` to `0x7E`) and the `\r\n` pairs is a `bad_request`. This includes the NUL byte (`'\0'`), which C string functions would silently cut at, and tabs.
- Extra bytes between requests (for example a stray blank line) are a `bad_request`, so the parser never has to guess where a request starts.
- When several requests arrive together, all are processed, one at a time, in order. Nothing is dropped.

After a framing or syntax error the server can't know where the next request would start, which is why `bad_request` closes the connection.

## Data model and semantics

- A **bucket** holds keys. A **key** holds an ordered list of one or more string values. The list is always an array in responses, even for one value, so clients handle a single shape.
- `SET` without `append` replaces the list with a single value. With `append: true` it adds the value to the end.
- `CREATE` on an existing bucket replies `exists` and changes nothing. A bucket is **not** created implicitly by `SET`.
- Operations on a missing bucket (`GET`, `SET`, `REMOVE`, `DELETE`) reply `not_found`.
- `GET` of a missing key replies `not_found`. `REMOVE` of a missing key also replies `not_found`.
- `REMOVE` removes the key and all its values. `DELETE` removes the bucket and all its keys.
- `GET` with `key: *` lists keys in ascending byte order of the key, so results are stable and easy to test. An empty bucket gives `found` with `"data":[]`.
- **Paging:** a listing returns at most `limit` entries. `"more":true` means more keys exist after the last one returned. To get the next page, send `after: <last key returned>`.
- Each request is applied completely or not at all. There is no state carried between requests on a connection.

## Status values

Response fields appear in the order shown, so a test can compare whole lines.

| Status | When | Fields | Connection after |
| --- | --- | --- | --- |
| `created` | `CREATE` made a bucket | `bucket` | stays open |
| `exists` | `CREATE` on an existing bucket | `bucket` | stays open |
| `found` | `GET` succeeded | `bucket`, `data`, `more` | stays open |
| `stored` | `SET` succeeded | `bucket`, `key`, `count` (values now stored) | stays open |
| `removed` | `DELETE` or `REMOVE` succeeded | `bucket` (and `key` for `REMOVE`) | stays open |
| `not_found` | bucket or key doesn't exist | `bucket` (and `key` if the request had one) | stays open |
| `rejected` | understood, but a limit blocks it (too many values, store full) | `bucket`, `error` | stays open |
| `bad_request` | malformed, too large, unknown version, timeout | `error` | **closes** |
| `error` | internal failure (out of memory, etc.) | `error` | **closes** |

The `error` text is always a fixed message chosen by the server, never copied from the request. Copying request bytes into a response would put attacker-controlled text into your JSON.

Error responses are responses too, so a bad request still gets a `bad_request` JSON line before the close.

The key is returned in a `"key"` field. Don't use the key itself as a JSON property name: a key called `status` or `bucket` would collide with your own fields.

## Limits and timeouts

| Limit | Value |
| --- | --- |
| Whole request, up to the blank line | 2 KB (the field limits above are chosen so a valid request always fits) |
| Bucket name / key / value | 64 / 256 / 1024 bytes |
| Values per key | 100 (more gives `rejected`) |
| Simultaneous connections | 256 (new connections beyond that are closed immediately) |
| Idle time between requests | 30 seconds |
| Total time to receive one request (first byte to blank line) | 10 seconds |
| Time to write one response | 10 seconds |
| Response size | 1 MB (see below) |

Keep both the idle and the request timeouts. A per-`recv` timeout alone is not enough: a client that sends one byte every 4 seconds never triggers it, yet holds a connection forever (a "slow loris" attack). With persistent connections the idle timeout matters even more, because every open connection occupies a thread or a slot.

**Response size:** a listing of 1000 keys with 100 values of 1 KB each could reach about 100 MB. The server stops adding entries once the response would pass 1 MB and sets `"more":true`, so the client continues with `after`. One key's values (at most about 100 KB) always fit.

---

# Summary

- PIKA is a request/response protocol over one persistent TCP connection. A request ends at a blank line, a response ends at a newline.
- Responses are compact one-line JSON, in the same order as the requests. Pipelining is allowed.
- Parse strictly: one way to read a message, and reject everything else.
- A key holds a list of values. Listings are sorted and paged with `limit`, `after` and `more`.
- `bad_request` and `error` close the connection. Every other status keeps it open.
- Bound everything: request size, value count, response size, and time (idle, request, and write).
- Retrying after a dropped connection is safe except for `SET` with `append: true`.





---
---
---


---
title: "Capstone Project: PIKA, a Toy Key-Value Protocol and server in C"
date: 2026-10-27
tags:
  - c
  - networking
  - sockets
  - pthreads
  - capstone
summary: A small text protocol over TCP in C, built as protocol, then sockets, then three concurrency models, then manual memory vs a garbage collector
---

# Overview

PIKA is a small text protocol for a key-value store, served over TCP. The goal is to practice the whole C toolkit on one project, in three stages: **protocol, sockets, multithreading**. 
<!-- Afterwards you rebuild the connection layer with an event loop (and optionally `io_uring`) and swap manual memory management for a garbage collector, then compare each against the first version. -->

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [The protocol](#1-the-protocol-pika10) | Wire format, strict parsing, limits, statuses |
| 2 | [Protocol Parser](#2-protocol-parser) | Write a Parser for protocol |
<!-- | 2 | [Framing in C](#2-framing-in-c) | Buffers, split reads, JSON encoder |
| 3 | [C-specific pitfalls](#3-c-specific-pitfalls) | Strings, hash table, ownership, clean close |
| 4 | [Build order](#4-build-order) | Steps and "done when" checks |
| 5 | [Concurrency models](#5-concurrency-models) | Threads, event loop, `io_uring` |
| 6 | [Memory management](#6-memory-management-manual-first-then-gc) | Allocator wrapper, Boehm GC, arenas |
| 7 | [Testing and tools](#7-testing-and-tools) | `nc`, fuzzing, sanitizers, load tests |
| 8 | [Extensions](#8-extensions) | Keep-alive, persistence, and more | -->

<!-- Terms used below:
- **Framing**: deciding where one message ends and the next begins.
- **EOF** (end of file): `recv` returning `0`, which means the other side closed the connection.
- **Ownership**: knowing which piece of code is responsible for freeing a piece of memory. -->

---

# 1. The protocol (PIKA/1.0)

A stripped-down, HTTP-like request/response protocol. This protocol was inspired by HTTP/1.0

```text
METHOD <bucket-name> PIKA/1.0\r\n
header: value\r\n
\r\n
```

A **bucket** is a named collection of key/value pairs, like a table. The methods:

| Method | Does | Request line | Headers |
| --- | --- | --- | --- |
| `CREATE` | create an empty bucket | `CREATE <bucket> PIKA/1.0` | none |
| `DELETE` | delete a bucket and everything in it | `DELETE <bucket> PIKA/1.0` | none |
| `GET` | read one key, or list keys with `key: *` | `GET <bucket> PIKA/1.0` | `key` (required), `limit`, `after` |
| `SET` | store a value under a key | `SET <bucket> PIKA/1.0` | `key`, `value` (required), `append` |
| `REMOVE` | remove one key | `REMOVE <bucket> PIKA/1.0` | `key` (required) |

`DELETE` and `REMOVE` differ by one letter and both reply `removed`. The response fields tell them apart (`REMOVE` has `key`), but it is an easy mix-up when testing.

Example requests:

```text
CREATE sessions PIKA/1.0

{ "status": "created", "bucket": "sessions" }
```

```text
DELETE sessions PIKA/1.0

{ "status": "removed", "bucket": "sessions" }
```

```text
GET sessions PIKA/1.0
Key: *
Limit: 10

{
  "status": "found",
  "data": [ 
    { "key": "13", "value": { "sess_13" } },
    { "userid": {"sess_12", "sess_14"}},
  ],
  "bucket": "sessions"
}

```

```text
GET sessions PIKA/1.0
Key: userid_1020

{ "status": "found", "data": { "key": "userid_12", "value": {"seess_12", "sess13" } }, "bucket": "sessions" }

```

```text
SET sessions PIKA/1.0
key: userid_1020
value: sessionid_1020

{ "status": "stored", "bucket": "sessions", "userid_1020": "sessionid_1020" }

```

```text
SET sessions?append=true PIKA/1.0
key: userid_1020
value: sessionid_1020

{ "status": "stored", "bucket": "sessions", "userid_1020": {"sessionid_1019", "sessionid_1020"} }

```

```text
REMOVE sessions PIKA/1.0
key: userid_1020

{ "status": "removed", "bucket": "sessions", "key": "userid_1020" }
```


```text
CREATE sessions PIKA/1.0
SET sessions PIKA/1.0
key: userid_1020
value: sessionid_1020

{ "status": "created", "bucket": "sessions" }

```

```text
SET sessions PIKA/1.1
key: userid_1020
value: sessionid_1020

{ "status": "bad_request", "error": "protocol version mismatch" }
```

<!-- ## Design decisions

| Decision | Choice | Why |
| --- | --- | --- |
| Request body | None. Arguments go in headers | Keeps parsing simple |
| Request framing | Message ends at the blank line (`\r\n\r\n`) | No `Content-Length` needed without a body |
| Response format | JSON | Matches the examples |
| Response framing | Server closes the connection if client dies or exits; client reads until EOF | Simplest rule, no length or delimiter needed (this is how HTTP/1.0 worked) | -->
<!-- | Keep-alive | Not in 1.0 | `PIKA/1.0` leaves room for a 1.1 (see [Extensions](#8-extensions)) | -->
<!-- | Unknown version | Reply `bad_request` and close | Never guess what a future version means | -->

## Parse strictly

The more ways one message can be read, the more ways two programs can read it differently. That gap is the root cause of request-smuggling bugs in real HTTP servers. Pick one reading and reject the rest:

- Exactly one space between the three parts of the request line, and exactly three parts.
- Line endings are strictly `\r\n`. A bare `\n` is a `bad_request`.
- Method names are case-sensitive: `GET`, not `get`.
- Header names are case-insensitive: `KEY`, `Key`, `key`, `kEy` are same.
- Split a header at the **first** colon only, so values may contain colons. Trim spaces around the value.
- Duplicate headers (two `key:` lines) are a `bad_request`, not "last one wins".
- Restrict bucket and key characters (for example letters, digits, `-`, `_`, `.`) and their length.
- A NUL byte (`'\0'`) anywhere in the request is a `bad_request`. C string functions would silently cut the data at it.
- If multiple request comes together, first one wins.
<!-- ## Consequences of headers-only requests

- Values cannot contain `\r\n`, and a value cannot exceed the header line limit.
- Values are text. Arbitrary binary data doesn't fit (see [Extensions](#8-extensions)).
- A `GET` with `Limit` can return many pairs, so the JSON `data` array holds them all. -->

<!-- ## Semantics to decide and document

Suggested defaults are marked, but each choice is yours:

| Question | Suggested default |
| --- | --- |
| `GET` with a `key` header: return just that key? | Yes. Without `key`, list up to `Limit` pairs |
| `Limit` default and maximum | default 10, maximum 1000, reject larger or non-numeric |
| `Order` values | `asc` (default) or `desc`, sorted by raw key bytes |
| `GET` on a missing bucket or key | `not_found` |
| `DEL` on a missing key | `not_found`. Treating it as success makes `DEL` idempotent (safe to repeat), but hides client typos. Pick one |
| `SET` on an existing key | Overwrite, reply `stored` |
| Bucket creation | Implicit on first `SET` | -->

## Status values

| Status | When | Close after? |
| --- | --- | --- |
| `created` | `CREATE` returned data | yes |
| `found` | `GET` returned data | yes |
| `stored` | `SET` succeeded | yes |
| `removed` | `DELETE` & `REMOVE` succeeded | yes |
| `not_found` | bucket or key doesn't exist | yes |
| `bad_request` | malformed, too long, unknown version | yes |
| `error` | internal failure (out of memory, etc.) | yes |

Error responses are responses too, so a bad request still gets a `bad_request` JSON reply before the close.

<!-- ## Limits (starting points)

| Limit | Value |
| --- | --- |
| Header line length | 4 KB |
| Total header size | 16 KB |
| Number of headers | 16 |
| Bucket name / key length | 64 / 256 bytes |
| Read timeout per `recv` | 5 seconds |
| **Total** time to receive a request | 10 seconds |
| Max simultaneous connections | 256 |

Keep both timeouts. A per-`recv` timeout alone is not enough: a client that sends one byte every 4 seconds never triggers it, yet holds a thread or slot forever (a "slow loris" attack). The total deadline stops that. -->

---

# 2. Protocol Parser (What you will implement)

## Project setup

your folder structure should be look like this or 

```
```

parser.h file
```
```
parser.c
```
```

main.c
```
```

## Testing



# 3. Implement the hash table to store the key-value

hashmap.h file
```
```

hashmap.c file 
```
```

main.c
```
```

## Testing



















<!-- # 2. Framing in C

TCP is a byte stream with no message boundaries. One `recv()` can return half a message, or two. So:

1. Keep a buffer and its current length.
2. Append each `recv` result to it.
3. After each append, look for `\r\n\r\n`. To find a terminator split across two reads, search from `max(0, old_len - 3)`, not just the new bytes.
4. If the buffer is full and there is still no terminator, reply `bad_request` (too large) and close.
5. Only parse once the terminator is found.

Decide what to do with bytes after the terminator. With no keep-alive, ignoring them or rejecting the request are both defensible, but pick one.

Parser test list (write these before any socket code):

- the whole message in one chunk
- the message fed **one byte at a time**
- the terminator split at each of the three possible places
- an empty request, and only a blank line
- a bare `\n` instead of `\r\n`
- a header line of exactly the limit, and one byte over
- a NUL byte inside a value
- a missing colon, a missing key, a duplicate header
- a non-numeric, negative, or huge `Limit`

## JSON encoder

You only need an **encoder**, not a parser. Escape every string value before writing it:

| Input | Output |
| --- | --- |
| `"` | `\"` |
| `\` | `\\` |
| newline, carriage return, tab | `\n`, `\r`, `\t` |
| other bytes below `0x20` | `\u00XX` |

Details that bite:
- The worst case is 6 output bytes per input byte (`\u00XX`), so size the output buffer for that, or grow it, and check bounds on every write.
- JSON must be valid UTF-8. Decide whether `SET` rejects values that aren't valid UTF-8 or the encoder replaces bad bytes. Silent corruption is the worst option.
- Write the escaper as its own function, test it with nasty inputs, and make the tests include a value containing `"}` to prove the output can't be broken out of.
- Use explicit lengths everywhere, not `strlen`, in case a `'\0'` ever gets through.

---

# 3. C-specific pitfalls

## Strings and numbers

- **Don't trust null terminators in network data.** Use lengths (`memchr`, `memcmp`). Never `strcpy` or `sprintf`. Use `snprintf` and check its return value. Cap every buffer.
- **Parse numbers carefully.** `strtol` with `errno` and `endptr` checks, then range-check the result. Reject `-5`, `10abc`, empty text, and overflow.
- **Watch size arithmetic.** `len + n` can overflow `size_t` when growing a buffer. Check before adding.

## The hash table

C has no built-in map, so write one: an array of buckets, each a linked list ("separate chaining").

- The table **copies** keys and values on `SET` and owns the copies. It never keeps pointers into the connection buffer, which is reused or freed.
- Choose a hash function (FNV-1a is short and fine). A table with a fixed hash can be attacked by sending many keys that collide, turning lookups into O(n). A random seed chosen at startup defends against that.
- Decide how it grows. Resizing changes every bucket, so it must happen under the write lock in the threaded version.
- `Order: asc|desc` needs sorted output: collect the keys at `GET` time and sort with `qsort` and a comparator (a function pointer, as in your pointer notes), or keep a sorted structure.
- Deleting a node needs care when another thread may be reading it. This is one reason to copy results out under the lock (see 5.1).

## Ownership

For every `malloc`, write down who frees it, including on error paths. A common pattern is one `cleanup:` label per function: every error path does `goto cleanup;` and the label frees everything that might be allocated (so initialize all pointers to `NULL` at the top, since `free(NULL)` is safe).

## Socket setup

- `SO_REUSEADDR` on the listener, set before `bind`.
- Ignore `SIGPIPE` with `signal(SIGPIPE, SIG_IGN)`, which works on both Linux and macOS. (`MSG_NOSIGNAL` is Linux only.)
- Loop on `send()` until everything is written, because it can write fewer bytes than you asked.
- `SO_RCVTIMEO` so a silent client can't hang a thread, plus the total deadline from section 1.

## Clean close

Closing a socket while the client's data is still unread in your receive buffer makes the OS send a TCP **reset** (RST, an abrupt abort instead of a normal close). The client may then throw away your response before reading it. This is most likely exactly when you reject a bad request early, without reading everything.

Safe sequence:
1. Send the response.
2. `shutdown(fd, SHUT_WR)`, which tells the client "no more data from me" but keeps the socket open for reading.
3. Read and discard incoming bytes for a short, bounded time (or until EOF).
4. `close(fd)`.

This is called a **lingering close**.

## Compile flags

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined main.c -o pika
```

Use these from day one. ThreadSanitizer needs a separate build (see [Testing and tools](#7-testing-and-tools)).

---

# 4. Build order

| Step | Do | Done when |
| --- | --- | --- |
| 1 | Parser as a pure function (buffer in, request struct out) | passes the whole test list in section 2, no sockets exist yet |
| 2 | JSON encoder | escaper tests pass |
| 3 | Hash table plus unit tests | add, overwrite, delete, sorted listing all verified |
| 4 | Single-client TCP server: `socket`, `bind`, `listen`, `accept` | `nc` gets correct replies, sanitizers are clean |
| 5 | Threads and a lock | many parallel clients, no races (TSan clean) |
| 6 | Limits, timeouts, lingering close | slow-loris and oversize tests handled |
| 7 | Stress test | many clients in parallel, no leaks, no crashes |

---

# 5. Concurrency models

Compare all three. The parser, store and JSON encoder are shared, and only the connection-handling layer changes.

| | Thread per connection | Event loop (`poll`) | `io_uring` |
| --- | --- | --- | --- |
| Threads | one per client | one | one or few |
| Lock on the store | yes | no | no (if single-threaded) |
| Many idle connections | costly (a thread each) | cheap | cheap |
| Uses several CPU cores | yes | no | optional |
| Code complexity | low | medium (state machines) | high |
| Platform | POSIX | POSIX (`epoll`: Linux, `kqueue`: macOS/BSD) | Linux only |

## 5.1 Thread per connection (the main exercise)

- `pthread_create` per accepted connection. Each thread reads, parses, handles, responds, and closes.
- Pass each thread its **own** heap-allocated connection struct (fd and client address). Passing the address of the loop's local variable is the `&i` bug from your threads note.
- Detach the threads (`pthread_detach`) since nobody joins them, or use a fixed pool from your threads note.
- **Limit the number of live threads** (a counter or semaphore). Without a cap, a flood of connections creates unbounded threads.
- Protect the store with `pthread_rwlock_t` (many readers or one writer). On glibc the default can let a steady stream of readers starve a writer, so test writes under heavy reading.
- **Never hold the lock while doing network I/O.** Copy the data you need out of the store under the lock, release it, then build and send the response. A slow client would otherwise freeze every other thread.
- Lesson: data races. Reproduce one on purpose first (many clients doing `SET` at once, no lock), watch it fail or get flagged by ThreadSanitizer, then add the lock.

## 5.2 Event loop (`poll`, then optionally `epoll`)

One thread handles every connection. "Ready" means "a read or write would not block right now" (this is **readiness-based** I/O).

- All sockets are non-blocking (`O_NONBLOCK`): `recv` or `send` that can't proceed returns `-1` with `EAGAIN` or `EWOULDBLOCK` instead of waiting.
- `poll` takes an array of `pollfd` and reports which are ready. Ask for `POLLIN` (readable) on the listener and on connections still reading.
- Ask for `POLLOUT` (writable) **only while you have data left to send**. Asking for it permanently makes the loop spin at 100% CPU.
- When the listener is readable, call `accept` in a loop until it returns `EAGAIN`.
- Handle `POLLHUP` and `POLLERR` by closing that connection.
- Each connection needs a `struct conn`: fd, read buffer and length, response buffer and bytes sent so far, state, and last-activity time.

State machine for one connection:

| State | Waits for | Moves to |
| --- | --- | --- |
| `READING` | `POLLIN`, appends to the buffer | `HANDLING` once `\r\n\r\n` found, or `CLOSING` on error or timeout |
| `HANDLING` | nothing, runs immediately | `WRITING` |
| `WRITING` | `POLLOUT`, remembers the offset | `CLOSING` when everything is sent |
| `CLOSING` | lingering close (see section 3) | removed from the poll set |

- Partial writes: remember how much was sent and wait for `POLLOUT` again.
- Idle and total timeouts are tracked by hand: use `poll`'s timeout argument and sweep connections for expired deadlines.
- No lock is needed, since only one thread touches the store.
- Downside: one slow handler blocks everyone (a huge sorted `GET`, say), and only one CPU core is used.
- `poll` rescans every descriptor on each call, so it slows down with thousands of connections. `epoll` (Linux) and `kqueue` (macOS/BSD) keep the interest list in the kernel and only report the ready ones. Terms: **level-triggered** reports "still ready" on every call until you drain it (simpler, matches `poll`), and **edge-triggered** reports only when the state changes (faster, but you must read until `EAGAIN` or you lose events).

## 5.3 `io_uring` (optional stretch)

Instead of asking "is it ready?" and then calling `recv` yourself (readiness), you submit "recv into this buffer" and the kernel tells you when it is **done** (**completion-based** I/O).

- Two ring buffers shared with the kernel: the **submission queue** (requests you make) and the **completion queue** (results). Use the `liburing` library, not raw syscalls.
- Each request carries a `user_data` value (use a pointer to your `struct conn`) that comes back with the completion, so you know which connection finished.
- Buffers must stay alive and untouched while an operation is in flight.
- Timeouts are separate requests, linked to the one they limit.
- Linux only. It first appeared in kernel 5.1, but accept and recv/send support came in later releases (roughly 5.5 to 5.6), so check the `liburing` documentation for the minimum your features need. Some container setups and cloud sandboxes disable it. If you are on macOS, run it in a Linux VM or container.

---

# 6. Memory management: manual first, then GC

Wrap allocation from day one, so the swap becomes a compile flag:

```c
void *pk_alloc(size_t n);
void *pk_realloc(void *p, size_t n);   // growing buffers need this too
void  pk_free(void *p);
```

- **Manual build:** these call `calloc`, `realloc`, `free`.
- **GC build** (`-DUSE_GC`): `pk_alloc` calls `GC_MALLOC`, `pk_realloc` calls `GC_REALLOC`, and `pk_free` does nothing.
- Make the manual `pk_alloc` return **zeroed** memory (`calloc`). `GC_MALLOC` returns zeroed memory, `malloc` does not. If the two builds differ here, a bug can appear in only one of them and your comparison isn't fair.
- Optionally add `pk_alloc_atomic` for data with no pointers inside (strings, byte buffers). A GC can skip scanning it, which is faster.

## Boehm GC (`libgc`)

- Install: `libgc-dev` on Debian/Ubuntu, `bdw-gc` with Homebrew on macOS. Call `GC_INIT()` once, link with `-lgc`.
- With threads, define `GC_THREADS` before `#include <gc.h>`, and include it before `<pthread.h>`, so the collector can track your threads.
- It is **conservative**: it can't tell integers from pointers, so it treats anything that looks like a pointer as one. The result is that it may keep some garbage alive (a small leak) but won't free anything still in use.
- **Important rule:** the collector only scans memory it knows about (stacks, globals, and its own objects). A pointer stored **only** inside plain `malloc`'d memory is invisible, and the object it points to can be freed while still in use. So in the GC build, route everything through `pk_alloc`, including the hash table nodes. Don't mix.
- AddressSanitizer can conflict with it. Test memory bugs in the manual build, and test the GC build without ASan.

## Arena alternative

An **arena** (bump allocator) is one big block you hand out from in order. You never free individual pieces, you release the whole block at once. Use one arena per request and free it right after the response. This fits PIKA 1.0, since each connection is short-lived.

The catch is ownership: anything that must outlive the request, such as keys and values stored in the hash table, has to be **copied out** to long-lived memory (`malloc`) before the arena is released.

## What to compare after both versions exist

| Measure | How |
| --- | --- |
| Lines of code and error-path complexity | how many `cleanup:` labels disappeared |
| Peak memory under the same stress test | Linux: `/usr/bin/time -v` or `VmRSS` in `/proc/self/status`. macOS: `/usr/bin/time -l`. Or `getrusage` (`ru_maxrss` is in KB on Linux and in **bytes** on macOS) |
| Latency pauses from collection | record the slowest requests (p99 and max), not just the average |
| What GC did **not** fix | file descriptors still need `close()`, and logic bugs remain |

Keep a **bug log** during the manual phase:

| Bug | Found by | Would GC have removed it? |
| --- | --- | --- |
| e.g. leak on parse-error path | ASan | yes |
| e.g. double free in `DEL` | ASan | yes |
| e.g. fd leak on timeout | `lsof` | no |

---

# 7. Testing and tools

Manual testing with `nc` (netcat):

```bash
printf 'GET /b PIKA/1.0\r\nLimit: 10\r\n\r\n' | nc 127.0.0.1 8080
```

`nc` flags and behavior differ between Linux and macOS versions. If it hangs after printing the reply, check `nc`'s timeout or "close on EOF" option for your version. Always use `printf` with `\r\n`, since typing in a terminal sends only `\n`.

Other tools:
- **Fuzzing**: because the parser is a pure function, you can feed it random bytes automatically with libFuzzer (`clang -fsanitize=fuzzer,address`) or AFL++. This finds crashes you would never write tests for.
- **Load generator**: write a small client (in C or Go) that opens many connections, sends mixed `SET`/`GET`/`DEL`, and checks the replies. A plain HTTP tool like `wrk` won't speak PIKA.
- `lsof -i :8080` shows connections. `strace -c` (Linux) counts syscalls, which is useful for comparing the three models.

Which build for which job:

| Build | Flags | Use |
| --- | --- | --- |
| Development | `-fsanitize=address,undefined -g` | day-to-day |
| Race check | `-fsanitize=thread -g` | threaded version (cannot combine with ASan) |
| GC | `-DUSE_GC -lgc`, no ASan | GC comparison |
| Leak check | `valgrind` (Linux), `leaks --atExit` (macOS) | extra confirmation |

---

# 8. Extensions

- **Keep-alive (PIKA/1.1):** let one connection carry several requests. This breaks the "response ends at EOF" rule, so responses need their own framing, for example a `Content-Length` header or a length prefix. This is the main reason the version token exists.
- **Binary values:** add a body with a length header (the one place the "no body" decision changes).
- **Persistence:** append every `SET` and `DEL` to a log file, and replay it at startup.
- **TTL:** an optional `ttl` header that expires keys.
- **Auth:** a token header, checked before anything else.
- **Metrics:** count requests, errors and timeouts, and expose them with a special `GET`.

---

# Summary

- The protocol is simple on purpose. Parse it **strictly**, and reject anything ambiguous.
- TCP has no message boundaries, so buffer, search for the terminator, and bound everything: size, time, and connection count.
- Use a total deadline, not just a per-read timeout.
- Close politely (`shutdown`, drain, `close`) so the client actually gets your reply.
- In the threaded version, don't hold locks during I/O, and pass each thread its own connection struct.
- Compare thread, `poll` and `io_uring` versions with the same tests and the same stress load.
- A conservative GC only sees memory it manages, so don't mix it with plain `malloc`.
- Keep the allocation wrapper identical in both builds (zeroed memory) so the comparison is fair. -->





---
title: "Capstone Project: PIKA, a Toy Key-Value Protocol and Server in C"
date: 2026-10-27
tags:
  - c
  - networking
  - sockets
  - capstone
summary: A text key-value protocol over TCP, built in C: wire format, strict parser, hash-table store and JSON responses (sockets and concurrency follow)
---

# Overview

PIKA is a small text protocol for a key-value store, served over TCP. The project has three stages: **protocol, sockets, multithreading**. This note covers the first stage: the wire format, a strict parser, the in-memory store and the response writer, all testable without a single socket. The socket server and the concurrency models come in the next notes.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [The protocol (PIKA/1.0)](#1-the-protocol-pika10) | Wire format, strict parsing, statuses, limits, open decisions |
| 2 | [Protocol parser](#2-protocol-parser) | Interface, validation rules, parsing steps, test table, Exercise 01 |
| 3 | [The key-value store](#3-the-key-value-store) | Data model, hash map spec, ownership, ordering, Exercise 02 |
| 4 | [Response writer](#4-response-writer) | JSON output, escaping, Exercise 03 |
| 5 | [Platform notes](#5-platform-notes) | macOS vs Linux for tools used here |
| 6 | [Common errors](#6-common-errors) | Bugs you will probably hit |
| 7 | [Summary](#summary) | Recap |

---

# 1. The protocol (PIKA/1.0)

A stripped-down, HTTP-like request/response protocol, inspired by HTTP/1.0.

```text
METHOD <bucket-name> PIKA/1.0\r\n
header: value\r\n
\r\n
```

A **bucket** is a named collection of key/value pairs, like a table. A key maps to a **list of values** (a plain `SET` replaces the list with one value, `append: true` adds to it). Requests have **no body**: the blank line ends the request.

| Method | Does | Request line | Headers |
| --- | --- | --- | --- |
| `CREATE` | create an empty bucket | `CREATE <bucket> PIKA/1.0` | none |
| `DELETE` | delete a bucket and everything in it | `DELETE <bucket> PIKA/1.0` | none |
| `GET` | read one key, or list keys with `key: *` | `GET <bucket> PIKA/1.0` | `key` (required), `limit`, `after` |
| `SET` | store a value under a key | `SET <bucket> PIKA/1.0` | `key`, `value` (required), `append` |
| `REMOVE` | remove one key | `REMOVE <bucket> PIKA/1.0` | `key` (required) |

`DELETE` and `REMOVE` differ by one letter and both reply `removed`. The response fields tell them apart (`REMOVE` has `key`), but it is an easy mix-up when testing.

## How to read the examples

In each block, the lines up to and including the first blank line are the **request** (the blank line is the request terminator, `\r\n\r\n`). Everything after it is the **response**. A response is just the JSON text; the server closes the connection afterwards, so the client reads until EOF and no length field is needed.

```text
CREATE sessions PIKA/1.0

{ "status": "created", "bucket": "sessions" }
```

```text
DELETE sessions PIKA/1.0

{ "status": "removed", "bucket": "sessions" }
```

List keys (`key: *`, at most 10 entries):

```text
GET sessions PIKA/1.0
Key: *
Limit: 10

{
  "status": "found",
  "data": [
    { "key": "userid_1019", "value": ["sessionid_1019"] },
    { "key": "userid_1020", "value": ["sessionid_1019", "sessionid_1020"] }
  ],
  "bucket": "sessions"
}
```

Read one key:

```text
GET sessions PIKA/1.0
Key: userid_1020

{ "status": "found", "data": { "key": "userid_1020", "value": ["sessionid_1019", "sessionid_1020"] }, "bucket": "sessions" }
```

Store a value (replaces any existing list):

```text
SET sessions PIKA/1.0
key: userid_1020
value: sessionid_1020

{ "status": "stored", "bucket": "sessions", "key": "userid_1020", "value": ["sessionid_1020"] }
```

Append to the list:

```text
SET sessions PIKA/1.0
key: userid_1020
value: sessionid_1020
append: true

{ "status": "stored", "bucket": "sessions", "key": "userid_1020", "value": ["sessionid_1019", "sessionid_1020"] }
```

Remove one key:

```text
REMOVE sessions PIKA/1.0
key: userid_1020

{ "status": "removed", "bucket": "sessions", "key": "userid_1020" }
```

Missing key:

```text
GET sessions PIKA/1.0
key: nope

{ "status": "not_found", "error": "key not found" }
```

Two requests in one send. Only the first is handled, the rest is discarded:

```text
CREATE sessions PIKA/1.0

SET sessions PIKA/1.0
key: userid_1020
value: sessionid_1020

{ "status": "created", "bucket": "sessions" }
```

Wrong version:

```text
SET sessions PIKA/1.1
key: userid_1020
value: sessionid_1020

{ "status": "bad_request", "error": "protocol version mismatch" }
```

## Parse strictly

The more ways one message can be read, the more ways two programs can read it differently. That gap is the root cause of **request smuggling** bugs in real HTTP servers (an attacker hides a second request inside the first so a proxy and a server disagree about where requests start). Pick one reading and reject the rest:

- Exactly one space between the three parts of the request line, and exactly three parts. No leading or trailing space.
- Line endings are strictly `\r\n`. A bare `\n` is a `bad_request`, and so is a bare `\r`.
- Method names are case-sensitive: `GET`, not `get`.
- Header names are case-insensitive: `KEY`, `Key`, `key`, `kEy` are the same header. No space between the name and the colon (`key : a` is a `bad_request`).
- Split a header at the **first** colon only, so values may contain colons. Trim spaces around the value.
- Duplicate headers are a `bad_request`, not "last one wins". "Duplicate" is judged after lowercasing the name, so `Key:` and `key:` collide.
- Unknown headers, and headers that don't belong to the method (`key:` on `CREATE`), are a `bad_request`.
- Restrict bucket and key characters (letters, digits, `-`, `_`, `.`) and their length. This also means `*` can never be a real key, which is why `key: *` can safely mean "list".
- A NUL byte (`'\0'`) anywhere in the request is a `bad_request`. C string functions would silently cut the data at it.
- If several requests arrive in one read, **only the first is handled** and the remaining bytes are discarded. This is safe because requests have no body, so the first blank line is always the end of the request.

## Status values

| Status | When | Close after? |
| --- | --- | --- |
| `created` | `CREATE` succeeded | yes |
| `found` | `GET` returned data | yes |
| `stored` | `SET` succeeded | yes |
| `removed` | `DELETE` & `REMOVE` succeeded | yes |
| `not_found` | bucket or key doesn't exist | yes |
| `bad_request` | malformed, too long, unknown version | yes |
| `error` | internal failure (out of memory, etc.) | yes |

Error responses are responses too, so a bad request still gets a `bad_request` JSON reply before the close.

## Limits

Every limit is a `#define` (not a `const int`, so it can size an array). These numbers are suggestions; the point is that **every** length has a cap.

| Limit | Suggested value | Why |
| --- | --- | --- |
| Whole request (up to the blank line) | 4096 bytes | bounds the read buffer, stops endless headers |
| Bucket name, key | 1 to 64 chars | fixed arrays, no malloc in the parser |
| Value | up to 1024 bytes | must fit inside the request limit |
| Header lines | 8 | there are only 5 known headers |
| Values per key | 100 | `append` in a loop must not eat all memory |
| Keys per bucket / buckets | 10000 / 1000 | same reason |
| `limit` header | 1 to 1000, default 100 | bounds a list response |

## Open decisions

The spec above doesn't say what happens in these cases. Pick an answer, write it down, and add a test for each.

| Case | Suggested answer |
| --- | --- |
| `CREATE` on an existing bucket | new status `exists` (close after) |
| `DELETE` / `SET` / `GET` / `REMOVE` on a missing bucket | `not_found` |
| `REMOVE` on a missing key | `not_found` |
| `GET` with `key: *` on an empty bucket | `found` with `"data": []` |
| `SET` with `append: true` on a missing key | create the key |
| `after` header | first key strictly greater than `after` in byte order; `after` need not exist |
| List response grows too large | stop at a size cap (for example 64 KB); the client continues with `after` = last key seen |

---

# 2. Protocol parser

The parser is a pure function: bytes in, a filled struct or an error out. No sockets, no globals, no `malloc`. That makes it easy to test and later easy to drive from a socket loop.

## Project setup

```text
pika/
├── Makefile
├── parser.h        interface and limits
├── parser.c        pika_parse() and helpers
├── parser_test.c   table-driven tests, has its own main()
├── hashmap.h       generic string-keyed map (section 3)
├── hashmap.c
├── store.h         buckets, keys, value lists (section 3)
├── store.c
├── response.h      JSON writer (section 4)
├── response.c
└── main.c          later: the socket server
```

Build line for the parser tests:

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined parser.c parser_test.c -o parser_test
```

## Interface

A suggested `parser.h`. It has declarations only; you write `parser.c`.

```c
#ifndef PIKA_PARSER_H
#define PIKA_PARSER_H

#include <stdbool.h>
#include <stddef.h>

#define PIKA_MAX_REQUEST 4096
#define PIKA_MAX_NAME    64     /* bucket and key */
#define PIKA_MAX_VALUE   1024
#define PIKA_DEFAULT_LIMIT 100
#define PIKA_MAX_LIMIT     1000

typedef enum {
    PIKA_CREATE, PIKA_DELETE, PIKA_GET, PIKA_SET, PIKA_REMOVE
} pika_method;

typedef enum {
    PARSE_OK,          /* one full request parsed */
    PARSE_INCOMPLETE,  /* no blank line yet: read more bytes */
    PARSE_BAD_REQUEST  /* reject; *err explains why */
} parse_status;

typedef struct {
    pika_method method;
    char   bucket[PIKA_MAX_NAME + 1];
    char   key[PIKA_MAX_NAME + 1];     /* "*" means list (GET only) */
    char   value[PIKA_MAX_VALUE + 1];
    char   after[PIKA_MAX_NAME + 1];
    bool   has_key, has_value, has_after;
    bool   append;
    size_t limit;                      /* PIKA_DEFAULT_LIMIT if absent */
    size_t consumed;                   /* bytes used by this request */
} pika_request;

/* buf is NOT NUL-terminated and may contain any byte. */
parse_status pika_parse(const char *buf, size_t len,
                        pika_request *out, const char **err);

#endif
```

`PARSE_INCOMPLETE` is for the socket stage, where a request can arrive in several `read()` calls. Designing it in now saves a rewrite.

## Validation rules

| Method | Required headers | Optional headers | Anything else |
| --- | --- | --- | --- |
| `CREATE` | none | none | `bad_request` |
| `DELETE` | none | none | `bad_request` |
| `GET` | `key` | `limit`, `after` (only when `key` is `*`) | `bad_request` |
| `SET` | `key`, `value` | `append` (`true` or `false` exactly) | `bad_request` |
| `REMOVE` | `key` | none | `bad_request` |

| Field | Rule |
| --- | --- |
| bucket, key | 1 to 64 chars from `[A-Za-z0-9._-]`; `key` may also be exactly `*` on `GET` |
| value | 1 to 1024 bytes of printable ASCII (`0x20` to `0x7E`); this keeps JSON output simple, see section 4 |
| `limit` | decimal digits only, 1 to 1000 |
| `append` | `true` or `false`, lowercase |

## Parsing steps (hints)

1. **Find the end of the head.** Search for `\r\n\r\n` yourself with a loop. `memmem` is not standard C (it exists on macOS and glibc, but don't depend on it). If you don't find it: if `len >= PIKA_MAX_REQUEST` return bad request ("request too large"), otherwise `PARSE_INCOMPLETE`.
2. **Scan the head once** for NUL, bare `\n` and bare `\r`. Do this before any splitting, so later steps can trust the bytes.
3. **Request line.** Find the two spaces by index. Exactly two, with non-empty parts between. Check version last or first, but **pick one order** and keep it, because your tests pin the error message.
4. **Headers.** For each line: find the first `:`, validate the name, lowercase it, trim the value, then look the name up in a small table of `{name, method mask}`.
5. **Duplicates.** Keep a bitmask `seen`; if the bit is already set, reject.
6. **Numbers.** `limit` must be digits only. If you use `strtoul`, check `errno`, that `endptr` reached the end, and that the string wasn't empty (it accepts a leading space and a `+` or `-`). A hand-written digit loop is just as good.
7. **Set `consumed`** to the offset just past the first `\r\n\r\n`.
8. **ctype gotcha.** `isalnum(c)` with a negative `char` is undefined behaviour. Cast to `unsigned char` first, or write your own range checks.
9. **Don't use `strtok`.** It collapses repeated delimiters, so it would silently accept two spaces and break the "exactly one space" rule.

## Exercise 01: Strict parser

Implement `pika_parse()` so every row of the test table below passes.

Requirements:
- No `malloc`, no globals, no static buffers.
- Never read past `len`, and never rely on a NUL terminator in `buf`.
- On `PARSE_BAD_REQUEST`, set `*err` to a string literal (it is later put in the JSON `error` field).
- Clean under `-fsanitize=address,undefined` with `-Wall -Wextra`.

Hints:
- Write `static` helper functions (`valid_name`, `trim`, `ascii_lower`) in `parser.c`; keep them out of the header.
- Table-driven tests: an array of `{ const char *name; const char *input; size_t len; parse_status want; }`.
- An input containing `\0` can't use `strlen`. Use `sizeof(literal) - 1` for its length.

## Testing

Inputs are C string literals. `\0` makes the case need an explicit length.

| # | Input | Expected |
| --- | --- | --- |
| 1 | `"CREATE sessions PIKA/1.0\r\n\r\n"` | OK, `CREATE`, bucket `sessions`, `consumed` 28 |
| 2 | `"GET sessions PIKA/1.0\r\nKey: *\r\nLimit: 10\r\n\r\n"` | OK, key `*`, limit 10 |
| 3 | `"GET sessions PIKA/1.0\r\nKEY: a\r\n\r\n"` | OK, key `a` (header name case-insensitive) |
| 4 | `"SET sessions PIKA/1.0\r\nkey: a\r\nvalue: x:y\r\n\r\n"` | OK, value `x:y` (first colon only) |
| 5 | `"SET sessions PIKA/1.0\r\nkey:   a  \r\nvalue: x\r\n\r\n"` | OK, key `a` (trimmed) |
| 6 | `"CREATE sessions PIKA/1.0\r\n"` | INCOMPLETE |
| 7 | `"CREATE sessions PIKA/1.0\r\n\r\nSET sessions PIKA/1.0\r\n\r\n"` | OK, `consumed` 28 (second request ignored) |
| 8 | `"get sessions PIKA/1.0\r\n\r\n"` | bad: method is case-sensitive |
| 9 | `"CREATE  sessions PIKA/1.0\r\n\r\n"` | bad: two spaces |
| 10 | `"CREATE sessions PIKA/1.0\n\n"` | bad: bare `\n` |
| 11 | `"CREATE sessions PIKA/1.1\r\n\r\n"` | bad: version mismatch |
| 12 | `"CREATE bad/name PIKA/1.0\r\n\r\n"` | bad: illegal bucket character |
| 13 | bucket name of 65 characters | bad: name too long |
| 14 | `"GET sessions PIKA/1.0\r\nkey: a\r\nKey: b\r\n\r\n"` | bad: duplicate header |
| 15 | `"GET sessions PIKA/1.0\r\nkey : a\r\n\r\n"` | bad: space before colon |
| 16 | `"CREATE sessions PIKA/1.0\r\nkey: a\r\n\r\n"` | bad: header not allowed on `CREATE` |
| 17 | `"SET sessions PIKA/1.0\r\nkey: a\r\n\r\n"` | bad: `value` missing |
| 18 | `"SET sessions PIKA/1.0\r\nkey: a\r\nvalue: x\0y\r\n\r\n"` | bad: NUL byte |
| 19 | `"GET sessions PIKA/1.0\r\nkey: a\r\nlimit: 10\r\n\r\n"` | bad: `limit` only with `key: *` |
| 20 | `"GET sessions PIKA/1.0\r\nkey: *\r\nlimit: 0\r\n\r\n"`, and `limit: abc`, `limit: 10x`, `limit: -1` | bad: invalid `limit` |
| 21 | 4097 bytes with no blank line | bad: request too large |

---

# 3. The key-value store

The store is independent of the parser and the network. It takes plain C values and returns status codes, so it can be tested on its own.

## Data model

```text
store
 └── bucket table          name  -> bucket
      └── key table        key   -> value list
           └── values      ["sessionid_1019", "sessionid_1020"]
```

Both tables are the same generic **hash map**, so you write the hash table once. A hash map turns a key into an array index with a hash function; two keys landing on the same index is a **collision**.

## Hash map interface

A suggested `hashmap.h`, declarations only:

```c
#ifndef HASHMAP_H
#define HASHMAP_H

#include <stddef.h>

typedef struct hashmap hashmap;

/* free_value is called on a value when it is replaced, removed or the map is destroyed. */
hashmap *hm_create(void (*free_value)(void *));
void     hm_destroy(hashmap *m);

int      hm_put(hashmap *m, const char *key, void *value);  /* 0 ok, -1 no memory */
void    *hm_get(const hashmap *m, const char *key);         /* NULL if absent */
int      hm_remove(hashmap *m, const char *key);            /* 1 removed, 0 absent */
size_t   hm_count(const hashmap *m);

/* stop early if fn returns non-zero */
void     hm_foreach(const hashmap *m,
                    int (*fn)(const char *key, void *value, void *ctx),
                    void *ctx);

#endif
```

## Design choices

| Choice | Suggested | Reason |
| --- | --- | --- |
| Collision handling | **chaining**: each slot is a singly linked list | simplest to get right, removal is easy |
| Hash function | FNV-1a, 64-bit (offset `14695981039346656037`, prime `1099511628211`, per byte: XOR then multiply) | tiny, good enough for strings |
| Slot count | power of two, start at 16 | index is `hash & (n - 1)`, no modulo |
| Growth | double when `count > 0.75 * slots` (**load factor** = entries per slot) | keeps chains short |
| Rehash | relink existing nodes into the new array; don't copy keys again | one allocation, no key churn |
| Key storage | the map owns its own copy of each key | the caller's buffer may disappear |

## Ownership rules

| Situation | Who owns what |
| --- | --- |
| `hm_put` succeeds | the map owns `value` from now on |
| `hm_put` fails (returns -1) | the caller still owns `value` and must free it |
| `hm_put` on an existing key | the old value goes to `free_value`, then the new one is stored |
| `hm_get` | returns a borrowed pointer, valid only until the next `hm_put`/`hm_remove` on that map |
| `hm_destroy` | frees every key, every value (via callback) and the table |

The value stored in the key table is a `struct { char **items; size_t count, cap; }`. Its free function frees each string, the array, then the struct.

## Store operations

| Operation | Result codes |
| --- | --- |
| `store_create(name)` | OK, EXISTS, NOMEM |
| `store_delete(name)` | OK, NO_BUCKET |
| `store_set(bucket, key, value, append)` | OK, NO_BUCKET, LIMIT, NOMEM |
| `store_get(bucket, key)` | OK + borrowed list, NO_BUCKET, NO_KEY |
| `store_remove(bucket, key)` | OK, NO_BUCKET, NO_KEY |
| `store_list(bucket, after, limit, ...)` | OK, NO_BUCKET |

## Listing and `after`

A hash map has **no useful order**: iteration follows the slots, and the order changes when the table grows. `after` only makes sense with a stable order, so for `store_list`:

1. Collect pointers to all keys with `hm_foreach` into an array.
2. Sort the array with `qsort` and `strcmp`.
3. Skip every key `<= after`, then take up to `limit`.

This costs O(n log n) per list call. That's fine for a toy, and you can note the trade-off. Keep the store behind this small interface so a lock can wrap it in the threading stage.

## Exercise 02: Hash map and store

Requirements:
- `hashmap.c` passes the tests below with no leaks and no ASan or UBSan reports.
- Check every allocation. If growth fails midway, the map must still be valid and keep every entry (the old table stays in use).
- `store.c` enforces the limits table from section 1 and returns `LIMIT` when exceeded.
- No `strcpy` or `strncpy`: use `malloc` + `memcpy` with the length you already know.

Hints:
- Write `hm_put`, `hm_get` and `hm_remove` first with a fixed slot count; add growth afterwards and check that old tests still pass.
- To remove from a singly linked chain, keep a pointer to the **link** (`node **`), not just the node.
- Test growth by inserting 10000 keys, then checking all of them, then removing all.

## Testing

| # | Test | Expect |
| --- | --- | --- |
| 1 | put `a`, get `a` | same pointer back |
| 2 | get a key never put | `NULL` |
| 3 | put `a` twice with different values | second value returned, `free_value` called once for the first |
| 4 | remove `a`, then get `a` | `NULL`, `hm_count` decreases |
| 5 | remove a key that doesn't exist | returns 0 |
| 6 | insert 10000 distinct keys | all retrievable, `hm_count` is 10000 |
| 7 | remove all 10000 | `hm_count` is 0, no leaks |
| 8 | two keys forced into one chain (small table) | both retrievable, removing one keeps the other |
| 9 | `store_set` without append, then with append | list has 1 value, then 2 in order |
| 10 | `store_list` after inserting `b`, `a`, `c` | order `a`, `b`, `c`; `after: a` gives `b`, `c` |
| 11 | `store_delete` on a bucket with keys | no leaks |

---

# 4. Response writer

Every response is JSON text built by hand. Because a stored value came from a client, the writer must make sure a value can never break out of its JSON string.

| Rule | Why |
| --- | --- |
| Allow only printable ASCII in values (section 2) | no control characters, no UTF-8 decoding in this stage |
| Escape `"` as `\"` and `\` as `\\` when writing a string | a value like `x", "status": "stored` would otherwise inject fields |
| Write into a caller-provided buffer with `snprintf`-style truncation checks | a long list can exceed the buffer |
| If the response doesn't fit, report it and let the caller send an `error` | never silently send half a JSON document |

## Exercise 03: JSON writer

Requirements:
- Write `response.c` with one function per shape: simple status, status with bucket, single data, list data, error.
- Each function writes into `char *buf, size_t cap` and returns the length written, or `0` / `(size_t)-1` when it doesn't fit (your choice, but document it).
- Test that an escaped output re-parses correctly (use `python3 -c "import json,sys; json.load(sys.stdin)"` on the output).

Hints:
- Write a tiny `append_str(buf, cap, pos, s)` helper that checks space every time, then build everything on top of it.
- Test with a value containing `"` and `\`.

---

# 5. Platform notes

| Task | Linux | macOS |
| --- | --- | --- |
| Compiler | `gcc` is GNU gcc | `gcc` is usually Apple clang; flags in this note work with both |
| Leak check | LeakSanitizer is part of `-fsanitize=address` | use `leaks --atExit -- ./parser_test` |
| Memory errors | `valgrind ./a.out` | valgrind isn't practical on Apple Silicon; stay on ASan |
| Sending raw bytes | `printf 'CREATE a PIKA/1.0\r\n\r\n' \| nc 127.0.0.1 7000` | same command; `nc -C` sends CRLF line endings |
| Who owns a port | `ss -ltnp` | `lsof -iTCP -sTCP:LISTEN -n -P` |

---

# 6. Common errors

| Symptom | Cause | Fix |
| --- | --- | --- |
| Duplicate-header test passes, but `Key` vs `key` isn't caught | compared names before lowercasing | lowercase first, then look up |
| `isalnum` crashes or misbehaves on high bytes | negative `char` passed to ctype | cast to `unsigned char` |
| Double spaces accepted | split with `strtok` | find spaces by index |
| `limit: 10x` accepted | `strtoul` result used without checking `endptr` | check `endptr` and `errno` |
| NUL test passes the wrong way | `strlen` used on a buffer containing `\0` | pass lengths around, use `memchr` |
| Next request starts at the wrong offset | `consumed` computed before the last `\r\n\r\n` | include all 4 bytes of the terminator |
| Reply JSON is broken for some values | value inserted without escaping | escape `"` and `\` |
| Crash after `hm_remove` | pointer from `hm_get` used after removal | copy what you need, or finish using it first |
| Leak after a failed `hm_put` | caller assumed the map took ownership | free the value when `hm_put` returns -1 |
| Map loses entries after growth | rehash used the old slot count | recompute `hash & (new_n - 1)` for every node |

---

# Summary

- PIKA is request line, headers, blank line, with JSON responses and a close after one response.
- Strict parsing means one reading only: exact spaces, strict `\r\n`, case rules, no duplicates, no unknown headers, NUL rejected.
- Every length has a cap, and requests have no body, so the first blank line always ends the request.
- The parser is a pure function of `(buf, len)` and can return `PARSE_INCOMPLETE` for the socket stage.
- A key maps to a list of values; `append` extends the list, a plain `SET` replaces it.
- The store is two layers of one generic hash map (chaining, FNV-1a, doubling at 0.75 load factor) with clear ownership rules.
- Listing needs a stable order, so sort the keys before applying `after` and `limit`.
- Values come from clients, so escape them when writing JSON.