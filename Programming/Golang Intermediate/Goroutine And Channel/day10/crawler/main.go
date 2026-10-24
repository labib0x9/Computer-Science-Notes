package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)
// https://books.toscrape.com/index.html


/**
	For simplicity, we will only extract links from <a> tags
**/

type ProblemId struct {
	problemUrl   string
	editorialUrl string
	tags         []string
}

type ContestId struct {
	contestType string
	contestUrl  string
	problemUrl  []ProblemId
}

var BaseUrl = "https://atcoder.jp"
var token = make(chan struct{}, 5)
var ERR error
var mu sync.Mutex
var contestPage = BaseUrl + "/contests/"
var contest []ContestId
var contestPageHash string
var contestFound map[string]struct{}

func FindLinks(url string) (links []string) {

	// (*http.Request, error)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : GET REQ")
		return
	}

	// Set User-Agent (Custom Header)
	mu.Lock()
	req.Header.Set("User-Agent", userAgent())
	mu.Unlock()

	// Make a client, Use pointer
	client := &http.Client{}

	// (*http.Response, error)
	resp, err := client.Do(req)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : DO REQ ", url)
		return
	}

	// fmt.Println(resp.StatusCode)

	defer resp.Body.Close()

	// ([]byte, error)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : READ BODY")
		return
	}

	// (*html.Node, error)
	content, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		// ERR = err
		log.Println("ERRERER : PARSE BODY")
		return
	}

	Traverse(content, &links)
	return
}

// Dfs search
func Traverse(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}
			link := attr.Val
			if strings.HasPrefix(link, "//") {
				link = "https:" + link
			} else if strings.HasPrefix(link, "https://") {
				// Do Nothing
			} else if strings.HasPrefix(link, "http://") {
				// Do Nothing
			} else {
				if !strings.HasPrefix(link, "/") {
					link = "/" + link
				}
				link = BaseUrl + link
			}

			// Check if the link is in the same domain
			if strings.HasPrefix(link, BaseUrl) {
				*links = append(*links, link)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Traverse(c, links)
	}
}

func crawl(url string, ch chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	mu.Unlock()
	token <- struct{}{}
	fmt.Println(url)
	if strings.HasPrefix(url, contestPage) && strings.Contains(url, "tasks") {
		contest = append(contest, ContestId{contestUrl: url})
		contestFound[url] = struct{}{}
		ch <- []string{}
	} else {
		ch <- FindLinks(url)
	}
	<-token
}

func validateUrl(url string) bool {
	// fmt.Println(url, contestPage)
	return strings.HasPrefix(url, contestPage)
}

func bfs(url string) {
	var wg sync.WaitGroup
	seen := make(map[string]bool)
	ch := make(chan []string)
	n := 1

	seen[url] = true
	wg.Add(1)
	go crawl(url, ch, &wg)

	// Receieve from goroutine via channel
	// and send to goroutine again

	for ; n > 0; n-- {
		wordlist := <-ch
		for _, link := range wordlist {
			if seen[link] {
				continue
			}
			if validateUrl(link) {
				n++
				wg.Add(1)
				go crawl(link, ch, &wg)
			}
			seen[link] = true
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func getNewContest() {
	links := FindLinks(contestArchivePage)
	for _, url := range links {
		if strings.HasPrefix(url, contestPage) && strings.Contains(url, "tasks") {
			contest = append(contest, ContestId{contestUrl: url})
			contestFound[url] = struct{}{}
		}
	}
}

func main() {
	bfs("https://atcoder.jp/contests/archive")
	// completeContest()

	// contestPageHash = getHash()
	// This process is a continious process
	// if time.Now().Weekday().String() == "Thursday" {
	// 	tempHash := getHash()
	// 	if tempHash != contestPageHash {
	// 		contestPageHash = tempHash
	// 		getNewContest()
	// 		completeContest()
	// 	}
	// }

	// fmt.Println(getHash())
}
