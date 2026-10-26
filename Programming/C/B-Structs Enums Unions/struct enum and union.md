---
title: Struct, Enum and Union in C Programming
date: 2026-10-24
tags:
  - c
  - struct
  - enum
  - union
---

# Overview

Sometimes we need to define our own data types. Structs, enums and unions are the tools C gives us for that.

| # | Section |
| --- | --- |
| 1 | [Struct](#struct) |
| 2 | [Enum](#enum) |
| 3 | [Union](#union) |

---

## Struct

A struct combines several values, possibly of different types, into one data type.

```c
typedef struct {
    void* ptr;
    int len, cap;
} slice;

typedef struct Node {
    struct Node* next;
    int value;
} node_t;
```

Here we declared a new data type named `slice`, which stores a pointer and two ints. A struct that points to itself, like `Node`, must use `struct Node*` inside its own definition, because the `node_t` alias does not exist yet at that point.

You can access the fields in two ways.

### Option 1: variable (use `.`)

```c
slice s;
s.len = 10;
```

### Option 2: pointer (use `->`)

```c
slice v;
slice* s = &v;
s->len = 10;   // same as (*s).len = 10;
```

### Struct Size

You can find the size of a struct with `sizeof`.

```c
typedef struct {
    bool alloc;
    long long c;
    int a;
} a_t;

printf("A = %zu\n", sizeof(a_t));
```

The size can change if you change the order of the fields, because of padding and alignment. The compiler adds unused bytes between fields so that each field starts at an address that suits its type. Try this:

```c
#include <stdio.h>
#include <stdbool.h>

typedef struct {
    bool alloc;      // 1 byte + 7 padding
    long long c;     // 8 bytes
    int a;           // 4 bytes + 4 padding at the end
} a_t;

typedef struct {
    long long c;     // 8 bytes
    int a;           // 4 bytes
    bool alloc;      // 1 byte + 3 padding at the end
} b_t;

int main(void) {
    printf("A = %zu\n", sizeof(a_t));
    printf("B = %zu\n", sizeof(b_t));
    return 0;
}
```

If you run it, you will see two different sizes, because of padding and alignment:

```text
A = 24
B = 16
```

Total size must be a multipler of size of the biggest field. Thus compiler adds padding at the end.

You can also find the offset of a field, which is its distance in bytes from the start of the struct. `offsetof` needs `<stddef.h>`.

```c
#include <stddef.h>

printf("%zu\n", offsetof(a_t, c));   // 8
```

---

## Enum

An enum gives names to a set of integer constants. Say you store the day of the week as a number: how would you know whether `2` means Sunday or Monday? You would need documentation to understand it. With an enum, you write `SUNDAY` or `MONDAY` instead, and the code explains itself.

```c
typedef enum DaysOfWeek {
    FRIDAY,     // 0
    SATURDAY,   // 1
    SUNDAY,     // 2
    MONDAY,     // 3
} days_of_week_t;

days_of_week_t day = FRIDAY;
```

The values start from 0 and go up by one, unless you set them yourself, for example `FRIDAY = 5`.

---

## Union

A union looks like a struct, but it works differently. In a struct, every field has its own memory, and the size is the sum of all fields plus padding. In a union, all fields share the same memory, and the size is the size of the largest field.

```c
typedef struct {
    int a;          // 4 bytes
    long long b;    // 8 bytes
} a_t;

typedef union {
    int a;          // 4 bytes
    long long b;    // 8 bytes
} b_t;

printf("%zu\n", sizeof(a_t));   // 16 bytes = 4 (int) + 4 (padding) + 8 (long long)
printf("%zu\n", sizeof(b_t));   // 8 bytes = max(4 (int), 8 (long long))
```

In a union, the same memory can be read as an `int` or as a `long long`. Only one field holds a meaningful value at a time: writing to one field overwrites the others.

## Source
[Source Code]()