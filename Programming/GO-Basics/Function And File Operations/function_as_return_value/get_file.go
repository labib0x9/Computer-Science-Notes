package main

import (
	"io"
	"log"
	"os"
)

func getFile(path string) (*os.File, func(), error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		file.Close()
	}, nil
}

func main() {

	file, closer, err := getFile("file.txt")
	if err != nil {
		log.Panic(err)
	}

	io.Copy(os.Stdout, file)
	defer closer()

}