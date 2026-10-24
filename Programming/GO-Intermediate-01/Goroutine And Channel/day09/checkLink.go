package main

import (
	"fmt"
	"net/http"
)

// Check if servers are up or down

func CheckLinks() {

	links := []string{
		"https://google.com",
		"https://facb.com",
		"https://go.dev",
	}

	ch := make(chan string)

	for _, link := range links {
		go checkLink(link, ch)
		// fmt.Println(<-ch)
	}

	// for i := 0; i < len(links); i++ {
	// 	fmt.Println(<-ch)
	// }

	// for {
	// 	go checkLink(<-ch, ch)
	// }

	// After receiving all channel it will still wait for receiving data
	// which will cause deadlock
	for link := range ch {
		fmt.Println(link)
	}
}

func checkLink(link string, ch chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "is Down")
		// ch <- "not ok"
		ch <- link
		return
	}
	fmt.Println(link, "is Up")
	// ch <- "ok"
	ch <- link
}
