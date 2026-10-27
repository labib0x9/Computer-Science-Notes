---
title: "Strings in C Programming"
date: 2026-10-27
tags:
  - c
  - string
  - memory
---

# Overview

A C string is an array of `char` that ends with a zero byte, `'\0'` (the **null terminator**). C has no real string type, so you manage the memory and the length yourself. Most string bugs in C (crashes, buffer overflows) come from forgetting that.

| # | Section | What it covers |
| --- | --- | --- |
| 1 | [Memory layout](#memory-layout) | `'\0'`, `strlen` vs `sizeof` |
| 2 | [Array vs literal](#array-vs-literal) | Writable copy vs read-only data |
| 3 | [Comparing](#comparing) | Why `==` is wrong |
| 4 | [Copying and buffer size](#copying-and-buffer-size) | `strcpy`, `snprintf`, `strncpy` |
| 5 | [Heap strings](#heap-strings) | `malloc`, `strdup`, `free` |
| 6 | [Reading input](#reading-input) | `fgets`, why not `gets` |
| 7 | [String overflow](#string-overflow) | The bug and the fix |
| 8 | [Library functions](#library-functions) | `<string.h>` cheat sheet |
| 9 | [UTF-8](#utf-8) | Bytes vs characters |
| 10 | [Common mistakes](#common-mistakes) | Bugs and how to avoid them |
| 11 | [Exercises](#exercises) | Practice |

---

# Memory layout

```c
char s[] = "hi";   // 3 bytes in memory: 'h', 'i', '\0'
```

```text
index:   0     1     2
value:  'h'   'i'   '\0'
byte:   0x68  0x69  0x00
```

- All string functions find the end by scanning for `'\0'`. If it is missing, they keep reading past the array (undefined behavior).
- `"hi"` takes 3 bytes, but `strlen("hi")` is 2, because `strlen` does not count the terminator.
- `'\0'` (a char with value 0) is not the same as `NULL` (a null pointer) and not the same as `'0'` (the digit character, value 48).

---

# Array vs literal

There are two ways to declare a string, and they behave differently.

```c
#include <stdio.h>
#include <string.h>

int main(void) {
    /* array: your own writable copy */
    char s[] = "Hello";
    s[0] = 'X';
    printf("%s\n", s);                // Xello
    printf("%zu\n", sizeof(s));       // 6 = 5 chars + '\0'
    printf("%zu\n", strlen(s));       // 5

    /* pointer to a string literal: read-only memory */
    const char *ss = "vauvau";
    // ss[0] = 'V';                   // crash: writing to read-only memory
    printf("%s\n", ss);               // vauvau
    printf("%zu\n", sizeof(ss));      // 8 = size of a pointer (64-bit)
    printf("%zu\n", sizeof(*ss));     // 1 = one char
    printf("%zu\n", strlen(ss));      // 6

    ss = "but i can!";                // the pointer can point elsewhere
    printf("%s\n", ss);

    ss = s;                           // or at the writable array
    printf("%s\n", ss);
    return 0;
}
```

| | `char s[] = "Hello";` | `const char *p = "Hello";` |
| --- | --- | --- |
| What it is | array of 6 chars | pointer (8 bytes) to a literal |
| Where the chars live | stack if local, data segment if global | read-only data (rodata) |
| Can modify chars | yes | no, crashes |
| Can reassign | no (array name is fixed) | yes, pointer can change |
| `sizeof` | 6 | 8 |

Notes:
- Writing to a literal is undefined behavior. On Linux it usually gives a segmentation fault (SIGSEGV), and on macOS often a bus error (SIGBUS).
- Always write `const char *` for literals, so the compiler catches the mistake for you.
- `sizeof` on an array gives its size in bytes. A function parameter like `char s[]` is really a pointer, so `sizeof` there gives 8. Pass the length separately.

---

# Comparing

```c
char a[] = "abc";
char b[] = "abc";

if (a == b) { }                    // wrong: compares addresses, false here
if (strcmp(a, b) == 0) { }         // right: compares contents
```

`strcmp` returns 0 if equal, a negative number if the first string is smaller, and a positive number if it is larger. Don't rely on the exact value, only on its sign. `strncmp(a, b, n)` compares at most `n` bytes.

---

# Copying and buffer size

A string of length `n` needs `n + 1` bytes.

```c
char dst[4];
strcpy(dst, "hello");      // BUG: needs 6 bytes, writes past the array
```

`strcpy` and `strcat` don't check the space, so long input overflows the buffer. That is a **buffer overflow**, a classic security bug.

Safe pattern:

```c
char dst[16];
int n = snprintf(dst, sizeof dst, "%s", src);
if (n < 0 || (size_t)n >= sizeof dst) {
    // src did not fit; dst holds a truncated, still valid string
}
```

`snprintf` never writes past the size you give it, always adds `'\0'`, and returns the length the full string would have had, so you can detect truncation.

`strncpy` looks safe but is a trap:

```c
char dst[4];
strncpy(dst, "hello", sizeof dst);   // dst = 'h','e','l','l' with NO '\0'
printf("%s\n", dst);                 // reads past the array
```

If the source is as long as or longer than `n`, `strncpy` does not write the terminator. If it is shorter, it pads the rest with zeros, which wastes time. Prefer `snprintf`, or set `dst[sizeof dst - 1] = '\0'` yourself.

---

# Heap strings

```c
char *copy = malloc(strlen(src) + 1);   // +1 for '\0'
if (copy != NULL) {
    strcpy(copy, src);                  // safe here: we sized it exactly
    // use copy
    free(copy);
}
```

`strdup(src)` does the same in one call. It is standard in C23 and available on Linux and macOS, but older strict C standards don't include it. You still must `free` the result.

Forgetting the `+ 1` is a very common bug: it writes one byte past the allocation.

---

# Reading input

```c
char buf[1024];
if (fgets(buf, sizeof buf, stdin) != NULL) {
    buf[strcspn(buf, "\n")] = '\0';   // strip the trailing newline, if any
}
```

- `fgets` reads at most `size - 1` characters and always adds `'\0'`. It returns `NULL` on end-of-file or error.
- It keeps the `'\n'`, which is why we strip it. `strcspn(buf, "\n")` returns the index of the first newline, or the string length if there is none, so it is safe either way.
- Never use `gets`. It has no size limit and cannot be used safely, so it was removed from the language in C11.

---

# String overflow

## Bug 1: `gets`

```c
char buffer[3];
gets(buffer);                 // typing "hello" writes 6 bytes into 3
printf("%s\n", buffer);
```

Undefined behavior. Modern compilers warn, and some systems refuse to link it.

## Bug 2: fixed-size destination

```c
char buf[1024];
fgets(buf, sizeof buf, stdin);
buf[strcspn(buf, "\n")] = '\0';

char newStr[3];
strcpy(newStr, buf);          // BUG: overflows if buf is longer than 2 chars
printf("%s\n", newStr);
```

Two things in your original version were also wrong:
- `fgets(&buf, ...)` should be `fgets(buf, ...)`. `&buf` has the wrong type (pointer to the whole array).
- `buffer[strcspn(buf, "\n")]` should be `buf[...]`. It used the wrong array.

## Fix

```c
char small[3];
snprintf(small, sizeof small, "%s", buf);   // keeps at most 2 chars + '\0'
printf("%s\n", small);
```

Try it: compile with AddressSanitizer and run the buggy version to see how it reports the error.

```bash
gcc -Wall -Wextra -g -fsanitize=address,undefined main.c -o a.out
```

---

# Library functions

All in `<string.h>` unless noted.

| Function | Does | Watch out |
| --- | --- | --- |
| `strlen(s)` | length, without `'\0'` | scans the whole string each call, so don't call it in a loop condition |
| `strcmp(a, b)` | compare contents | 0 means equal |
| `strchr(s, c)` | first occurrence of a char | returns `NULL` if not found |
| `strstr(s, sub)` | first occurrence of a substring | returns `NULL` if not found |
| `strcat(d, s)` | append | no bounds check |
| `memcpy(d, s, n)` | copy `n` raw bytes | ranges must not overlap |
| `memmove(d, s, n)` | same, overlap is safe | |
| `memset(p, c, n)` | fill `n` bytes | |
| `snprintf(d, n, fmt, ...)` | formatted, bounded write | check the return value |
| `toupper`, `isalpha`, `isdigit` | per-char helpers (`<ctype.h>`) | |
| `strtol`, `strtod` | text to number (`<stdlib.h>`) | |

## `strtok` (splitting)

`strtok` splits a string into tokens, but it **modifies the string** by writing `'\0'` into it, and it keeps hidden state between calls.

```c
char line[] = "a,b,c";            // must be a writable array, not a literal
for (char *t = strtok(line, ","); t != NULL; t = strtok(NULL, ",")) {
    printf("%s\n", t);
}
```

Calling it on a literal crashes, and it is not thread-safe. This matters for the `Split` exercise below, where you may prefer to write your own.

## Text to number

Prefer `strtol` over `atoi`, because `atoi` can't tell you whether the conversion failed.

```c
char *end;
long v = strtol("123abc", &end, 10);   // v = 123, end points at "abc"
if (end == "123abc") { /* no digits were read */ }
```

---

# UTF-8

C strings are bytes, not characters. In UTF-8 an English letter is 1 byte, but a Bangla letter is 3 bytes and an emoji is 4.

```c
const char *bn = "বাংলা";
printf("%zu\n", strlen(bn));   // number of bytes (15), not 5 characters
```

So `strlen`, indexing and `SubString` count bytes. Cutting a string in the middle of a multi-byte character produces invalid text. Keep this in mind for your string builder: decide whether your type is "bytes" or "characters" and document it.

---

# Common mistakes

| Bug | What happens | How to avoid |
| --- | --- | --- |
| Missing `'\0'` | functions read past the end | use `snprintf`, check `strncpy` |
| `malloc(strlen(s))` | one byte too small | `strlen(s) + 1` |
| `==` on strings | compares addresses | `strcmp` |
| Writing to a literal | crash | use `char[]`, or `const char *` to catch it |
| `strcpy` into a small buffer | buffer overflow | `snprintf` with `sizeof` |
| Using `gets` | unbounded read | `fgets` |
| `sizeof` on a parameter | gives 8, not the length | pass the length |
| `strtok` on a literal | crash | copy into an array first |
| Returning a local array | dangling pointer | return a `malloc`'d string or let the caller pass a buffer |

---

# Exercises

## Exercise 01: string builder

Build a small string type that owns its memory.

```c
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

typedef struct String {
    size_t len;     // length without '\0'
    char *data;     // heap memory, always '\0' terminated
} string_t;

string_t NewString(const char *str);                    // copy str to the heap
bool AppendString(string_t *a, const string_t *b);      // a = a + b
string_t SubString(const string_t *a, size_t start, size_t len);
// stretch: Split, returning an array of string_t

void FreeString(string_t *str) {
    if (str->data) free(str->data);
    str->data = NULL;
    str->len = 0;
}
```

Things to decide (and test) yourself:
- What should `NewString(NULL)` return? An empty string with `data == NULL`, or `data == ""` on the heap? Pick one and keep it consistent.
- `AppendString` needs `realloc`. If `realloc` fails it returns `NULL` and the old block is still valid, so don't overwrite `a->data` until you know it worked.
- Keep the `'\0'` at `data[len]` after every operation so you can still pass `data` to `printf`.
- `SubString`: what if `start + len` goes past the end? Clamp it, or return an error?
- `Split` needs to return several strings. Think about who frees them.
- Test with an empty string, a one-character string, and appending a string to itself.

The original `SubString(a, len)` had no start index. A substring needs a start position, so I added one.

## Exercise 02: write your own libc

Implement `my_strlen`, `my_strcpy`, `my_strcmp` and `my_strchr` using only pointer arithmetic, with no indexing and no library calls.

## Exercise 03: reverse and palindrome

Reverse a string in place with two pointers, then use it to check for a palindrome.

## Exercise 04: safe line reader

Write `int read_line(char *buf, size_t size)` that wraps `fgets`, strips the newline, and returns an error if the line was longer than the buffer. Hint: if there is no `'\n'` in the buffer and you did not hit end-of-file, the line was cut off.

---

# Summary

- A string is `char` bytes ending in `'\0'`. A string of length `n` needs `n + 1` bytes.
- `char s[] = "..."` is a writable copy. `const char *p = "..."` points to read-only memory.
- Compare with `strcmp`, never `==`.
- Use `snprintf` and `fgets` with `sizeof`. Avoid `gets`, `strcpy` into fixed buffers, and `strncpy`.
- Strings are bytes, so UTF-8 text needs care.
- Compile with `-fsanitize=address` while learning.