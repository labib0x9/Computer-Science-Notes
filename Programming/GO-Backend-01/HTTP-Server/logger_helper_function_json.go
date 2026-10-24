package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	loggerHelper(r)
	fmt.Fprintf(w, "API HANDLER")
}

func nopeHandler(w http.ResponseWriter, r *http.Request) {
	loggerHelper(r)
	fmt.Fprintf(w, "Nope Handler")
}

func loggerHelper(r *http.Request) {
	url := r.URL.String()
	reqType := r.Method
	protocol := r.Proto
	reqBodyLen := 0
	if reqType == "POST" {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Logger", "Url", url, "Request Type", reqType, "Body Length", reqBodyLen, "Protocol", protocol)
			return
		}
		reqBodyLen = len(reqBody)
	}
	slog.Info("Logger", "Url", url, "Request Type", reqType, "Body Length", reqBodyLen, "Protocol", protocol)
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api", apiHandler)
	mux.HandleFunc("POST /nope", nopeHandler)

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	log.Panic(http.ListenAndServe(addr, mux))
}
