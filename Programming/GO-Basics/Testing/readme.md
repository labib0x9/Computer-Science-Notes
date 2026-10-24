go build -> does't include *_test.go files
go test -> includes *_test.go files

filename: Filename_test.go
function: TestFunction(t *testing.T) -> Test prefix is fixed

go test .
go test -v .
go test --cover .
go tool cover --html=coverage.go

// Testing functions in check_test.go
func TestPalindrome(t *testing.T) { ... }
func TestReverse(t *testing.T) { ... }
func TestSomethingElse(t *testing.T) { ... }

// tests only TestPalindrome
// from root folder
go test ./FileName -run=Palindrome -v

go test -run='(Palindrome|Reverse)' -v