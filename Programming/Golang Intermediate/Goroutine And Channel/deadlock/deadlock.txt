---
title: Deadlock, Lock Ordering, and Copying Locks
date: 2026-10-24
tags:
  - golang
  - concurrency
  - deadlock
summary: Three versions of one Go program showing a deadlock, the lock-ordering fix, and why copying a mutex breaks locking
---

Three versions of the same program (`printSum` with two `value` structs, each holding a mutex). Same code shape, three different behaviors.

```go
type value struct {
    mu    sync.Mutex
    value int
}
```

`printSum(id, v1, v2)` does: lock `v1` -> sleep 20s -> lock `v2` -> print `v1.value + v2.value`.

| Version | What changed | Result |
|---|---|---|
| 1 | `a` calls `(a, b)`, `b` calls `(b, a)` | **Deadlock** |
| 2 | `b` swaps so it also locks `a` first | **No deadlock** |
| 3 | Passes `value` **by value** (copies) | **No deadlock, but locks protect nothing** |

---

## Version 1: The deadlock

```go
go printSum("a", &a, &b)   // locks a, then wants b
go printSum("b", &b, &a)   // locks b, then wants a
```

Timeline:

```
t=0s    G_a locks a.mu           G_b locks b.mu
t=0-20s both sleep (each HOLDING one lock)
t=20s   G_a wants b.mu (held by G_b) -> blocks
        G_b wants a.mu (held by G_a) -> blocks
        main is blocked in wg.Wait()
        => every goroutine is asleep => runtime crashes
```

**Why it is a deadlock:** each goroutine holds a lock the other needs, and neither will release until it gets the second lock. Nobody can move.

**The four conditions (all four must hold for a deadlock):**

1. **Mutual exclusion:** a lock has one owner at a time.
2. **Hold and wait:** a goroutine holds one lock while waiting for another.
3. **No preemption:** you cannot forcibly take a lock from its owner.
4. **Circular wait:** a cycle of goroutines, each waiting for the next one's lock.

Prevention means breaking at least one. Version 2 breaks #4.

**Why `time.Sleep(20 * time.Second)` is there:** it widens the window between "lock first" and "lock second", so the deadlock happens every time instead of occasionally. Without the sleep, the bug is still present, but a goroutine might finish both locks before the other starts, so it would only deadlock on unlucky timing. That makes it a **race condition** too (correctness depends on timing), which is exactly why deadlocks are hard to catch in testing.

**Reading the crash:**

- It appears after about 20 seconds, not instantly. During the sleep, the goroutines are not "stuck", because a timer will wake them. The runtime only declares deadlock once nothing can ever wake anyone.
- `goroutine 4 [sync.Mutex.Lock]` and `goroutine 5 [sync.Mutex.Lock]` are the two workers, both at the line `v2.mu.Lock()`.
- The argument order shows the swap. G4 got `(..120, ..130)`, G5 got `(..130, ..120)`. The `lockSlow(0x...)` address is the mutex each one is waiting for.
- `goroutine 1 [sync.WaitGroup.Wait]` is main, waiting for workers that will never finish.
- Goroutine IDs (4, 5, 7, 8, 19, 20...) are only labels. They differ between runs and you should never build logic on them.

---

## Version 2: Why swapping fixes it

```go
if id == "b" {
    v1, v2 = v2, v1     // now b also locks a first, then b
}
```

Now both goroutines lock in the **same order**: `a`, then `b`.

```
G_a: lock a -> sleep -> lock b -> print -> unlock b -> unlock a
G_b: lock a (waits until G_a unlocks a) -> sleep -> lock b -> print -> ...
```

Whoever gets `a.mu` first wins. The other one waits on `a.mu` while holding **nothing**, so it cannot be part of a cycle.

**Why consistent ordering makes deadlock impossible (the proof idea):**
Give every lock a rank (`a` = 1, `b` = 2). If every goroutine only ever waits for a lock with a *higher* rank than the ones it holds, then a cycle would need ranks that go up forever around a loop, which cannot happen. No cycle means no circular wait means no deadlock.

**Things to notice:**

- The swap changes only the function's **local** copies of the pointers. The caller's `a` and `b` are untouched.
- **Total time is about 40 seconds, not 20.** The second goroutine cannot even start its sleep until the first one releases `a`. The ordering fix makes the two serialize (run one after the other).
- Output `a + b = 0` for both is expected, since nobody changed the values.
- The `race` detector reports nothing here, and it would also report nothing for Version 1. **The race detector does not find deadlocks.** It finds data races only.

**The weakness of this fix:** `if id == "b"` is a hard-coded hack. It works for exactly two named goroutines. A general fix is to give each lockable thing a stable id and always lock the lower id first:

```go
type value struct {
    id    int
    mu    sync.Mutex
    value int
}
// inside printSum: if v1.id > v2.id { v1, v2 = v2, v1 }
```

This is the same idea used for "transfer money between accounts A and B": lock the lower account id first, regardless of which direction the money moves.

---

## Version 3: Why copying "fixes" it (it doesn't, really)

```go
printSum := func(id string, v1, v2 value) { ... }   // value, not *value
go printSum("a", a, b)
go printSum("b", b, a)
```

Passing a struct **by value** copies the whole struct, including the `sync.Mutex` inside it. So:

- Goroutine "a" gets its own private copies of `a` and `b`, with their own mutexes.
- Goroutine "b" gets a different pair of private copies.
- The four mutexes are all different objects. Nobody is waiting on anybody, so there is no deadlock.

**But the program is now broken in a quieter way:** the locks no longer protect anything shared. If another goroutine changed the real `a.value`, these goroutines would not see it, and the copies could be read while someone else writes them (a data race on the copy step). The mutex has become decoration.

**Rules for `sync` types (`Mutex`, `WaitGroup`, `RWMutex`, `Once`, ...):**

1. **Never copy them after first use.** The Go docs say this explicitly. A mutex is state (locked/unlocked plus a wait queue), and a copy gets a snapshot of that state.
2. **If you copy a *locked* mutex, the copy is locked too**, and nobody holds the original's unlock responsibility for it. That can create a deadlock out of nowhere.
3. Pass structs containing locks **by pointer** (`*value`), as in Versions 1 and 2.

**Tool that catches this:** `go vet` flags it with a message like `passes lock by value` (the `copylocks` check). Run `go vet ./...` regularly.

---

## Side-by-side summary

| | Same locks shared? | Lock order | Deadlock? | Correct? |
|---|---|---|---|---|
| V1 | Yes (pointers) | Opposite (`a,b` vs `b,a`) | **Yes** | No |
| V2 | Yes (pointers) | Same (`a,b` both) | No | **Yes** |
| V3 | **No** (copies) | Irrelevant | No | No (locks are meaningless) |

---

## Checklist: avoiding deadlocks in real code

- **Always acquire multiple locks in one global order** (by id, by name, by rank).
- **Hold fewer locks, for less time.** Do not sleep, do I/O, or call unknown code while holding a lock.
- **Do not copy structs that contain locks.** Use pointers. Run `go vet`.
- **Use `defer mu.Unlock()`** so early returns cannot leave a lock held.
- **Use timeouts** (`context.WithTimeout` with `select`) for channel operations, so a stuck wait becomes an error instead of a silent hang.
- **Prefer one owner goroutine plus channels** when the design allows, because then there are no locks to order.

## Debugging a deadlock

- If *all* goroutines are blocked, Go crashes with `fatal error: all goroutines are asleep - deadlock!` and prints every stack.
- If *some* goroutine is still running (a server, a timer loop), Go stays silent. Get a stack dump with `kill -QUIT <pid>` or Ctrl+\, or open `/debug/pprof/goroutine?debug=2` if `net/http/pprof` is enabled.
- In the dump, look for goroutines in `sync.Mutex.Lock`, `chan send`, `chan receive`, or `sync.WaitGroup.Wait` for a long time, and match the lock addresses to find the cycle.
- `GOTRACEBACK=system` also shows the runtime's own helper goroutines (GC, finalizer, etc.). They are idle, not stuck.

## Terms used

- **Data race:** two goroutines access the same variable, at least one writes, no synchronization.
- **Race condition:** the correct result depends on timing or order. Can exist without a data race.
- **Deadlock:** goroutines blocked forever, each waiting on something another blocked goroutine holds.
- **Livelock:** goroutines keep running and reacting but make no progress.
- **Starvation:** one goroutine never gets the resource because others keep winning.
- **Serialize:** make operations happen one after another instead of at the same time.