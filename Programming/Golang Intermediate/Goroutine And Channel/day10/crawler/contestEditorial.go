package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

var tags = []string{
	"tree",
	"depth-first search",
	"dfs",
	"bfs",
	"graph",
	"dijkstra",
	"shortest path",
	"hashing",
	"warshall-floyd",
	"strongly connected component",
	"bipartite graph",
	"maxflow",
	"mincostflow",
	"atcoder::mcf_graph",

	"disjoint det union",
	"dsu",

	"cumulative sum",
	"prefix sum",
	"prefix-sum",

	"priority_queue",
	"stack",
	"queue",

	"square root decomposition",
	"segment tree",
	"fenwick tree",
	"lazy segment tree",
	"lca",
	"euler tour",
	"virtual tree",
	"heavy-light decomposition",
	"centroid decomposition",
	"binary lifting",

	"dynamic programming",
	"knapsack",
	"digit dp",
	"subset dp",
	"inline dp",

	"data structure",
	"algorithm",

	"sliding window",
	"binary search",
	"upperbound",
	"lowerbound",

	"next_permutation",

	"atcoder::convolution",
	"fast fourier transform",
}

func extractTags(url string) (tag []string) {
	// (*http.Request, error)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : GET REQ")
		return
	}

	// Set User-Agent (Custom Header)
	// mu.Lock()
	req.Header.Set("User-Agent", userAgent())
	// mu.Unlock()

	// Make a client, Use pointer
	client := &http.Client{}

	// (*http.Response, error)
	resp, err := client.Do(req)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : DO REQ ", url)
		return
	}

	defer resp.Body.Close()

	// ([]byte, error)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// ERR = err
		log.Println("ERRERER : READ BODY")
		return
	}

	bodyStr := strings.ToLower(string(body))
	for _, key := range tags {
		if strings.Contains(bodyStr, key) {
			tag = append(tag, key)
		}
	}
	return
}

func completeContest() {

	for i := range contest {

		contest[i].contestType = (strings.Split(contest[i].contestUrl, "/")[4])[:3]

		if len(contest[i].problemUrl) == 0 {
			for _, link := range FindLinks(contest[i].contestUrl) {
				contestId := strings.Split(link, "/")[4]
				if strings.Count(link, contestId) == 2 {
					contest[i].problemUrl = append(contest[i].problemUrl, ProblemId{problemUrl: link})
				}
			}
		}

		if contest[i].problemUrl[0].editorialUrl == "" {
			for j := range contest[i].problemUrl {
				for _, link := range FindLinks(contest[i].problemUrl[j].problemUrl) {
					if strings.Contains(link, "editorial") {
						contest[i].problemUrl[j].editorialUrl = link
						contest[i].problemUrl[j].tags = extractTags(link)
					}
				}
			}
		}
	}

}
