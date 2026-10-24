package main

import (
	"fmt"
	"io"
	"syscall"
)

type syscallReader struct {
	fd int	// file descriptor
}

// Necessary method to implement io.Reader
func (s *syscallReader) Read(p []byte) (n int, err error) {
	n, err = syscall.Read(s.fd, p)
	return
}

// Extra method
func (s *syscallReader) Temp() {
	fmt.Println("TEMPORARY CALL")
}

func ReadTenBytes(r io.Reader) {
	b := make([]byte, 10)
	if _, err := r.Read(b); err != nil {
		return
	}
	fmt.Println(b)
}

func main() {

	r := &syscallReader{0} // syscallReader implements io.Reader, 0 is file descriptor for stdin
	ReadTenBytes(r)
}
