package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Get(w http.ResponseWriter, r *http.Request) {
	rawReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rawReq))
}

func main() {
	// --- //
}
