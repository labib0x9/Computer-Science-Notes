package main

import "fmt"

/*

Is it a goroutine leak ?? how can i tell the gorutine for cancellation ?

labib0x9@TestGoCode $ go run main.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
        /Users/labib0x9/Desktop/TestGoCode/main.go:22 +0x38

goroutine 19 [chan receive (nil chan)]:
main.main.func1.1()
        /Users/labib0x9/Desktop/TestGoCode/main.go:13 +0x90
created by main.main.func1 in goroutine 1
        /Users/labib0x9/Desktop/TestGoCode/main.go:11 +0x80
exit status 2

*/

func main() {
	chanOwner := func(stream <-chan int) <-chan struct{} {
		complete := make(chan struct{})
		go func() {
			defer close(complete)
			for val := range stream { // ranging over a nil channel, so defer never executes, blocks on READ operation.
				fmt.Println(val)
			}
		}()
		return complete
	}

	// stream := make(chan int, 4)
	complete := chanOwner(nil)
	<-complete
}
