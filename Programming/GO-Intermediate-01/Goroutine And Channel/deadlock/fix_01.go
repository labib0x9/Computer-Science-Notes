package main

import (
	"fmt"
	"sync"
	"time"
)

/*

labib0x9@TestGoCode $ go build -race -o out main.go
labib0x9@TestGoCode $ ./out
ID= b  | a + b = 0
ID= a  | a + b = 0

*/

type value struct {
	mu    sync.Mutex
	value int
}

func main() {

	var wg sync.WaitGroup
	printSum := func(id string, v1, v2 *value) {
		// Why this ordering fix deadlock ???
		if id == "b" {
			v1, v2 = v2, v1
		}
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(20 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Println("ID=", id, " | a + b =", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)

	go printSum("a", &a, &b)
	go printSum("b", &b, &a)

	wg.Wait()
}
