package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/*


labib0x9@TestGoCode $  curl -i 'http://localhost:8080/download?filename=go.mod'
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Fri, 23 Oct 2026 17:38:02 GMT
Transfer-Encoding: chunked

module test

go 1.26.2

*/

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	fmt.Println(filename)

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil {
		return
	}
	contentType := http.DetectContentType(buf[:n])
	w.Header().Set("Content-Type", contentType)

	f.Seek(0, 0)

	io.Copy(w, f)
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /download", downloadHandler)

	log.Fatal(http.ListenAndServe(addr, mux))
}
