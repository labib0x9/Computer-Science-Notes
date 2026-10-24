package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestcheckInput(t *testing.T) {

	reader := bufio.NewReader(strings.NewReader("5\n"))

	message, _ := checkInput(reader)

	_ = message
}

func TestreadUserInput(t *testing.T) {
	
	doneChan := make(chan bool)
	
	var stdin bytes.Buffer
	stdin.Write([]byte("5\nq\n"))

	go readUserInput(&stdin, doneChan)
	<-doneChan
}