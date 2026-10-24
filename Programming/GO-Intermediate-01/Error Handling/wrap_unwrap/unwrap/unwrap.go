package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Remove customized part and get the original error

func OpenFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("[error] [%v] OpenFile <---> %w", time.Now().Format("2006-01-02 15:04:05"), err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("[error]  Stat() function fails <---> %w", err)
	}
	_ = info
	return nil
}

// For errors.As()
type B struct {
	msg string
	val int64
}

func (b *B) Error() string {
	return fmt.Sprintf("%d cannot fit in %s data types", b.val, b.msg)
}

func checkBoundaryFor8ByteInt(n int64) error {
	if n > (1 << 8) {
		return &B{val: n, msg: "8byte"}
	}
	return nil
}

func main() {

	err := OpenFile("a.txt")
	fmt.Println(err)
	err = errors.Unwrap(err) // unwrap
	fmt.Println(err)
	fmt.Println()

	// If the full message is a customized one, then wrapping would give nil
	err = fmt.Errorf("aa")
	fmt.Println(err)
	err = errors.Unwrap(err) // unwrap
	fmt.Println(err)
	fmt.Println()

	// err is a customized one
	err = OpenFile("a.txt")
	fmt.Println(err)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file doesn't exist")
	}
	fmt.Println()
	// os.ErrClosed
	// os.ErrDeadlineExceeded
	// os.ErrExist
	// os.ErrInvalid
	// os.ErrNoDeadline
	// os.ErrNotExist
	// os.ErrPermission
	// os.ErrProcessDone
	// os.PathError

	// io.EOF
	// io.ErrClosedPipe
	// io.ErrNoProgress
	// io.ErrShortBuffer
	// io.ErrShortWrite
	// io.ErrUnexpectedEOF

	// net.ErrClosed
	// net.ErrWriteToConnected
	// net.UnknownNetworkError
	// net.InvalidAddrError
	// net.OpError

	// context.Canceled
	// context.DeadlineExceeded

	err = os.Mkdir("temp", 0755)
	if errors.Is(err, os.ErrExist) {
		fmt.Println("directory already exists")
	}
	fmt.Println()

	// target must be a pointer to an interface or struct type
	// err must not be a nil pointer
	err = OpenFile("a.txt")
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		fmt.Println("OP :", pathError.Op)
		fmt.Println("ERR :", pathError.Err)
		fmt.Println("PATH :", pathError.Path)
	}
	fmt.Println()

	// errors.Is() vs errors.As()
	// errors.Is() -> checks if a error is a sentinel like os.ErrInvalid (var)
	// errord.As() -> checks if a error is a certain data types (interface, struct)

	err = checkBoundaryFor8ByteInt(1 << 15)
	var boundaryErr *B
	if errors.As(err, &boundaryErr) {
		fmt.Println("value :", boundaryErr.val)
		fmt.Println("msg :", boundaryErr.msg)
	}
}
