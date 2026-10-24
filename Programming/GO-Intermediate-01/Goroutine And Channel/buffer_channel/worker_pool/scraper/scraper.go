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
var limiter = time.NewTicker(1000 * time.Millisecond)

func pool(n int, links []string, title chan string, wg *sync.WaitGroup) {
	jobs := make(chan string, len(links))
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(jobs, title, wg)
	}

	for _, url := range links {
		jobs <- url
	}
	close(jobs)
}

func worker(jobs chan string, title chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		<-limiter.C
		t := FindMetadata(job, true)
		if len(t) == 0 {
			title <- ""
		} else {
			title <- t[0]
		}
	}
}

func main() {

	defer limiter.Stop()

	url := "https://books.toscrape.com/index.html"
	links := FindMetadata(url, false)
	for i := range links {
		links[i] = "https://books.toscrape.com/" + links[i]
	}

	title := make(chan string, len(links))

	var wg sync.WaitGroup
	pool(4, links, title, &wg)

	go func() {
		wg.Wait()
		close(title)
	}()

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
