package main

/*

labib0x9@TestGoCode $ go run main.go
Created = 8

*/

import (
	"fmt"
	"sync"
)

func main() {
	numObjCreation := 0
	mypool := &sync.Pool{
		New: func() any {
			numObjCreation++
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// seed
	mypool.Put(mypool.New())
	mypool.Put(mypool.New())
	mypool.Put(mypool.New())
	mypool.Put(mypool.New())

	const worker = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(worker)

	for i := 0; i < worker; i++ {
		go func() {
			defer wg.Done()
			mem := mypool.Get().(*[]byte) // type assert
			defer mypool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Println("Created =", numObjCreation)
}
