package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var tokens = make(chan struct{}, 1) // For maximum 1 goroutine to run concurrently

type Result struct {
	filename                                   string
	byteCount, lineCount, wordCount, charCount int
	err                                        error
}

// Process rune by rune from input and count results
func inputProcess(r io.Reader) (result Result) {
	reader := bufio.NewReader(r)
	inWord := false

	for {
		in, size, err := reader.ReadRune() // Reads a single rune
		if err == io.EOF {                 // Check if for End Of File
			break
		}
		if err != nil {
			log.Println(err)
			result.err = err
			break
		}

		if in == '\n' { // Newline count
			result.lineCount++
		}

		if unicode.IsSpace(in) {
			inWord = false
		} else if !inWord {
			inWord = true
			result.wordCount++ // Word count
		}

		result.charCount++       // Char count
		result.byteCount += size // Byte count
	}
	return
}

// Print output
func PrintOutput(result Result) {
	output := make([]string, 0)
	output = append(output, strconv.Itoa(result.lineCount))
	output = append(output, strconv.Itoa(result.wordCount))
	output = append(output, strconv.Itoa(result.byteCount))
	output = append(output, strconv.Itoa(result.charCount))
	if result.err == nil {
		fmt.Println("  " + strings.Join(output, "   ") + "   " + result.filename)
	} else {
		log.Println(result.err)
	}
}

func main() {

	ch := make(chan Result) // Channel
	var wg sync.WaitGroup   // Waitgroup for goroutine

	files, err := os.ReadDir("test")
	if err != nil {

	}

	for _, filename := range files {

		wg.Add(1) // Wait for 1 goroutine

		go func(filepath string, ch chan Result) {
			defer wg.Done()
			tokens <- struct{}{}
			f, err := os.Open(filepath)
			if err != nil {
				r := Result{err: err}
				r.filename = filepath
				ch <- r
				<-tokens
				return
			}
			defer f.Close()
			r := inputProcess(f)
			r.filename = filepath
			ch <- r
			<-tokens
		}("test/" + filename.Name(), ch)

	}

	// Wait for all goroutine to finish and close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Receive data from channel
	for out := range ch {
		PrintOutput(out)
	}

}
