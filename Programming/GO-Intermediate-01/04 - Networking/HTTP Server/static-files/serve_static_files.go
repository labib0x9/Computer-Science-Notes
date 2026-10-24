package main

import (
	"fmt"
	"net/http"
)

// logging
func printRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before", r.Method, r.URL.Path)
		next.ServeHTTP(w, r) // without this middleware doesn't work
		fmt.Println("After", r.Method, r.URL.Path)
	})
}

func main() {

	// // // http://127.0.0.1:8080/static.go -> ./static/static.go

	// // use /static as file loading
	root := http.Dir("./static")
	fs := http.FileServer(root) // root -> http.FileSystem

	// without -> http://host:port/ -> ./static
	// with ->    http://host:port/static/ -> ./static
	// Route: /static/* → ./static/*
	fsHandler := http.StripPrefix("/static/", fs)

	// http://127.0.0.1:8080/static/.env -> exposed
	// directory trversal -> handles correctly
	http.Handle("/", printRequest(fsHandler))

	http.ListenAndServe(":8080", nil)
}
