package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var mu sync.Mutex
var requestCount int

func FanOut(n int, links <-chan string, rateLimiter *time.Ticker) []<-chan string {
	titles := make([]<-chan string, n)

	for i := 0; i < n; i++ {
		titles[i] = worker(i, links, rateLimiter)
	}
	return titles
}

func FanIn(titles []<-chan string) <-chan string {
	var wg sync.WaitGroup
	title := make(chan string)

	for _, outChan := range titles {
		wg.Add(1)
		go func(out <-chan string) {
			defer wg.Done()
			for r := range out {
				title <- r
			}
		}(outChan)
	}

	go func() {
		wg.Wait()
		close(title)
	}()
	return title
}

func worker(n int, links <-chan string, rateLimiter *time.Ticker) chan string {
	title := make(chan string)
	go func() {
		defer close(title)
		// 1 http request per tick, not per worker, any worker may put the request
		for url := range links {
			<-rateLimiter.C // wait
			t := FindMetadata(url, true)
			var res string
			if len(t) != 0 {
				res = t[0]
			}
			title <- res
		}
	}()
	return title
}

func main() {

	input := make(chan string, 3)
	go func() {
		defer close(input)
		url := "https://books.toscrape.com/index.html"
		links := FindMetadata(url, false)
		for i := range links {
			input <- "https://books.toscrape.com/" + links[i]
		}
	}()

	rateLimiter := time.NewTicker(300 * time.Millisecond)
	defer rateLimiter.Stop()

	chanOut := FanOut(3, input, rateLimiter)
	title := FanIn(chanOut)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()
			fmt.Println("Requests/sec:", requestCount)
			requestCount = 0
			mu.Unlock()
		}
	}()

	for t := range title {
		fmt.Println(t)
	}
}

func FindMetadata(url string, title bool) (res []string) {

	mu.Lock()
	requestCount++
	mu.Unlock()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ERROR")
		return
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal("Error parsing body: ", url)
		return
	}

	extract(&res, doc, title)
	return
}

func extract(res *[]string, n *html.Node, title bool) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if title {
				if n.FirstChild != nil {
					*res = append(*res, n.FirstChild.Data)
				}
				return
			}
			fallthrough
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					*res = append(*res, a.Val)
				}
			}
		}
	}

	// Recurse on children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extract(res, c, title)
	}
}
