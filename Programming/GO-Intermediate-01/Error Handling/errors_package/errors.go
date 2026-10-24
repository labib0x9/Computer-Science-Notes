package main

import "errors"

func main() {

	err1 := errors.New("file1 cannot be opened")
	err2 := errors.New("file2 cannot be closed")
	errors.Join(err1, err2)
}