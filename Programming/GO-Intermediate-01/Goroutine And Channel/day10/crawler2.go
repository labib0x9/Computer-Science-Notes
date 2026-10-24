package main

import (
	"fmt"
)

// Need to check, acquire token mechanism

// Enforce a limit of 20 concurrent request at once
var tokens = make(chan struct{}, 20)

func crawler2() {
	wordlist := make(chan []string)
	var n int

	n++
	tokens <- struct{}{}
	go Push2("http://gopl.io", wordlist)

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		lists := <-wordlist
		for _, link := range lists {
			if seen[link] {
				continue
			}
			seen[link] = true
			n++
			tokens <- struct{}{} // Acquire a token
			go Push2(link, wordlist)
		}
	}
}

func Push2(link string, wordlist chan []string) {
	wordlist <- crawl2(link)
}

func crawl2(url string) []string {
	fmt.Println(url)
	// tokens <- struct{}{}	// Acquire a token
	lists := FindLinks(url)
	<-tokens // Release a token
	return lists
}
