package main

import (
	"io"
	"os"
	"strings"
)

// Incomplete

func main() {
	io.Copy(os.Stdout, os.Stdin)

	r1 := strings.NewReader("Labib")
	r2 := strings.NewReader("Faisal")
	r := io.MultiReader(r1, r2)

	w1 := os.Stdout
	w2 := os.Stderr
	w := io.MultiWriter(w1, w2)

	io.CopyN(w, r, 10)

	buf := make([]byte, 10)
	io.CopyBuffer(w, r, buf)

	_ = io.TeeReader(r, w)

	pr, pw := io.Pipe()
	
	pr.Close()
	pr.CloseWithError()
	pr.Read()

	pw.Close()
	pw.CloseWithError()
	pw.Write()
}
