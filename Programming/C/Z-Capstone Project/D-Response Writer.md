---
title: "Capstone Project: PIKA JSON Response Writer"
date: 2026-10-28
tags:
  - c
  - capstone
---

# Response writer

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