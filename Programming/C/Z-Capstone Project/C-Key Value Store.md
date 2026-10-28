---
title: "Capstone Project: PIKA Key-Value Store Using Hashmap"
date: 2026-10-28
tags:
  - c
  - capstone
---

# The key-value Store

The store is independent of the parser and the network. It takes plain C values and returns status codes, so it can be tested on its own.

## Data model

```text
store
 └── bucket table          name  -> bucket
      └── key table        key   -> value list
           └── values      ["sessionid_1019", "sessionid_1020"]
```

Both tables are the same generic **hash map**, so you write the hash table once. A hash map turns a key into an array index with a hash function; two keys landing on the same index is a **collision**.

## Hash map interface

A suggested `hashmap.h`, declarations only:

```c
#ifndef HASHMAP_H
#define HASHMAP_H

#include <stddef.h>

typedef struct hashmap hashmap;

/* free_value is called on a value when it is replaced, removed or the map is destroyed. */
hashmap *hm_create(void (*free_value)(void *));
void     hm_destroy(hashmap *m);

int      hm_put(hashmap *m, const char *key, void *value);  /* 0 ok, -1 no memory */
void    *hm_get(const hashmap *m, const char *key);         /* NULL if absent */
int      hm_remove(hashmap *m, const char *key);            /* 1 removed, 0 absent */
size_t   hm_count(const hashmap *m);

/* stop early if fn returns non-zero */
void     hm_foreach(const hashmap *m,
                    int (*fn)(const char *key, void *value, void *ctx),
                    void *ctx);

#endif
```

## Design choices

| Choice | Suggested | Reason |
| --- | --- | --- |
| Collision handling | **chaining**: each slot is a singly linked list | simplest to get right, removal is easy |
| Hash function | FNV-1a, 64-bit (offset `14695981039346656037`, prime `1099511628211`, per byte: XOR then multiply) | tiny, good enough for strings |
| Slot count | power of two, start at 16 | index is `hash & (n - 1)`, no modulo |
| Growth | double when `count > 0.75 * slots` (**load factor** = entries per slot) | keeps chains short |
| Rehash | relink existing nodes into the new array; don't copy keys again | one allocation, no key churn |
| Key storage | the map owns its own copy of each key | the caller's buffer may disappear |

## Ownership rules

| Situation | Who owns what |
| --- | --- |
| `hm_put` succeeds | the map owns `value` from now on |
| `hm_put` fails (returns -1) | the caller still owns `value` and must free it |
| `hm_put` on an existing key | the old value goes to `free_value`, then the new one is stored |
| `hm_get` | returns a borrowed pointer, valid only until the next `hm_put`/`hm_remove` on that map |
| `hm_destroy` | frees every key, every value (via callback) and the table |

The value stored in the key table is a `struct { char **items; size_t count, cap; }`. Its free function frees each string, the array, then the struct.

## Store operations

| Operation | Result codes |
| --- | --- |
| `store_create(name)` | OK, EXISTS, NOMEM |
| `store_delete(name)` | OK, NO_BUCKET |
| `store_set(bucket, key, value, append)` | OK, NO_BUCKET, LIMIT, NOMEM |
| `store_get(bucket, key)` | OK + borrowed list, NO_BUCKET, NO_KEY |
| `store_remove(bucket, key)` | OK, NO_BUCKET, NO_KEY |
| `store_list(bucket, after, limit, ...)` | OK, NO_BUCKET |

## Listing and `after`

A hash map has **no useful order**: iteration follows the slots, and the order changes when the table grows. `after` only makes sense with a stable order, so for `store_list`:

1. Collect pointers to all keys with `hm_foreach` into an array.
2. Sort the array with `qsort` and `strcmp`.
3. Skip every key `<= after`, then take up to `limit`.

This costs O(n log n) per list call. That's fine for a toy, and you can note the trade-off. Keep the store behind this small interface so a lock can wrap it in the threading stage.

## Exercise 02: Hash map and store

Requirements:
- `hashmap.c` passes the tests below with no leaks and no ASan or UBSan reports.
- Check every allocation. If growth fails midway, the map must still be valid and keep every entry (the old table stays in use).
- `store.c` enforces the limits table from section 1 and returns `LIMIT` when exceeded.
- No `strcpy` or `strncpy`: use `malloc` + `memcpy` with the length you already know.

Hints:
- Write `hm_put`, `hm_get` and `hm_remove` first with a fixed slot count; add growth afterwards and check that old tests still pass.
- To remove from a singly linked chain, keep a pointer to the **link** (`node **`), not just the node.
- Test growth by inserting 10000 keys, then checking all of them, then removing all.

## Testing

| # | Test | Expect |
| --- | --- | --- |
| 1 | put `a`, get `a` | same pointer back |
| 2 | get a key never put | `NULL` |
| 3 | put `a` twice with different values | second value returned, `free_value` called once for the first |
| 4 | remove `a`, then get `a` | `NULL`, `hm_count` decreases |
| 5 | remove a key that doesn't exist | returns 0 |
| 6 | insert 10000 distinct keys | all retrievable, `hm_count` is 10000 |
| 7 | remove all 10000 | `hm_count` is 0, no leaks |
| 8 | two keys forced into one chain (small table) | both retrievable, removing one keeps the other |
| 9 | `store_set` without append, then with append | list has 1 value, then 2 in order |
| 10 | `store_list` after inserting `b`, `a`, `c` | order `a`, `b`, `c`; `after: a` gives `b`, `c` |
| 11 | `store_delete` on a bucket with keys | no leaks |