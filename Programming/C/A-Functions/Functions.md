---
title: Functions in C Programming
date: 2026-10-24
tags:
  - function
  - c
  - constructor
---

I am not going to explain the fundamentals and basic concepts here. If you are familiar with programming, you already know what a function is.

# Overview

| # | Section |
| --- | --- |
| 1 | [Constructor](#constructor) |
| 2 | [Destructor](#destructor)
| 2 | [Function Visibility](#function-visibility) |

## Constructor

You might have heard that programs always start from the `main()` function, right? That is not always true. With GCC and Clang, you can mark a function so that it runs before `main()`. This is called a constructor.

### Example 1

```c
#include <stdio.h> // needed for printf()

// runs before main()
__attribute__((constructor)) void init() {
    printf("Hello From init()\n");
}

int main() {
    printf("Hello From main()\n");
    return 0;
}
```

### How to run

```bash
```

### Example 2

```c
#include <stdio.h> // needed for printf()

#define public __attribute__((visibility("default")))
#define private __attribute__((visibility("hidden")))

private void hello_world() {
    printf("Hello World\n");
}

public void hello_world_public() {
    printf("Hello world public\n");
}

int main() {
    printf("Hello From main()\n");
    return 0;
}
```

### How to run

```bash
```

## Destructor

```c
#include <stdio.h> // needed for printf()

// runs before main()
__attribute__((constructor)) void init() {
    printf("Hello From init()\n");
}

// runs after main()
__attribute__((destructor)) void clean_up() {
    printf("Hello from cleanup()\n")
}

int main() {
    printf("Hello From main()\n");
    return 0;
}
```

## Function Visibility

In a shared library (a `.so` file), the `visibility` attribute controls which functions other programs can use.

- `default` ("public"): code outside the library can call the function.
- `hidden` ("private"): only code inside the same library can call it.

In a plain executable like Example 2, both functions behave the same, so the difference only shows up when you build a shared library.

---

## Source

[Source Code]()