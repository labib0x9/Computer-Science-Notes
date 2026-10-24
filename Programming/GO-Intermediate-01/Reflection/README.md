# Reflection
---
Rule
- Reflection goes from interface value to reflection object
```
interface{} -> reflect.Value, reflect.Type
```
- Reflection goes from reflection object to interface value
```
reflect.Value, reflect.Type -> interface{}
```
- To modify a reflection object, the value must be settable
```
v.CanSet() == True
```

---

Resource
- https://go.dev/blog/laws-of-reflection
- https://www.youtube.com/watch?v=Jvask1Hq_KE&t=226s&ab_channel=GopherConIsrael

- https://reliasoftware.com/blog/reflection-in-golang
- https://medium.com/capital-one-tech/learning-to-use-go-reflection-822a0aed74b7
- https://medium.com/capital-one-tech/learning-to-use-go-reflection-part-2-c91657395066
- https://research.swtch.com/interfaces
- https://www.slingacademy.com/article/using-reflection-for-generic-serialization-functions-in-go/
- https://www.youtube.com/watch?v=Kqt6EI4ypsk&ab_channel=EsotericTech
- https://jimmyfrasche.github.io/go-reflection-codex/
- https://medium.com/kokster/go-reflection-creating-objects-from-types-part-i-primitive-types-6119e3737f5d
- https://medium.com/kokster/go-reflection-creating-objects-from-types-part-ii-composite-types-69a0e8134f20
- https://www.geeksforgeeks.org/go-language/reflect-makeslice-function-in-golang-with-examples
- https://www.geeksforgeeks.org/go-language/reflect-new-function-in-golang-with-examples/
- https://www.tutorialspoint.com/go/go_reflection.htm
- https://pkg.go.dev/reflect
