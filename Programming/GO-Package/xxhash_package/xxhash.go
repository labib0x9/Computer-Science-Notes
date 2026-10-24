package main

import (
	"fmt"
	"math/rand"

	"github.com/cespare/xxhash/v2"
)

func R() uint64 {
	return rand.Uint64()
}

func main() {

	h := xxhash.NewWithSeed(11091021361343833604)

	h.MarshalBinary()
	h.UnmarshalBinary()

	h.Write()
	h.BlockSize()

	h.WriteString("Heelo")
	debug(h.Sum64())

	h.ResetWithSeed(11091021361343833604)
	h.WriteString("Heelo")
	debug(h.Sum64())
}

func debug(n any) {
	fmt.Println("Here, ", n)
}
