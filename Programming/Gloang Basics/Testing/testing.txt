---
title: Testing in Go
date: 2026-10-25
tags:
  - golang
  - testing
summary: How go test finds tests, how to run them, and how to measure coverage
---

## Test files and functions

- `go build` does not include `*_test.go` files.
- `go test` includes `*_test.go` files.
- File name: it must end with `_test.go`, for example `check_test.go`.
- Function name: `TestXxx(t *testing.T)`. The `Test` prefix is fixed, and the next letter must be uppercase (`TestPalindrome` runs, `Testpalindrome` does not).

```go
// Tests in check_test.go
func TestPalindrome(t *testing.T) { ... }
func TestReverse(t *testing.T) { ... }
func TestSomethingElse(t *testing.T) { ... }
```

## Running tests

```bash
go test .          # run tests in the current package
go test -v .       # verbose: print every test and its result
go test -cover .   # also print the coverage percentage
```

## Coverage report

Coverage is the share of your code lines that the tests actually run. To see it as an HTML page, write a profile file first, then open it:

```bash
go test -coverprofile=coverage.out .
go tool cover -html=coverage.out
```

## Running specific tests

`-run` takes a regular expression (a text pattern). It matches any test whose name *contains* the pattern.

```bash
# from the root folder: run tests in the ./FileName package whose name contains "Palindrome"
go test ./FileName -run=Palindrome -v

# run tests whose name contains "Palindrome" or "Reverse"
go test -run='(Palindrome|Reverse)' -v

# run exactly TestPalindrome and nothing else
go test -run='^TestPalindrome$' -v
```