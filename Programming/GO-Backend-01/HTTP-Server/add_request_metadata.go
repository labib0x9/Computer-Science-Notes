package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type requestContextKey struct{}
type requestContextValue struct {
	requestId string
}

func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestId: requestID,
	}
	curCtx := r.Context()
	newCtx := context.WithValue(curCtx, requestContextKey{}, c)
	return r.WithContext(newCtx)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	reqID := "abc-123-xyz"
	r = addRequestID(r, reqID)
	loggerHelper(r)
	w.Header().Add("X-Request-ID", reqID)
	fmt.Fprintf(w, "API HANDLER")
}

func loggerHelper(r *http.Request) {
	url := r.URL.String()
	reqType := r.Method
	protocol := r.Proto
	reqBodyLen := 0
	var reqId requestContextValue
	v := r.Context().Value(requestContextKey{})
	if m, ok := v.(requestContextValue); ok {
		reqId = m
	}
	if reqType == "POST" {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Logger", "ID", reqId.requestId, "Url", url, "Request Type", reqType, "Body Length", reqBodyLen, "Protocol", protocol)
			return
		}
		reqBodyLen = len(reqBody)
	}
	slog.Info("Logger", "ID", reqId.requestId, "Url", url, "Request Type", reqType, "Body Length", reqBodyLen, "Protocol", protocol)
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api", apiHandler)

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	log.Panic(http.ListenAndServe(addr, mux))
}
