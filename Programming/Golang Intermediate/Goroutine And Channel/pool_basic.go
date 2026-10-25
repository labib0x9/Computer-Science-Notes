package main

/*

labib0x9@TestGoCode $ go run main.go
Init..
Init..

*/

import (
	"fmt"
	"sync"
)

func main() {
	mypool := &sync.Pool{
		New: func() any {
			fmt.Println("Init..")
			return struct{}{}
		},
	}

	mypool.Get()
	instance := mypool.Get()
	mypool.Put(instance)
	mypool.Get()
}
