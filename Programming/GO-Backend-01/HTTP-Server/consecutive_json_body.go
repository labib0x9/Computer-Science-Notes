package main

/*

labib0x9@TestGoCode $ curl -i http://localhost:8080/api -d '{"ip": "127.0.0.1", "event": "dgg"}'
HTTP/1.1 200 OK
Date: Wed, 21 Oct 2026 18:13:35 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

API HANDLER%


labib0x9@TestGoCode $ curl -i http://localhost:8080/api -d '{"ip": "127.0.0.1", "event": "dgg"} {"ip": "127.0.0.2", "event": "hello i am double"}'
HTTP/1.1 200 OK
Date: Wed, 21 Oct 2026 18:15:18 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

API HANDLER%



labib0x9@TestGoCode $ go run main.go
{"time":"2026-10-22T00:13:35.834917+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.1","Event":"dgg"}
{"time":"2026-10-22T00:15:18.091179+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.1","Event":"dgg"}

*/

/*
labib0x9@TestGoCode $ curl -i http://localhost:8080/api2 -d '{"ip": "127.0.0.1", "event": "dgg"} {"ip": "127.0.0.2", "event": "hello i am double"}'
HTTP/1.1 200 OK
Date: Wed, 21 Oct 2026 18:36:09 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

API HANDLER%

labib0x9@TestGoCode $ go run main.go
{"time":"2026-10-22T00:36:09.471028+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.1","Event":"dgg"}
{"time":"2026-10-22T00:36:09.471062+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.2","Event":"hello i am double"}
*/

/*
labib0x9@TestGoCode $ curl -i http://localhost:8080/api2 -d '{"ip": "127.0.0.1", "event": "dgg"}} {"ip": "127.0.0.2", "event": "hello i am double"}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 21 Oct 2026 18:45:47 GMT
Content-Length: 12

BAD REQUEST

labib0x9@TestGoCode $ go run main.go
{"time":"2026-10-22T00:45:47.505674+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.1","Event":"dgg"}
*/

/*

labib0x9@TestGoCode $ curl -i http://localhost:8080/api3 -d '{"ip": "127.0.0.1", "event": "dgg"}} {"ip": "127.0.0.2", "event": "hello i am double"}'
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 21 Oct 2026 18:54:43 GMT
Content-Length: 12

BAD REQUEST

labib0x9@TestGoCode $ go run main.go
{"time":"2026-10-22T00:54:43.873591+06:00","level":"INFO","msg":"Handler","IP":"127.0.0.1","Event":"dgg"}
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type ReqType struct {
	Ip    string `json:"ip"`
	Event string `json:"event"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var req ReqType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	slog.Info("Handler", "IP", req.Ip, "Event", req.Event)

	fmt.Fprintf(w, "API HANDLER")
}

func api2Handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	for {
		var req ReqType
		err := decoder.Decode(&req)
		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, "BAD REQUEST", http.StatusBadRequest)
			return
		}

		slog.Info("Handler", "IP", req.Ip, "Event", req.Event)
	}

	fmt.Fprintf(w, "API HANDLER")
}

func api3Handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e *json.UnmarshalTypeError
	for {
		var req ReqType
		err := decoder.Decode(&req)
		if err == io.EOF {
			break
		}

		if errors.As(err, &e) {
			slog.Error("Unmarshal errror", "err", err)
			continue
		}

		if err != nil {
			http.Error(w, "BAD REQUEST", http.StatusBadRequest)
			return
		}

		slog.Info("Handler", "IP", req.Ip, "Event", req.Event)
	}

	fmt.Fprintf(w, "API HANDLER")
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", apiHandler)
	mux.HandleFunc("POST /api2", api2Handler)
	mux.HandleFunc("POST /api3", api3Handler)

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	log.Fatal(http.ListenAndServe(addr, mux))
}
