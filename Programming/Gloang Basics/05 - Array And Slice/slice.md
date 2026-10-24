---
title: Array and Slice in Go
date: 2026-10-25
tags:
  - golang
  - slices
  - arrays
summary: How arrays and slices differ in Go (len, cap, backing array, reallocation), with practice questions
---

Precise reasoning about `len`, `cap`, the backing array, and reallocation.

## Slice operations

- `make()`
- `append()`
- `len()`, `cap()`
- `copy()`
- Slicing with `[:]`

## Notes

- An array has a fixed size. A slice is a small header struct that points to a backing array (we often call it a dynamic array).
- When you pass an array by value, the entire array is copied. For a slice, only the header is copied, which is 24 bytes on a 64-bit machine (a pointer, a length, and a capacity).
- Inside a function (pass by value), changing an array does not affect the original. Changing a slice element does affect the original, as long as no reallocation happened first.
- Arrays are not flexible. Slices are.
- `append()` returns a slice, which may or may not share the old backing array. Always use the returned value.
- `s[lo:hi]` has `cap = cap(s) - lo`. `s[lo:hi:max]` has `cap = max - lo`, so the third index limits the capacity.
- `a1` below is an array of zero length and zero capacity.

```go
a1 := [...]int{}         // same as [0]int
a2 := [...]int{1, 2, 3}  // same as [3]int{1, 2, 3}
```

Why use `[...]`? The compiler counts the elements for you, so you don't have to update the size by hand when you add or remove one.

- 2D array (`n` and `m` must be constants):

```go
var x [n][m]int
```

- Array vs slice (use different variable names, because `:=` twice with the same name in one scope does not compile):

```go
s := []int{1, 2, 3}     // slice
a := [3]int{1, 2, 3}    // array
```

## Quiz

These questions are from ChatGPT. I will try to add more real interview questions.

**Q1. Array vs slice.** List three fundamental differences that affect performance or behavior.

**Q2. Memory model.** What does a slice header contain? Where is the underlying array allocated?

**Q3. Value vs reference semantics.** When you pass a slice to a function, what is copied? What happens if the function appends beyond capacity?

**Q4.** Why can `append` be O(1) sometimes and O(n) other times?

**Q5.** How does Go decide when to allocate a new underlying array during `append`?

**Q6.** How can you force `append` not to modify the original slice's underlying array?

**Q7.** What exactly does `copy()` copy?

**Q8.** What is the output?

```go
s := []int{1, 2, 3}
t := s
t[0] = 100
fmt.Println(s)
```

**Q9.** What is the output? Explain what is happening.

```go
s := make([]int, 0, 2)
s = append(s, 1)
s = append(s, 2)
t := append(s, 3)

fmt.Println(s, t)
```

**Q10.** What is the output, and how do you fix it?

```go
s := []int{1, 2, 3, 4}
a := s[1:3]
a[0] = 100
fmt.Println(s)
```

**Q11.** What is wrong here? How do you fix it without global variables?

```go
func modify(s []int) {
    s = append(s, 10)
    s[0] = 99
}

func main() {
    s := []int{1, 2, 3}
    modify(s)
    fmt.Println(s)
}
```

**Q12.** What is the output?

```go
s := make([]int, 2, 4)
s[0], s[1] = 1, 2

a := s[:2]
b := append(a, 3)
b[0] = 100

fmt.Println(s, a, b)
```

**Q13.** What is the output?

```go
s := []int{1, 2, 3, 4}
a := s[:2:2]   // note the third index
b := append(a, 100)

fmt.Println(s, a, b)
```

**Q14.** What is the output?

```go
s := []int{1, 2, 3}
var p []*int

for _, v := range s {
    p = append(p, &v)
}

fmt.Println(*p[0], *p[1], *p[2])
```

**Q15.** What is the output?

```go
var a []int
b := []int{}

fmt.Println(a == nil, b == nil)
fmt.Println(len(a), len(b))
fmt.Println(cap(a), cap(b))
```

**Q16.** What is the output?

```go
func g(s []int) []int {
    s[0] = 999
    return append(s, 5)
}

func main() {
    s := []int{1, 2, 3}
    t := g(s)
    fmt.Println(s, t)
}
```

**Q17.** What is the output? Why do `b` and `c` behave differently?

```go
s := []int{1, 2, 3, 4, 5}

a := s[1:3:3]
b := append(a, 100)
c := append(s[1:3], 200)

fmt.Println(s)
fmt.Println(a, b, c)
```

**Q18.** What is the output? Why does this look undefined but is valid Go?

```go
s := []int{1, 2}
a := append(s, 3)
b := append(s, 4)

fmt.Println(a, b)
```

**Q19.** What is the output? Which parts are isolated, and which are not?

```go
s := []int{1, 2, 3}
t := make([]int, 2)
copy(t, s)

t = append(t, 100)
t[0] = 999

fmt.Println(s, t)
```

**Q20.** What is the output? Under exactly which capacity condition does this break?

```go
func mutate(s []int) {
    s = append(s, 5)
    s[0] = 777
}

func main() {
    s := []int{1, 2, 3}
    mutate(s[:2])
    fmt.Println(s)
}
```

**Q21.** How does `make()` work?