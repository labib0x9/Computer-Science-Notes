package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

/*

labib0x9@TestGoCode $ curl -i http://localhost:8080/decode -d '{"user_ip": "127.0.0.1", "event": "dgg", "user_data": "my personal information"}'
HTTP/1.1 200 OK
Date: Wed, 21 Oct 2026 19:10:23 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

OK%

labib0x9@TestGoCode $ go run main.go
127.0.0.1 dgg



labib0x9@TestGoCode $ curl -i http://localhost:8080/decode -d '{"user_ip": "127.0.0.1", "event": "dgg", "user_data": "my personal information"}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 21 Oct 2026 19:12:59 GMT
Content-Length: 32

json: unknown field "user_data"

*/

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// See commented and uncommented...
	decoder.DisallowUnknownFields()

	for {
		var l logLine
		err := decoder.Decode(&l)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(l.UserIP, l.Event)
	}
	fmt.Fprintf(w, "OK")
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /decode", decodeHandler)

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	log.Fatal(http.ListenAndServe(addr, mux))
}
