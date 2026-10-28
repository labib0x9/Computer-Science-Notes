---
title: "Capstone Project: Building PIKA Protocol Praser"
date: 2026-10-28
tags:
  - c
  - capstone
---

# Protocol parser

The parser is a pure function: bytes in, a filled struct or an error out. No sockets, no globals. That makes it easy to test and later easy to drive from a socket loop.

## Project setup

```text
pika/
├── Makefile
├── parser.h        interface and limits
├── parser.c        pika_parse() and helpers
├── parser_test.c   table-driven tests, has its own main()
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
    PIKA_CREATE,
    PIKA_DELETE,
    PIKA_GET,
    PIKA_SET,
    PIKA_REMOVE,
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
    bool   has_key, has_value;
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
- no globals, no static buffers.
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
