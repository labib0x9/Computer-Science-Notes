package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// What is the perpose of middleware ?
func handleMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Response Header
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Random-Header", "true")

		// Set Request Header
		r.Header.Set("X-Hackerone-Id", "9xzer0")

		slog.Info(
			"Incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote-addr", r.RemoteAddr,
		)

		//
		next.ServeHTTP(w, r)

		// dbQuery, external service, etc etc
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Why ?
	// Try http://localhost:8080/test with and without this if condition
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to homepage")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	//
	msg := map[string]string{
		"title": "Testing",
		"msg":   "Contetnt not found",
	}

	// use json encoder
	json.NewEncoder(w).Encode(msg)
}

func main() {

	http.HandleFunc("/", homeHandler)

	// Go through middleware
	http.Handle("/echo", handleMiddleware(http.HandlerFunc(jsonHandler)))

	// Start server
	fmt.Println("Starting server at 127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)

	// Can we direct two different handler with one path.
	// like homeHandler, jsonHandler in both "/" path
}
