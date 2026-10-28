---
title: Reflection in Go
date: 2026-10-24
tags:
  - golang
  - reflection
summary: The three laws of reflection in Go, with a reading list
---

## The three laws of reflection

1. Reflection goes from an interface value to a reflection object.

```go
t := reflect.TypeOf(x)    // reflect.Type
v := reflect.ValueOf(x)   // reflect.Value
```

2. Reflection goes from a reflection object back to an interface value.

```go
i := v.Interface()        // interface{}
```

3. To modify a reflection object, the value must be settable.

```go
v.CanSet() == true
```

## Resources

- https://go.dev/blog/laws-of-reflection
- https://pkg.go.dev/reflect
- https://research.swtch.com/interfaces
- https://jimmyfrasche.github.io/go-reflection-codex/
- https://www.youtube.com/watch?v=Jvask1Hq_KE&t=226s
- https://www.youtube.com/watch?v=Kqt6EI4ypsk
- https://reliasoftware.com/blog/reflection-in-golang
- https://medium.com/capital-one-tech/learning-to-use-go-reflection-822a0aed74b7
- https://medium.com/capital-one-tech/learning-to-use-go-reflection-part-2-c91657395066
- https://medium.com/kokster/go-reflection-creating-objects-from-types-part-i-primitive-types-6119e3737f5d
- https://medium.com/kokster/go-reflection-creating-objects-from-types-part-ii-composite-types-69a0e8134f20
- https://www.slingacademy.com/article/using-reflection-for-generic-serialization-functions-in-go/
- https://www.geeksforgeeks.org/go-language/reflect-makeslice-function-in-golang-with-examples
- https://www.geeksforgeeks.org/go-language/reflect-new-function-in-golang-with-examples/
- https://www.tutorialspoint.com/go/go_reflection.htm