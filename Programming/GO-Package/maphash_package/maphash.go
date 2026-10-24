package main

import (
	"fmt"
	"hash/maphash"
)

func main() {

	// same seed, same hash
	// different seed, different hash
	// Important :: You can't store seed for next use
	seed := maphash.MakeSeed()

	// create a new maphash.Hash
	var h maphash.Hash
	h.SetSeed(seed)

	// Hash
	h.Write([]byte("hello"))
	fmt.Println(h.Sum64())

	h.Reset()
	h.Write([]byte("world"))
	fmt.Println(h.Sum64())

	h.Reset()
	h.WriteString("hello")
	h.WriteString(" ")
	h.Write([]byte("world"))

	h.Sum64() // hash for "hello world"
}
