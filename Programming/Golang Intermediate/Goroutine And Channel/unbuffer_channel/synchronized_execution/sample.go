package main

import "fmt"

func main() {

	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func (n int)  {
			fmt.Println("Worker: ", i)
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}