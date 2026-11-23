package main

import (
	"fmt"
	"net/http"
)

/*

labib0x9@TestGoCode $ go run main.go
Status= 200 OK
Error= Get "http://resolvefailed.com": dial tcp: lookup resolvefailed.com: no such host
*/

type Result struct {
	Error    error
	Response *http.Response
}

func main() {

	// ownership + error passsing
	getStatus := func(done <-chan struct{}, urls []string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, u := range urls {
				resp, err := http.Get(u)
				result := Result{Error: err, Response: resp}
				select { // pass result or shutdown..
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://google.com", "http://resolvefailed.com"}
	status := getStatus(done, urls)

	for result := range status {
		if result.Error != nil {
			fmt.Println("Error=", result.Error)
			continue
		}
		fmt.Println("Status=", result.Response.Status)
	}

}
