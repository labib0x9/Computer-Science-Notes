---
title: Communicate with the Docker Daemon Using the Docker API
date: 2026-10-24
tags:
  - golang
  - docker
summary: Connecting to the Docker daemon from Go, and a go.mod version note
---

## Important

Sometimes you have to pin a specific version of the Docker package, because the package has changed between versions.

In `go.mod`, under `require`:

```go
require github.com/docker/docker v28.5.2+incompatible
```