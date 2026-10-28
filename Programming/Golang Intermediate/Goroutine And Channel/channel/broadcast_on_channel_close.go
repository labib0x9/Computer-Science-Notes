package main

/*

labib0x9@TestGoCode $ go run main.go
Unblocking goroutines !!
10 has begun!
15 has begun!
14 has begun!
13 has begun!
12 has begun!
11 has begun!
19 has begun!
20 has begun!
18 has begun!
17 has begun!
16 has begun!
labib0x9@TestGoCode $

*/

import (
	"fmt"
	"sync"
)

func main() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 10; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun!\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines !!")
	close(begin)
	wg.Wait()
}
