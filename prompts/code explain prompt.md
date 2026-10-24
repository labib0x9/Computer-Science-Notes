You are writing documentation for a personal notes repository. Read this whole brief before doing anything. Work in phases and stop where the brief says to stop.

## 1. Goal

For every source code file with the extension .c, .cpp, .go or .asm in this repo, create ONE markdown (.md) file that explains that code. The .md file goes in the SAME folder as the code file. Each .md file is later rendered as a blog-style page by my notes website, so it must be accurate, readable and beginner-friendly.

## 2. File naming (exact)

- The md name is the code file's name with the extension replaced by `.md`. The name is otherwise unchanged: same capitalisation, underscores, spelling mistakes and spaces.
  - `create_a_container.go` becomes `create_a_container.md`
  - `hello_world.asm` becomes `hello_world.md`
  - `calc_test.go` becomes `calc_test.md`
- Collision rule. If the plain name would clash, use the longer form `<name>.<ext>.md` (for example `merge.c.md`). A clash happens when:
  1. another code file in the same folder has the same name without the extension (for example `merge.c` and `merge.go`; in that case BOTH get the long form: `merge.c.md` and `merge.go.md`), or
  2. an md file with that name already exists in that folder (for example `select.go` next to an existing `select.md`; the existing file is never touched).
- Do not rename, move, edit or delete any existing file. Do not touch any existing .md file, including every `readme.md`, `roadmap.md` and `README.md`.
- Only these extensions are in scope: .c, .cpp, .go, .asm. Ignore everything else: .h, .i, .o, .s, .exe, .gif, .mod, .sum, .sh, .html, .proto, .yml, .txt, .bak. If a .h file or other support file is needed to understand a code file, read it and mention it in the explanation, but do not create an md for it.
- If a code file is empty, or you cannot understand what it does, still create the md, but say so plainly in the summary (see section 5) instead of guessing.

## 3. Hard safety rules

1. NEVER run, compile, build or execute any code in this repo. Many files touch Docker, namespaces, chroot, mounts, processes, sockets, the filesystem, or the GitHub API (some can delete repositories). Only read files. Everything you write about behaviour comes from reading the code.
2. NEVER run `git commit`, `git push`, `git add` or anything that changes git state. I will review and commit myself. A push triggers my website to rebuild.
3. Secrets: if a file contains something that looks like a token, key, password or private key, do NOT copy it into the md. Write `<REDACTED>` in its place, and list the file in your final report.
4. Do not modify any code file.

## 4. Accuracy rules (important)

I check my notes carefully and dislike overstated claims. So:
- Explain only what is in the code. Do not invent features, history, or intent.
- Since you do not run the code, never write "output:" as if you had run it. Use the heading "Expected behavior (from reading the code, not run)".
- If something is uncertain, say "as far as I can tell" or "unclear from the code". Do not fill gaps with confident guesses.
- Mention version or platform dependence when it is real. Examples: Linux-only system calls (epoll, namespaces, `mmap` flags), POSIX-only APIs, behavior that changed between Go versions (such as loop variable semantics before and after Go 1.22), or compiler-specific extensions.
- If the code has a bug, an unsafe pattern, or undefined behavior, say so clearly and explain why. Do not "fix" it in the code.
- Some files are scratch or unfinished (for example `temp.go`, `test.go`, `need_to_test_setting_null_value.go`). Say that they appear to be experiments, and explain what they do anyway, briefly.
- Where a folder name or file name hints at a topic, you may use it to decide what to explain, but the explanation itself must come from the code.

## 5. Contents of each md file

Start with this frontmatter (the website reads it). Every field must be specific to THIS file. Never copy frontmatter from another note.

```
---
title: "<readable title, for example: Create a Container (Go)>"
date: 2026-10-25
tags:
  - <language: c, cpp, golang or asm>
  - <1 to 3 topic tags taken from the code, lowercase, hyphenated>
  - code-explained
summary: "<one sentence, under 160 characters, saying what this code does>"
---
```

- Quote `title` and `summary` with double quotes so colons do not break the YAML.
- Do not put a `# Title` heading in the body. The website already shows the frontmatter title.

Body, in this order. Skip a section only if it truly has nothing to say.

1. **One-paragraph overview.** What the program does and why someone would write it.
2. **Source:** a link to the code file on GitHub, written as
   `https://github.com/labib0x9/cs-notes/blob/main/<path-to-file>`
   with spaces in the path written as `%20`. Use the file's real path from the repo root. (Repo name is `cs-notes`; if the actual folder is named differently, use the folder name that exists in the workspace.)
3. **Concepts used.** A short list of the ideas this file demonstrates (for example "goroutines and channels", "pointer arithmetic", "BIOS interrupt 21h").
4. **How it works.** A step-by-step walkthrough in plain language, in the order the code runs. Refer to functions, structs and labels by NAME. Do not use line numbers, because they change.
5. **Key excerpt.** At most 25 lines of the most important code, in a fenced block with the right language tag (`c`, `cpp`, `go`, `asm`). Do not paste the whole file; link to it instead. If the file is under 25 lines, the whole file may be shown.
6. **Expected behavior (from reading the code, not run).** What a run would print or do, including the command to run it (for example `go run file.go`, `gcc file.c -o file && ./file`, `jwasm -mz file.asm` for the 8086 files). Include the command only; do not run it.
7. **Things to watch.** Real pitfalls, bugs, undefined behavior, platform limits, or version dependence found in THIS code. If there are none, omit the section.
8. **Terms.** One-line definitions of technical terms that a beginner might not know and that appear in this explanation. Define each term in plain words. Do not assume the reader knows jargon.
9. **Related.** Links to other notes, ONLY if the target exists. Use ordinary relative links to .md files, such as `[roadmap](../roadmap.md)` or `[deadlock notes](./deadlock.md)`. Before writing any link, check that the target file exists in the workspace (or is another md file you are creating in this run). Never link to a file that does not exist.

Last line of the body, in italics:
`*This explanation was generated by an AI from reading the source, not from running it. Verify before relying on it.*`

Length: proportional to the code. A 15-line program needs about 150 to 300 words. A large file may run to about 700 words. Do not pad.

Tone: clear, neutral, instructional. Do not write in the first person ("I learned...") because the notes are mine. Do not use hype words. Use "the code" or the function names as the subject.

## 6. Language-specific guidance

- **Go:** explain goroutines, channels, defer, interfaces, receivers, generics and similar features by what they do in the code. For `_test.go` files, explain what each test checks and show the command `go test -v ./...` run from that folder (do not run it). For `main` packages, explain the entry point first.
- **C and C++:** explain memory (stack, heap, ownership), pointers, headers and macros as used. For code that frees or reuses memory, say who owns what. Flag any buffer overflow, use-after-free, leaks, or uninitialised reads you can see. For `.c` files that depend on a local header (for example `math.h`, `custom header/`), read that header and mention it.
- **Assembly (8086, MASM-style `.asm`):** explain the segments (`.data`, `.code`, stack), the registers used, each interrupt call (for example `int 21h` with the function number in `ah`), and what each procedure does. Mention that these are 16-bit DOS programs meant for DOSBox.

## 7. Process

**Phase 0: inventory.** List every in-scope file, grouped by folder, and the md name each will get, including the collision renames. Print the totals per extension. Print the list of collisions. Do not ask for confirmation unless a collision is ambiguous.

**Phase 1: pilot.** Create md files for ONLY these 6 and then STOP and wait for my review:
1. `Programming/C/thread/create_thread.c`
2. `Programming/Golang System/02-Container/chroot.go`
3. `Programming/Assembly/8086-microprocessor/hello_world.asm`
4. `Programming/C++/pointer and dma/smart_pointer.cpp`
5. `Data Structures/Array/merge.c` and `Data Structures/Array/merge.go` (to prove the collision rule)
6. `Programming/Gloang Basics/05 - Array And Slice/slice_declare.go`

If any of these paths does not exist (folder names may change), pick the closest file of the same language and say which you used.

**Phase 2: the rest.** Only after I reply "continue", do all remaining files, folder by folder. Keep going without asking, but print a progress line after each top-level folder.

**Idempotent:** if an md that you would create already exists AND starts with the line `*This explanation was generated by an AI`... anywhere in it, skip it and count it as "already done". Never overwrite any other existing md.

## 8. Final report (print in chat; do NOT save it as a file in the repo)

- Number of md files created, per language.
- Number of files skipped, with the reason (already existed, empty, out of scope).
- The list of collision renames (`name.ext.md` forms).
- Files that need my attention: secrets found and redacted, code that looks buggy or unsafe, files you could not understand, files whose explanation you are unsure about.
- Things you did not do or could not verify. State this plainly. Do not claim the explanations are correct; say they come from reading only.

## 9. Do not

- Do not create index pages, readmes, folder summaries, or any file other than the md files described here.
- Do not add dependencies, scripts or config to the repo.
- Do not link to the internet except the GitHub source link and well-known official documentation, and only when you are sure the URL is correct. If unsure, leave the link out.