---
title: 8086 Microprocessor Assembly on Mac M1
date: 2026-10-24
tags:
  - assembly
  - microprocessor
summary: Environment setup for 8086 microprocessor assembly on a Mac M1 using DOSBox
---

## Mac (M1) Setup

- Install DOSBox Staging: https://www.dosbox-staging.org/

Create a folder for your assembly files. Run this in the macOS terminal:

```bash
mkdir -p ~/dos/asm
```

Mount that folder as the C drive. Run this inside DOSBox:

```
mount c ~/dos/asm
c:
```

- Get `jwasm.exe` and `link.exe` (a MASM-compatible assembler and a linker) and put them in the folder.

Assemble and run your program inside DOSBox:

```
jwasm -mz hello.asm
hello
```

![Running an ASM program in DOSBox](run-asm.gif)

## VS Code

- Extension: ASM Code Lens