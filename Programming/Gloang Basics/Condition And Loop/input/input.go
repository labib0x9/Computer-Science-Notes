package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func scanfInput() {
	var a, b int
	fmt.Scanf("%d %d", &a, &b)
}

// Efficient
func bufioInput() {
	// Reader
	bufReader := bufio.NewReader(os.Stdin)
	a, err := bufReader.ReadByte()
	b, err := bufReader.ReadBytes(byte('\n'))
	c, _, err := bufReader.ReadLine()
	d, _, err := bufReader.ReadRune()
	e, err := bufReader.ReadSlice('\n')
	f, err := bufReader.ReadString(',')

	if err != nil {
		log.Println("error: ", err)
	}
	fmt.Println(a, b, c, d, e, f)

	// Scanner
	bufScanner := bufio.NewScanner(os.Stdin)
	found := bufScanner.Scan()
	text := bufScanner.Text()
	bufScanner.Split(bufio.ScanWords)
	if err := bufScanner.Err(); err != nil {
		log.Println("error: ", err)
	}
	fmt.Println(found, text)
}

func chunkInput() {
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 4096) // 4kb buffer
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Process buf[:n]
		}
		if err == io.EOF {
			break
		}
		// Handle err
		if err != nil {
			
		}
	}
}

func main() {
	bufioInput()
}
