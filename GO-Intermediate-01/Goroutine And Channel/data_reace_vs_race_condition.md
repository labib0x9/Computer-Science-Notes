# Go: Data Race vs Race Condition

Two **independent** properties. Neither implies the other.

| | Definition | Detected by `-race`? |
|---|---|---|
| **Data race** | ≥2 goroutines access the same memory concurrently, ≥1 is a write, no synchronization (no happens-before edge). | Yes (only on paths actually executed) |
| **Race condition** | Program correctness depends on timing/ordering of goroutines. A logic bug. | **No** |

- Data race = a *memory-level* property. Fixed by mutex / atomic / channel.
- Race condition = a *semantic* property. Fixed by correct coordination design (WaitGroup, channels, ordering rules, atomic check-and-act). Whether it's a bug depends on your spec: if every ordering is acceptable, it's not a bug.

## 2×2 matrix

| | **No race condition** | **Race condition** |
|---|---|---|
| **No data race** | Ex 4: mutex + WaitGroup | Ex 2: mutex only · Ex 5: last-writer-wins |
| **Data race** | Ex 6: two goroutines write same constant | Ex 1: map · Ex 3: WaitGroup only |

Test: `go build -race` / `go run -race` / `go test -race`.

---

## Ex 1 — Data race + race condition (map)

```go
func main() {
	mp := make(map[string]int)

	go func() {
		mp["foo"] = 1 // write, goroutine
	}()

	mp["bar"] = 2 // write, main

	for k, v := range mp { // read, main
		fmt.Println("K=", k, "V=", v)
	}
}
```

**Data race:** goroutine and main write `mp` with no sync.

**Race condition:** `main` never waits for the goroutine, so `foo` may or may not exist when the loop runs. Output varies:

```
K= bar V= 2
K= foo V= 1      <- sometimes
```

Race detector:

```
WARNING: DATA RACE
Write at 0x... by goroutine 7:
  runtime.mapassign_faststr()
  main.main.func1()
Previous write at 0x... by main goroutine:
  runtime.mapassign_faststr()
  main.main()
```

Notes:
- Without `-race`, the runtime itself may abort with `fatal error: concurrent map writes` or `concurrent map iteration and map write`. That is a best-effort check, not a guarantee, and it is **unrecoverable** (not a `panic`, `recover()` won't catch it).
- Only mutex added → data race gone, race condition stays (`foo` still may be missing).
- Mutex **and** `wg.Wait()` before the `range` → both gone. (Iteration *order* is still random, but that's map semantics, not a race: the key set is deterministic.)

Fixed:

```go
func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	mp := make(map[string]int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		mp["foo"] = 1
		mu.Unlock()
	}()

	mu.Lock()
	mp["bar"] = 2
	mu.Unlock()

	wg.Wait() // happens-before: all writes done
	for k, v := range mp { // safe, no concurrent writers left
		fmt.Println("K=", k, "V=", v)
	}
}
```

---

## Ex 2 — NO data race, but race condition (mutex only)

```go
func main() {
	var mu sync.Mutex
	var value int

	for range 10000 {
		go func() {
			mu.Lock()
			value++
			mu.Unlock()
		}()
	}

	mu.Lock()
	fmt.Println("It should be 10000, GOT=", value)
	mu.Unlock()
}
```

```
It should be 10000, GOT= 5762    (varies every run)
```

- Every access to `value` is locked → **no data race**, `-race` is silent.
- `main` doesn't wait for goroutines to *finish*, only for the lock → prints whatever count has happened so far. Mutex gives **mutual exclusion**, not **completion ordering**. → **race condition**.
- Also leaks unfinished goroutines when `main` returns.

---

## Ex 3 — Data race + race condition (WaitGroup only)

```go
func main() {
	var wg sync.WaitGroup
	var value int

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value++ // read-modify-write, unsynchronized
		}()
	}

	wg.Wait()
	fmt.Println("It should be 10000, GOT=", value)
}
```

- `wg.Wait()` guarantees all goroutines finished, so the *final read* is safe.
- But `value++` is **load → add → store**. Two goroutines interleave → lost update → **data race** and a **race condition** (wrong total).
- Race report: `Read ... by goroutine N` / `Previous write ... by goroutine M`, both at the `value++` line.
- Lost updates are probabilistic. In a re-run under `-race` the sum happened to print `10000` while the detector still reported the race. **Correct-looking output does not mean no race.** The detector is the authority.

---

## Ex 4 — Neither (mutex + WaitGroup)

```go
func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var value int

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			value++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("It should be 10000, GOT=", value) // always 10000
}
```

- Mutex → mutual exclusion (no data race).
- WaitGroup → completion ordering (no race condition).
- Idiomatic alternative for a plain counter: `atomic.Int64` + `Add(1)`.

---

## Ex 5 — No data race, race condition (nondeterministic winner)

```go
func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var finisherId int

	for i := range 10000 {
		wg.Add(1)
		go func(id int, finisherId *int) {
			defer wg.Done()
			mu.Lock()
			*finisherId = id
			mu.Unlock()
		}(i, &finisherId)
	}

	wg.Wait()

	switch finisherId % 3 {
	case 0: fmt.Println("Tic")
	case 1: fmt.Println("Tac")
	case 2: fmt.Println("Toe")
	}
}
```

- Writes are locked, final read is after `wg.Wait()` → **no data race**.
- Final value = id of whichever goroutine ran **last**, decided by the scheduler → output changes between runs (`Tic`/`Tac`/`Toe`). **Race condition.**
- Proves: perfectly synchronized memory access can still produce nondeterministic results.
- Whether it's a *bug* depends on intent. If "any winner is fine", it isn't.

---

## Ex 6 — Data race WITHOUT race condition (the missing quadrant)

```go
func main() {
	var wg sync.WaitGroup
	x := 0

	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			x = 1 // both write the same value
		}()
	}

	wg.Wait()
	fmt.Println("x =", x) // always 1
}
```

- Result is deterministic (always `1`) → no race condition.
- Two unsynchronized writes to `x` → `-race` still reports **DATA RACE**.
- In Go, "benign" data races don't exist. The memory model gives no guarantees for racy programs (compiler/CPU reordering, torn writes for multiword values like strings, slices, interfaces). Fix it anyway.

---

## Cheat sheet

| Symptom | Likely cause | Tool |
|---|---|---|
| `-race` reports | Data race | mutex / `sync/atomic` / channel |
| Wrong/varying result, `-race` clean | Race condition | `WaitGroup`, channels, redesign ordering, atomic check-and-act |
| `fatal error: concurrent map writes` | Data race on map | mutex or `sync.Map` |
| Total too low after goroutines | Lost update (`x++` not atomic) | mutex or `atomic.Add` |
| Reads partial results | Missing completion sync | `wg.Wait()` / channel close |

Rules:
1. Mutex = **"one at a time"**. WaitGroup = **"wait until done"**. You usually need to know which one your bug needs.
2. `-race` finds data races that *executed* in that run. Clean run ≠ proven safe. Run in CI with `go test -race ./...` and with enough load.
3. Check-then-act (`if !exists { create }`) split across two lock sections is a race condition even if each section is locked. Hold the lock across the whole check+act.
4. Prefer sharing by communicating (channels) or ownership (one goroutine owns the data).

## Build commands

```bash
go build -race -o out main.go        # detector build (slower, ~5-10x CPU, ~5-10x mem)
go run -race main.go
go test -race ./...
go build -gcflags="all=-N -l" -o out main.go   # no optimization/inlining; makes debugging easier, doesn't cause or fix races
```

`for range 10000` and per-iteration loop variables need Go 1.22+.