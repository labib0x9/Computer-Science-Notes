package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// Serialization format
// Incomplete

type User struct {
	Name string
}

func main() {

	var buf bytes.Buffer

	en := gob.NewEncoder(&buf) // writes to buf
	de := gob.NewDecoder(&buf) // reads from buf

	// When to use Register ?
	// Encoding/decoding interface{}
	// Interface field inside a struct
	// gob.Register(User{})	// No needed

	// struct, map, slices dont need register

	// Encode Data
	data := User{Name: "Labib"}
	en.Encode(data)

	// In-Memory cache
	cache := make(map[string][]byte)
	cache["Labib"] = buf.Bytes()

	// Decode Data
	var deData User
	de.Decode(&deData)

	// Files
	// Important fact: while reading and writing in the same function
	// use seek to reset pointer, else the pointer will be a EOF
	f, _ := os.Create("data.gob")
	defer f.Close()

	en = gob.NewEncoder(f)

	// en.EncodeValue()
	// de.DecodeValue()

}

func File() {
	filename := "test.gob"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	if err := enc.Encode(10); err != nil {
		return
	}
	dec := gob.NewDecoder(file)

	// Run with and without seek
	// with seek the file pointer is at 0
	// without seek the file pointer is at EOF
	file.Seek(0, io.SeekStart)
	var n int
	if err := dec.Decode(&n); err != nil {
		if err == io.EOF {
			fmt.Println("I am at EOF, set pointer to offset 0")
		} else {
			return
		}
	}

	fmt.Println(n)
}

// Use Cases:

//     Fast IPC (Inter-Process Communication) between Go services.

//     Local persistence of Go structs (e.g., caching to file).

//     RPC serialization using net/rpc.
// Send serialized structs via channels for logging, processing, or queuing systems.
