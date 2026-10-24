---
title: Memory Management in Go
date: 2026-10-24
tags:
  - golang
  - memory
summary: Escape analysis in Go and how to inspect it with the compiler
---

## Escape analysis

Escape analysis is the compiler's check of whether a variable can live on the stack or must be moved to the heap. To see its decisions, build with the `-m` flag:

```bash
go build -gcflags="-m" main.go
```