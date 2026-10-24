---
title: "Smart Pointers and Memory Ownership in C++"
date: 2026-10-25
tags:
  - cpp
  - memory-management
  - smart-pointers
  - code-explained
summary: "Demonstrates std::unique_ptr and std::shared_ptr for automatic memory management and RAII in C++."
---

## Overview

This program demonstrates modern C++ dynamic memory management using smart pointers from the `<memory>` header. It showcases single-owner exclusive semantics with `std::unique_ptr` and shared reference-counted ownership with `std::shared_ptr`. C++ developers use smart pointers to adhere to RAII (Resource Acquisition Is Initialization), eliminating manual `delete` calls and avoiding memory leaks or dangling pointers.

## Source

https://github.com/labib0x9/Programming-Notes/blob/main/Programming/C%2B%2B/pointer%20and%20dma/smart_pointer.cpp

## Concepts used

- RAII (Resource Acquisition Is Initialization)
- Exclusive ownership with `std::unique_ptr`
- Ownership transfer via `std::move`
- Dynamic reallocation with `reset()`
- Shared ownership and reference counting with `std::shared_ptr`
- Deterministic destructor execution at block scope boundaries

## How it works

The program includes `<iostream>` for output and `<memory>` for smart pointer templates, bringing the `std` namespace into scope.

A struct `Tree` demonstrates lifecycle tracking. Its constructor prints `"construction"`, its destructor prints `"destruction"`, and its member method `Print()` outputs `"Hello unique pointer"`.

The function `Get(int n = 0)` returns a `std::unique_ptr<int>`. Inside, it allocates an integer on the heap, logs the pointer address to the console, and returns `v`. Because `std::unique_ptr` cannot be copied, modern C++ moves the pointer into the caller's receiving variable upon return.

The `main` function demonstrates various smart pointer operations:
- Allocates an integer on the heap managed by `unique_ptr<int> value(new int(15))`, printing its pointer address and dereferenced value (`*value`).
- Calls `value.reset(new int(42))`. The `reset` method deallocates the existing heap integer (`15`) and takes ownership of the newly allocated integer (`42`).
- Modifies the stored value in place via dereference assignment: `*value = 35`.
- Creates a `unique_ptr<string> strPtr` managing a dynamically allocated `std::string` and populates it using `strPtr->assign("Hello World")`.
- Demonstrates move semantics: `unique_ptr<int> valueNew(std::move(value))` transfers heap ownership from `value` to `valueNew`. After the move, `value` becomes `nullptr`, confirmed by evaluating `value == nullptr`.
- Calls `Get()`, assigning the returned unique pointer to `vPtr` and logging its address in `main()`.
- Enters a local block scope `{ ... }` declaring `unique_ptr<Tree> newTree(new Tree())` and invokes `newTree->Print()`. When the execution exits the block, `newTree` goes out of scope and its destructor executes automatically.
- Demonstrates `std::shared_ptr`: declares `shared_ptr<int> sValue(new int(19))` and copies it into `nNewValue`, incrementing the internal reference count.
- Enters nested block scopes for `shared_ptr<Tree>`. The outer variable `sTree` takes an assignment from the inner variable `st`. When the inner scope ends, the `Tree` object remains alive because `sTree` still holds a reference. When the outer scope ends, the reference count drops to zero and the destructor executes.
- Returns `0`.

## Key excerpt

```cpp
    // no copy, but you can move
    unique_ptr<int> valueNew(std::move(value)); // transfer ownership to valueNow, value becomes nullptr
    cout << valueNew << " " << *valueNew << endl;
    // cout << value << " " << *value << endl;
    cout << (value == nullptr ? "true" : "false") << endl;

    // what if, the value escapes the scope ??
    unique_ptr<int> vPtr = Get();
    cout << "main(): " << vPtr << endl;

    {
        unique_ptr<Tree> newTree(new Tree());
        // unique_ptr<Tree> newTree = make_unique<Tree>(); // something is wrong here ?
        newTree->Print();
    }
```

## Expected behavior (from reading the code, not run)

To compile and execute the program using a C++11 compliant compiler:
```bash
g++ -std=c++11 smart_pointer.cpp -o smart_pointer && ./smart_pointer
```

Expected behavior:
The program outputs addresses and values as memory is allocated and transferred. Pointer addresses (represented as `<hex_address>`) will differ per run:
```text
<hex_address> 15
<hex_address> 42
<hex_address> 35
<hex_address> Hello World
<hex_address> 35
true
Get(): <hex_address>
main(): <hex_address>
construction
Hello unique pointer
destruction
Inner scope: construction
Outer scope: destruction
```

## Things to watch

- **Direct use of `new`:** The code passes raw `new` expressions to constructors (e.g. `unique_ptr<int>(new int(15))`). In modern C++ (C++14 for `std::make_unique` and C++11 for `std::make_shared`), factory functions are preferred for exception safety and single-allocation optimization for shared pointers.
- **Pointer address variability:** Directly printing smart pointers (`cout << value`) displays virtual memory addresses, which vary across operating systems and execution runs.
- **Language standard version:** `std::unique_ptr`, `std::shared_ptr`, and `std::move` require C++11 or newer; compilation will fail under older C++98/C++03 standards.

## Terms

- **Smart Pointer:** A class template wrapping a raw pointer that provides automatic lifetime management and deallocates memory when no longer needed.
- **std::unique_ptr:** A smart pointer that retains sole ownership of an object through move-only semantics.
- **std::shared_ptr:** A smart pointer that manages an object's lifetime through reference counting, freeing the resource when the last owner is destroyed.
- **RAII:** Resource Acquisition Is Initialization, a C++ idiom tying the acquisition and release of resources to object construction and destruction.

## Related

There are no related notes in this directory.

*This explanation was generated by an AI from reading the source, not from running it. Verify before relying on it.*
