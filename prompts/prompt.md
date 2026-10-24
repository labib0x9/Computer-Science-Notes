You are working inside my personal programming notes repository.

The repository contains code organized into folders by language and topic. Examples include:

* Go
* C
* C++
* Assembly
* JavaScript
* TypeScript
* Backend
* Systems
* Linux
* Networking
* Databases
* Operating Systems
* Compilers
* other technical topics

The exact directory structure is not fixed. You must inspect the repository first.

## OBJECTIVE

For **every folder that contains at least one source-code file**, create a Markdown documentation file that explains the code contained in that folder.

The Markdown file should be generated from the actual code in that folder.

Do NOT create generic tutorials unrelated to the repository.

The repository's code is the primary source of truth.

---

# 1. DISCOVER THE REPOSITORY

First recursively inspect the repository.

Determine:

* directory structure
* programming languages
* topics
* source files
* existing README/Markdown files
* configuration files
* scripts
* test files
* examples
* dependencies

Do not modify existing source files.

Do not delete anything.

Do not move anything.

Do not rename anything.

Before generating documentation, determine which folders qualify as documentation targets.

A folder qualifies if it contains one or more meaningful code files directly inside it.

Examples:

```text
go/
├── basics/
│   ├── variables.go
│   ├── loops.go
│   └── functions.go
```

Generate:

```text
go/basics/README.md
```

Another example:

```text
systems/
├── epoll/
│   ├── server.c
│   └── client.c
```

Generate:

```text
systems/epoll/README.md
```

---

# 2. DO NOT DOCUMENT EVERY DIRECTORY BLINDLY

Do NOT create a README merely because a directory exists.

For example:

```text
go/
go/basics/
go/concurrency/
```

If `go/` contains no code directly, do not automatically create:

```text
go/README.md
```

unless there is already a meaningful existing README that should be preserved.

The target is:

> One generated Markdown file for each folder containing code.

---

# 3. PRESERVE EXISTING DOCUMENTATION

Before generating anything:

* inspect existing `.md` files
* do not overwrite manually written documentation blindly
* if a README already exists, determine whether it appears generated or manually written
* preserve useful existing content

If an existing README is clearly manually maintained, either:

1. update it carefully while preserving its content, or
2. create a separate generated documentation file if that is safer.

Never destroy existing knowledge.

---

# 4. IDENTIFY THE CODE

For every target folder, inspect all relevant source files.

Recognize common extensions including:

### C

```text
.c
.h
```

### C++

```text
.cpp
.cc
.cxx
.hpp
.h
```

### Go

```text
.go
```

### Assembly

```text
.s
.S
.asm
```

### JavaScript

```text
.js
.mjs
.cjs
```

### TypeScript

```text
.ts
.mts
.cts
```

### Python

```text
.py
```

### Rust

```text
.rs
```

and other source languages if they appear in the repository.

Do not assume the repository only contains the languages listed above.

---

# 5. UNDERSTAND THE CODE BEFORE DOCUMENTING IT

For each folder:

1. Read every relevant code file.
2. Determine what the folder is teaching or demonstrating.
3. Identify the concepts used.
4. Understand how the files relate to each other.
5. Identify important implementation details.
6. Identify dependencies between files.
7. Identify interesting edge cases.
8. Identify potentially dangerous misconceptions.
9. Determine whether the code demonstrates a language feature, algorithm, system concept, architecture pattern, API, or experiment.

Do not merely summarize filenames.

---

# 6. GENERATED README STRUCTURE

Each generated Markdown file should generally follow this structure:

````md
# <Folder Topic>

## Overview

Short explanation of what this folder contains and what concept it demonstrates.

## Learning Objectives

After studying this folder, I should understand:

- ...
- ...
- ...

## Files

| File | Purpose |
|------|---------|
| `example.go` | ... |
| `server.go` | ... |

## Concepts

### 1. <Concept>

Explanation.

Relevant code:

```go
...
````

Explain what is happening.

### 2. <Concept>

...

## How the Code Works

Explain the execution flow.

If multiple files interact:

```text
file A
   ↓
function X
   ↓
file B
   ↓
function Y
```

## Important Details

Explain subtle or non-obvious behavior.

## Common Mistakes

List mistakes someone studying this code might make.

## Language Notes

If applicable, explain language-specific syntax or semantics.

## Related Concepts

* ...
* ...
* ...

## Questions to Test Myself

1. ...
2. ...
3. ...

## Further Experiments

Suggest small experiments that can be performed by modifying the existing code.

## Source Files

* `file1.go`
* `file2.go`

````

Adapt the structure when necessary.

Do NOT force irrelevant sections onto every folder.

---

# 7. CODE EXPLANATIONS MUST USE THE ACTUAL CODE

For example, if the folder contains:

```go
func worker(ch <-chan int) {
    for value := range ch {
        fmt.Println(value)
    }
}
````

The documentation should explain the actual code:

* receive-only channel
* `range` over channel
* blocking behavior
* channel closure
* goroutine interaction if present

Do NOT replace it with a generic explanation of channels.

The documentation should answer:

> "What is this code doing and why?"

---

# 8. LANGUAGE-AWARE DOCUMENTATION

Adapt explanations to the language.

For Go, pay attention to:

* slices
* maps
* structs
* interfaces
* methods
* pointers
* goroutines
* channels
* context
* interfaces
* error handling
* generics
* defer
* embedding
* package structure

For C:

* pointers
* memory
* stack/heap
* structs
* function pointers
* preprocessing
* compilation
* system calls
* undefined behavior
* manual memory management

For C++:

* RAII
* references
* pointers
* classes
* inheritance
* templates
* STL
* move semantics
* smart pointers
* concurrency
* object lifetime

For Assembly:

* registers
* instructions
* stack
* calling convention
* memory addressing
* control flow
* ABI
* interaction with C/C++

For JavaScript:

* objects
* arrays
* closures
* callbacks
* promises
* async/await
* event loop
* modules
* prototypes
* `this`

For TypeScript:

* types
* interfaces
* generics
* unions
* narrowing
* type assertions
* structural typing
* runtime vs compile-time behavior

For backend folders:

* request flow
* handlers
* services
* repositories
* databases
* transactions
* caching
* queues
* authentication
* concurrency
* error handling

For systems folders:

* processes
* threads
* memory
* syscalls
* files
* networking
* scheduling
* IPC
* kernel interaction

For networking folders:

* packets
* sockets
* TCP/UDP
* DNS
* HTTP
* routing
* namespaces
* interfaces
* protocols

Do not blindly include all of these. Include only concepts actually demonstrated by the code.

---

# 9. CROSS-LANGUAGE COMPARISONS

When useful, relate concepts to other languages present in the repository.

For example, if documenting a Go folder about slices:

````md
## Comparison

### Go

```go
s := []int{1, 2, 3}
````

### C++

```cpp
std::vector<int> v{1, 2, 3};
```

### C

```c
int values[] = {1, 2, 3};
```

Explain the important semantic differences.

````

Do this only when it genuinely improves understanding.

Do not create comparisons merely for the sake of comparison.

---

# 10. SYSTEM-LEVEL CODE

For low-level code, explain what happens underneath the source code whenever the code makes that relevant.

For example:

```c
read(fd, buffer, size);
````

Explain:

```text
user-space program
       ↓
read()
       ↓
system call boundary
       ↓
kernel
       ↓
file/socket/device
       ↓
kernel
       ↓
user-space buffer
```

For assembly, explain:

```text
source code
    ↓
compiler
    ↓
assembly
    ↓
assembler
    ↓
object file
    ↓
linker
    ↓
executable
```

Only include such explanations when relevant to the actual code.

---

# 11. EXECUTION FLOW

If the code has meaningful execution flow, document it.

For example:

```text
main()
  │
  ├── create listener
  │
  ├── accept connection
  │
  ├── spawn worker
  │
  └── worker()
        │
        ├── read
        ├── process
        └── write
```

Use ASCII diagrams rather than generating images.

---

# 12. COMPILATION / EXECUTION

Where relevant, include commands to build and run the code.

Examples:

### Go

```bash
go run .
```

### C

```bash
gcc main.c -o main
./main
```

### C++

```bash
g++ main.cpp -o main
./main
```

### Assembly

Use the appropriate assembler/linker commands based on the actual architecture and build system.

Do not invent commands.

If the repository contains a Makefile, CMake configuration, Go module, Dockerfile, etc., use the repository's actual build mechanism.

---

# 13. EXPLAIN UNUSUAL CODE

If the code contains something non-obvious, explain it.

Examples:

```go
defer func() {
    ...
}()
```

or:

```c
*(ptr + offset)
```

or:

```cpp
std::move(x)
```

or:

```asm
mov %rsp, %rbp
```

or:

```ts
const result = await Promise.all(...)
```

Explain why it works, not merely what syntax means.

---

# 14. DO NOT HIDE BUGS

If the code contains:

* a bug
* race condition
* memory leak
* undefined behavior
* incorrect error handling
* unsafe assumption
* deadlock possibility
* security issue
* inefficient algorithm
* portability issue

document it explicitly.

Use:

```md
## Issues / Risks

> This code has a potential race condition because ...

```

Do not silently "fix" the source code.

Documentation should describe the repository as it actually exists.

---

# 15. DO NOT FABRICATE INTENT

If you cannot determine why something exists, say:

> "The code appears to demonstrate X, based on ..."

Do not invent explanations about the author's intention.

---

# 16. LEARNING LEVEL

Assume the reader is an experienced programming student who understands:

* C
* C++
* Go
* algorithms
* data structures
* basic operating systems
* networking
* databases

Therefore:

DO NOT explain:

> "A variable stores a value."

Instead explain the language-specific behavior.

For example:

> "Unlike a Go slice, JavaScript's Array is a dynamic object with indexed properties and a collection of methods such as `map`, `filter`, and `reduce`."

Prioritize semantic differences and implementation details.

---

# 17. CREATE A TOPIC INDEX

After generating the individual Markdown files, create or update a root-level:

```text
KNOWLEDGE_INDEX.md
```

It should contain:

```md
# Knowledge Index

## Languages

### Go

- [Basics](...)
- [Concurrency](...)
- [Interfaces](...)

### C

- [Pointers](...)
- [Memory](...)

### C++

...

### Assembly

...

## Backend

...

## Systems

...

## Networking

...

## Databases

...
```

Only include directories/files that actually exist.

---

# 18. CONCEPT INDEX

Also create:

```text
CONCEPT_INDEX.md
```

Map concepts to the folders where they appear.

Example:

```md
# Concept Index

## Concurrency

- Go: `go/concurrency/`
- Backend: `backend/worker-pool/`
- C++: `cpp/threading/`

## Memory

- C: `c/pointers/`
- C++: `cpp/memory/`
- Assembly: `asm/stack/`

## Networking

- C: `systems/socket/`
- Go: `go/networking/`
- Backend: `backend/http/`
```

This should be generated from the repository structure.

---

# 19. DO NOT DUPLICATE LARGE AMOUNTS OF CODE

Documentation should explain the source code.

Do not copy entire source files into Markdown.

Use small relevant excerpts.

For example:

```go
ch := make(chan int)

go worker(ch)
```

Then explain it.

Only reproduce complete code when the file itself is very small and the complete example is necessary.

---

# 20. PRESERVE CODE REFERENCES

Whenever discussing code, mention the actual source file.

Example:

```md
The worker is created in `worker.go` and consumes values from the channel.
```

For functions:

```md
`worker()` in `worker.go`
```

For important lines, include line numbers when practical.

---

# 21. GENERATED DOCUMENTATION MARKER

Every generated Markdown file should contain a small marker near the top:

```md
<!-- AUTO-GENERATED FROM SOURCE CODE -->
```

Do not put this marker into manually written files unless you are intentionally converting them to generated documentation.

---

# 22. IDEMPOTENCY

The process must be safe to run repeatedly.

Running the documentation generator twice should not produce:

```text
README.md
README-1.md
README-2.md
README-final.md
```

Instead, regenerate/update the generated documentation deterministically.

Do not duplicate sections.

Do not duplicate entries in indexes.

---

# 23. IMPORTANT: DOCUMENTATION IS DERIVED FROM CODE

The hierarchy is:

```text
SOURCE CODE
     ↓
ANALYSIS
     ↓
DOCUMENTATION
     ↓
INDEX
```

Never modify source code merely to make documentation easier.

If the code is confusing, explain the confusion.

---

# 24. FINAL REPORT

After completing the generation, provide a summary:

```text
Documentation generated.

Folders analyzed: X
Folders documented: X
Languages detected: ...

Generated:
- ...
- ...
- ...

Skipped:
- ...

Existing documentation preserved:
- ...

Potential issues found:
- ...
```

Also report any folders that could not be confidently documented.

---

# 25. CRITICAL CONSTRAINT

Do not turn this repository into a generic tutorial website.

This is my **personal code-based knowledge base**.

The generated Markdown should answer:

> "What did I learn from the code in this folder?"

rather than:

> "What does a generic programming textbook say about this topic?"

Use my actual code, terminology, experiments, and implementation details as the center of every document.

Start by inspecting the entire repository and its directory structure. Then generate the documentation systematically.
