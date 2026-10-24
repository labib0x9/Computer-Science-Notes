package main

var token = make(chan struct{}, 2)

// semaphore mechanism

// What is happening here, buffer channel details

// Code - 01

func crawl1(url string, ch chan []string) {
	ch <- FindLinks(url) // Send to ch channel
	<-token              // Release a token
}

func code1() {
	seen := make(map[string]bool)
	ch := make(chan []string)
	n := 1

	token <- struct{}{} // acquire a token
	seen["http://go.dev"] = true
	go crawl1("http://go.dev", ch)

	for ; n > 0; n-- {
		word := <-ch // Wait to receive
		for _, link := range word {
			if seen[link] {
				continue
			}
			seen[link] = true
			n++
			token <- struct{}{} // acquire a token, Blocks if the token is full
			go crawl1(link, ch)
		}
	}
}

// Token should be acquired and released in the same goroutine, Why ??

// Code - 02

func crawl2(url string, ch chan []string) {
	token <- struct{}{}
	ch <- FindLinks(url)
	<-token
}

func code2() {
	seen := make(map[string]bool)
	// ch := make(chan []string)
	ch := make(chan []string, 100)
	n := 1

	seen["http://go.dev"] = true
	go crawl2("http://go.dev", ch)

	for ; n > 0; n-- {
		word := <-ch
		for _, link := range word {
			if seen[link] {
				continue
			}
			seen[link] = true
			n++
			go crawl2(link, ch)
		}
	}
}

// Code - 03

func code3() {
	
}

// how to close this channels and waitgroups -> Code3

func FindLinks(url string) []string {
	return []string{}
}
