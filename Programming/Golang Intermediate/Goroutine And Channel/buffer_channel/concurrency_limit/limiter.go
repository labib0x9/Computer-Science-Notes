package main

func main() {

	// maximum 3 goroutine running concurrently
	token := make(chan struct{}, 3)

	for i := 0; i < 15; i++ {
		token <- struct{}{}
		go func(n int) {
			// do work
		}(i)
		<-token
	}

}
