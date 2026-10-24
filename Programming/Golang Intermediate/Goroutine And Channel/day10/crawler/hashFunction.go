package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

var contestArchivePage = "https://atcoder.jp/contests/archive"

func removeNonContestLink(links []string) string {
	var contestIds string
	for _, link := range links {
		link, found := strings.CutPrefix(link, contestPage)
		if found {
			if len(contestIds) == 0 {
				contestIds = link
			} else {
				contestIds = contestIds + "+" + link
			}
		}
	}
	return contestIds
}

func generateHash(contestIds string) string {
	hashTemp := sha256.Sum256([]byte(contestIds))
	hash := hex.EncodeToString(hashTemp[:])
	return hash
}

func getHash() string {

	links := FindLinks(contestArchivePage)
	contestIds := removeNonContestLink(links)

	fmt.Println(contestIds == "")

	return generateHash(contestIds)
}
