package main

import (
	"bytes"
	"fmt"
	"io"
)

// Incomplete

func main() {
	// Most of the function are from strings package, which are not shown in this file

	var buffer bytes.Buffer

	buffer.WriteByte(byte('a'))
	buffer.WriteRune(rune('a'))
	buffer.WriteString("str")

	fmt.Println(buffer.String())

	b, err := buffer.ReadByte()
	r, _, err := buffer.ReadRune()

	line, err := buffer.ReadBytes(byte('\n'))
	str, err := buffer.ReadString('\n')

	if err != nil && err != io.EOF {

	}

	fmt.Println(b, line, r, str)

	nextBBytes := buffer.Next(5)
	fmt.Println(string(nextBBytes))

	buffer.Reset()

	buffer.Bytes()
	buffer.Truncate(0)
	buffer.Available()
}