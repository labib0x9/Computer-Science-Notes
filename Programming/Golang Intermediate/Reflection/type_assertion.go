package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {

	var r io.Reader
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	r = file // r = {type: *os.File, value: 0xc0000a2000}

	fmt.Println("Type of file:", reflect.TypeOf(file)) // *os.File
	fmt.Println("Type of r:", reflect.TypeOf(r))       // *os.File

	// Both file and r point to the same *os.File value
	// What i did, i just copied value from one interface to another interface
	fmt.Printf("Address: %v\n", *file)                               // 0xc0000a2000
	fmt.Printf("Address: %v\n", reflect.ValueOf(file).Elem().Addr()) // 0xc0000a2000
	fmt.Printf("Address: %v\n", reflect.ValueOf(r).Elem().Addr())    // 0xc0000a2000

	var w io.Writer
	w = r.(io.Writer) // Now w is assigned the value of r
	_ = w             // w = {type: *os.File, value: 0xc0000a2000}

	w.Write([]byte("Hello, World!"))

	var e interface{}
	e = w // e = {type: *os.File, value: 0xc0000a2000}
	_ = e
}
