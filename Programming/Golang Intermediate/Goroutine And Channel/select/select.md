---
title: Select
date: 2026-10-24
tags:
  - golang
  - concurrency
  - deadlock
summary: 
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