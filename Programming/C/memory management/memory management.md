---
title: AddressSanitizer (ASan) in C
date: 2026-10-25
tags:
  - c
  - gcc
  - debugging
summary: Compile and run a C program with AddressSanitizer to catch memory bugs
---

## Command

```bash
gcc -fsanitize=address -g a.c && ./a.out
```

- `-fsanitize=address` turns on AddressSanitizer (ASan). It adds checks around memory accesses, and the program stops with a report when it finds a bug.
- `-g` adds debug information, so the report shows file names and line numbers and not just addresses.
- `&&` runs `./a.out` only if the compile succeeded. The parentheses in the original are not needed.

## What it catches

- Buffer overflow on the heap, the stack and global variables
- Use after free
- Double free
- Memory leaks (see the note below)

## Tips

- Add `-fno-omit-frame-pointer` for cleaner stack traces in the report.
- ASan makes the program slower and use more memory, so use it for development and testing, not in production.
- Do not combine it with another sanitizer that has to be run separately, such as `-fsanitize=thread`. Use one at a time.