package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func FindLinks(url string) (urls []string) {
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

	extractUrl(&urls, doc)
	return
}

func extractUrl(links *[]string, n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					*links = append(*links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractUrl(links, c)
	}
}