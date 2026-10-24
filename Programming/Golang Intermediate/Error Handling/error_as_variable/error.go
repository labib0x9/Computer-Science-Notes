package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	NotFoundError = errors.New("Not found")
	FileOpenFailed = errors.New("File open failed")
	ErrOutOfIndex = errors.New("Slice is too short to access")
)

func main() {

	var err error
	_, err = os.Open("a.txt")

	// Handle error like variable
	if err != nil {

	}

	err = fmt.Errorf("%w : ", FileOpenFailed)

}