package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// implements io.ByteScanner
type byteReader struct {
	r         io.Reader
	buf       [1]byte
	unread    bool
	prevBuf   byte
	firstCall bool
}

func (b *byteReader) ReadByte() (byte, error) {
	var bb byte
	var err error
	b.firstCall = false

	if b.unread {
		bb = b.prevBuf
		b.unread = false
		return bb, err
	}

	n, err := io.ReadFull(b.r, b.buf[:]) // read 1 byte
	if n != 1 {
		return bb, errors.New("longer")
	}
	bb = b.buf[0]
	b.prevBuf = bb
	return bb, err
}

func (b *byteReader) UnreadByte() error {
	var err error
	if b.firstCall || b.unread {
		return errors.New("no byte to unread()")
	}
	b.unread = true
	return err
}

type reader struct {
	r io.ByteScanner
}

// we can to read from os.Stdin, which doesn't implement io.ByteScanner
func read(r io.Reader) (byte, error) {
	rdr := reader{}
	if rr, ok := r.(io.ByteScanner); ok {
		fmt.Println("Implements")
		rdr.r = rr
	} else {
		fmt.Println("Not implements")
		rdr.r = &byteReader{r: r, unread: false, firstCall: true}
	}

	return rdr.r.ReadByte()
}

func main() {

	/**		**/
	r := &byteReader{r: os.Stdin, unread: false, firstCall: true}

	// if err := r.UnreadByte(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Println(r.ReadByte()) // buf = {a} -> buf = {}
	r.UnreadByte()            // buf = {} -> buf = {a}
	fmt.Println(r.ReadByte()) // buf = {a} -> buf = {}

	/**		**/

	fmt.Println(read(os.Stdin))

}
