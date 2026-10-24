package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// http.FileSystem is an interface, with only `Open(name string) (File, error)` method
type wrapFS struct {
	fh http.FileSystem
}

func (w wrapFS) Open(name string) (http.File, error) {
	clean := path.Clean(name)
	fmt.Println("WrapFS=", clean)
	base := path.Base(clean)
	if strings.HasPrefix(base, ".") == true {
		return nil, errors.New("No no, you can't see.")
	}
	return w.fh.Open(clean)
}

// logging
func printRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming=", r.Method, r.URL.Path)
		next.ServeHTTP(w, r) // without this middleware doesn't work
	})
}

func main() {

	// use /static as file loading
	// without wrapFS http://127.0.0.1:8080/static/.env -> 200 OK
	// with wrapFS http://127.0.0.1:8080/static/.env -> 500 Internal Server Error

	// root is a http.Dir which has Open() method, so root implements http.FileSystem
	root := http.Dir("./static")

	// wrapFs stores http.FileSystem, wfs stores the behavior of http.Dir (root).
	wfs := wrapFS{fh: root}

	// http.FileServer takes a http.FileSystem
	// But why wfs ? WrapFs
	// because WrapFS has a method named Open(), which implements http.FileSystem.
	fs := http.FileServer(wfs)


	/*
		Function Interception / Structural Typing

		FileServer calls filesystem.Open()	-> filesystem (interface) is now WrapFS
		filesystem calls wfs.Open()		-> WrapFS
		wfs.Open() calls root.Open()	-> http.Dir
	*/


	fsHandler := http.StripPrefix("/static/", fs)

	http.Handle("/", printRequest(fsHandler))

	http.ListenAndServe(":8080", nil)
}
