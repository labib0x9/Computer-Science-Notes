# select 
---

```go
select {
    case <-done :
    case msg := <-ch :
    case <- time.After(expireDuration) :
    default :
        // Default
}
```

---