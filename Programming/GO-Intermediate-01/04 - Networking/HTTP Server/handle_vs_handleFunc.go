package main

import "net/http"

// what we want to do, when /api is called.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World\r\n"))
}

func main() {

	// Both does the same perpose
	// http.Handle registers object
	// http.HandleFunc registers function

	// handles /api path
	http.HandleFunc("/api", apiHandler)

	// handles /api/ path
	http.Handle("/api/", http.HandlerFunc(apiHandler))

	http.ListenAndServe(":8080", nil)
}
