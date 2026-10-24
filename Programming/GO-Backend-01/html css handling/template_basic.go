package main

import (
	"html/template"
	"net/http"
)

// parse the template file once..
var templ = template.Must(template.ParseFiles("./hello.html"))

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// get query
	q := r.URL.Query().Get("name")

	// html format, no encoding
	data := template.HTML(q)

	// fmt.Println(data)

	// 
	templ.Execute(w, map[string]any{
		"Name": data,
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("http://ip:port/echo?name=Labib"))
}

func main() {

	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":8080", nil)
}
