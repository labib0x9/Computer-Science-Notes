# Array And Slice
---
Precise reasoning: len, cap, backing array, reallocation.
---
# Slice Operations
- make()
- append()
- len(), cap()
- copy()
- Slicing [:]

---
# Notes
- Array has fixed size, Slice is a header struct with a backed array (we call it dynamic array).
- When we pass a array using Pass by value, entire array is being copied, for slice it is just 24byte header. 
- In function (pass by value), for array it doesn't affect the original array, but for slice it may modify if reallocation not occurs.
- Array not flexible, slice is flexible.
- `append()` returns a new slice.
- s[lo:hi:max] sets max capacity to slice. s[lo:hi], cap = cap(s) - lo, s[lo:hi:max], cap = max - lo
- This `a1` is array of zero length and capacity.
```go
a1 := [...]int{}         // same as [0]int
a2 := [...]int{1, 2, 3}  // same as [3]int{1, 2, 3}, Then why it is used ?
```
- 2D array
```go
var x [n][m]int
```
- Array vs Slice
```go
x := []int{1, 2, 3}     // Slice
x := [3]int{1, 2, 3}    // Array
```
---

# Quiz
This quizs are from ChatGPT. I will try to add more real interview quizs.

```text
Q1. Array vs Slice
    - List three fundamental differences that affects performance or behaviour.
Q2. Memory Model
    - What does a slice header contain?
    - Where is the underlying array allocated?
Q3. Value vs Reference Semantics
    - When you pass a slice to a function, what is copied?
    - What happens if the function appends beyond capacity?
Q4. Why can append be O(1) sometimes and O(n) other times?
Q5. How does Go decide when to allocate a new underlying array during append?
Q6. How can you force append not to modify the original slice’s underlying array?
Q7. What exactly does copy() copy?
```

```text
Q8. What is the output ?
    s := []int{1, 2, 3}
    t := s
    t[0] = 100
    fmt.Println(s)
```

```text
Q9. What is the output and explain what is happening.
    s := make([]int, 0, 2)
    s = append(s, 1)
    s = append(s, 2)
    t := append(s, 3)

    fmt.Println(s, t)
```

```text
Q10. What is the output and how to fix it ?
    s := []int{1, 2, 3, 4}
    a := s[1:3]
    a[0] = 100
    fmt.Println(s)
```

```text
Q11. What is wrong here? How do you fix it without global variables?
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

```text
Q12. What is the output ?
    s := make([]int, 2, 4)
    s[0], s[1] = 1, 2

    a := s[:2]
    b := append(a, 3)
    b[0] = 100

    fmt.Println(s, a, b)
```


```text
Q13. What is the output ?
    s := []int{1, 2, 3, 4}
    a := s[:2:2]   // note the third index
    b := append(a, 100)

    fmt.Println(s, a, b)
```

```text
Q14. What is the output ?
    s := []int{1, 2, 3}
    var p []*int

    for _, v := range s {
        p = append(p, &v)
    }

    fmt.Println(*p[0], *p[1], *p[2])
```

```text
Q15. What is the output ?
    var a []int
    b := []int{}

    fmt.Println(a == nil, b == nil)
    fmt.Println(len(a), len(b))
    fmt.Println(cap(a), cap(b))
```

```text
Q16. What is the output ?
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

```text
Q17. What is the output ? Why do b and c behave differently?
    s := []int{1, 2, 3, 4, 5}

    a := s[1:3:3]
    b := append(a, 100)
    c := append(s[1:3], 200)

    fmt.Println(s)
    fmt.Println(a, b, c)
```

```text
Q18. What is the output ? Why is this undefined-looking but valid Go?
    s := []int{1, 2}
    a := append(s, 3)
    b := append(s, 4)

    fmt.Println(a, b)
```

```text
Q19. What is the output ? Which parts are isolated? Which are not?
    s := []int{1, 2, 3}
    t := make([]int, 2)
    copy(t, s)

    t = append(t, 100)
    t[0] = 999

    fmt.Println(s, t)
```

```text
Q20. What is the output ? Under exactly which capacity condition does this break?
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

```text
Q21. How does make() works?
```
---