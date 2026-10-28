package main

import (
	"fmt"
)

/*

labib0x9@TestGoCode $ go run main.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
        /Users/labib0x9/Desktop/TestGoCode/main.go:15 +0x98

goroutine 35 [chan receive]:
main.main.func1()
        /Users/labib0x9/Desktop/TestGoCode/main.go:12 +0x28
created by main.main in goroutine 1
        /Users/labib0x9/Desktop/TestGoCode/main.go:11 +0x88
exit status 2

*/

// Channels are typed.

func main() {
	intStream := make(chan int)
	stringStream := make(chan string)

	go func() {
		fmt.Println(<-intStream)
	}()

	stringStream <- "hello from channel"
	// time.Sleep(1 * time.Second)
}
