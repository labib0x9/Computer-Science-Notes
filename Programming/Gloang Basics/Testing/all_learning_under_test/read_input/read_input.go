package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func isPrime(n int) (bool, string) {
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not a prime", n)
	}
	if n == 2 {
		return true, fmt.Sprintf("%d is a prime", n)
	}
	if n%2 == 0 {
		return false, fmt.Sprintf("%d is not a prime, divisible by 2", n)
	}

	// i * i <= n
	// i <= n
	for i := 3; i <= n; i += 2 {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not a prime, divisible by %d", n, i)
		}
	}
	return true, fmt.Sprintf("%d is a prime", n)
}

func readUserInput(in io.Reader ,doneChan chan bool) {
	reader := bufio.NewReader(in)
	for {
		message, done := checkInput(reader)
		if done {
			doneChan <- true
		}

		fmt.Println(message)
	}
}

func checkInput(reader *bufio.Reader) (string, bool) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", false
	}
	line = strings.TrimSpace(line)
	if line == "q" {
		return "", true
	}

	n, err := strconv.Atoi(line)
	if err != nil {
		return fmt.Sprintf("%s is not a valid number", line), false
	}

	_, msg := isPrime(n)
	return msg, false
}

func main() {

	done := make(chan bool, 1)
	go readUserInput(os.Stdin, done)
	<-done
}
