package main

// Incomplete

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	// reader.ReadByte()
	// reader.ReadBytes()
	// reader.ReadLine()

	writer := bufio.NewWriter(os.Stdout)
	writer.WriteByte('a')
	// writer.WriteRune()
	// writer.WriteString()

	// What is the perpose of Flush() ??
	// Write any buffered data to the underlying output (stdout)
	// Without Flush(), buffered data might never be displayed
	writer.Flush()

	// var p []byte            // If we use thid p, then there will be no output as p is a zero capacity []byte
	p := make([]byte, 1024) // Now p is a 1kb buffer
	reader.Read(p)          // Read upto 1kb into p
	writer.Write(p)         // Write p

	reader.Peek(10)        // returns next 10bytes from reader without advancing the reader pointer
	reader.Discard(5)      // skips next 5bytes and advances the pointer
	reader.Reset(os.Stdin) // resets the reader

	// Use case of reset
	fileReader := bufio.NewReader(nil)
	for _, entry := range []string{"a.txt", "b.txt"} {
		file, _ := os.Open(entry)
		defer file.Close()
		fileReader.Reset(file)
	}

	writer.Available()      // how many buffers are available before flusing
	writer.Buffered()       // how many bytes waiting to flush
	writer.Reset(os.Stdout) // resets the writer

	writer = bufio.NewWriterSize(os.Stdout, 10) // small buffer for demo

	fmt.Println("Available:", writer.Available()) // 10 (buffer size)
	fmt.Println("Buffered:", writer.Buffered())   // 0

	writer.WriteString("abc")
	fmt.Println("Available after Write:", writer.Available()) // 7 (10 - 3)
	fmt.Println("Buffered after Write:", writer.Buffered())   // 3

	writer.Flush()
	fmt.Println("Available after Flush:", writer.Available()) // 10
	fmt.Println("Buffered after Flush:", writer.Buffered())   // 0
}
