package main

import "fmt"

// What is the problem in this code ??

func Crawler1() {

	wordlist := make(chan []string)

	// go func() {
	// 	wordlist <- []string{"http://gopl.io"}
	// }()
	go Push("http://gopl.io", wordlist)

	seen := make(map[string]bool)
	for list := range wordlist {
		for _, link := range list {
			if seen[link] {
				continue
			}
			seen[link] = true
			// go func(link string) {
			// 	wordlist <- crawl(link)
			// }(link)
			go Push(link, wordlist)
		}
	}
}

func Push(link string, wordlist chan []string) {
	wordlist <- crawl(link)
}

func crawl(link string) []string {
	fmt.Println(link)
	list := FindLinks(link)
	return list
}




















// func crawl(url string) []string {
// 	fmt.Println(url)
// 	list, err := Extract(url)
// 	if err != nil {
// 		return []string{}
// 	}
// 	return list
// }
