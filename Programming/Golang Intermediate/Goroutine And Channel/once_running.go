package main

/*
labib0x9@TestGoCode $ go run main.go
A

Why prints only A ?
*/

import (
	"fmt"
	"sync"
)

func main() {
	fn1 := func() { fmt.Println("A") }
	fn2 := func() { fmt.Println("B") }

	var once sync.Once

	once.Do(fn1)
	once.Do(fn2)
}
