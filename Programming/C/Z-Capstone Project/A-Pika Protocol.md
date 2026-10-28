---
title: "Capstone Project: PIKA Protocol Defination"
date: 2026-10-28
tags:
  - c
  - capstone
---

# The protocol (PIKA/1.0)

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
| `GET` | read one key, or list keys with `key: *` | `GET <bucket> PIKA/1.0` | `key` (required), `limit` |
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
