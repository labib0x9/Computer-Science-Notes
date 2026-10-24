// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func checkLink(link string, ch chan string) {
// 	_, err := http.Get(link)
// 	if err != nil {
// 		// fmt.Println(link, "is Down")
// 		ch <- "not ok"
// 		return
// 	}
// 	ch <- "ok"
// }

// func main() {

// 	links := []string{
// 		"https://google.com",
// 		"https://facb.com",
// 		"https://go.dev",
// 	}

// 	ch := make(chan string)

// 	for _, link := range links {
// 		go checkLink(link, ch)
// 	}

// 	for range links {
// 		select {
// 		case res := <-ch:
// 			fmt.Println(res)
// 		case <-time.After(2 * time.Second):
// 			fmt.Println("Error")
// 		}
// 	}
// }

package main

import (
	"fmt"
	"time"
)

// only Send
func producer(ch chan<- int) {

	for i := 1; i <= 5; i++ {
		fmt.Println("Producing : ", i)
		ch <- i
		fmt.Println("Produced")
		time.Sleep(time.Millisecond * 200)
	}
	fmt.Println("Producing done")

	// Run without close() and with close()
	close(ch)
}

// only receive
func consumer(ch <-chan int) {
	time.Sleep(time.Second * 5)
	for val := range ch {
		fmt.Println("Consumed: ", val)
		time.Sleep(time.Second * 5)
	}

	_, ok := <-ch
	if !ok {
		fmt.Println("Channel Closed", ok)
	}

	fmt.Println("Consuming done")
}

func main() {

	ch := make(chan int, 2)

	go producer(ch)
	go consumer(ch)

	time.Sleep(35 * time.Second)
}
