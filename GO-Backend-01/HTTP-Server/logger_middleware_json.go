package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Practical Go Book - exercise 5.1 a request logger

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API HANDLER")
}

func nopeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Nope Handler")
}

func loggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		reqType := r.Method
		reqBodyLen := 0
		if reqType == "POST" {
			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Error...")
				return
			}
			reqBodyLen = len(reqBody)
		}
		protocol := r.Proto
		slog.Info("Logger", "Url", url, "Request Type", reqType, "Body Length", reqBodyLen, "Protocol", protocol)
		h.ServeHTTP(w, r)
	})
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.Handle("GET /api", loggerMiddleware(http.HandlerFunc(apiHandler)))
	mux.Handle("POST /nope", loggerMiddleware(http.HandlerFunc(nopeHandler)))

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	log.Panic(http.ListenAndServe(addr, mux))
}
