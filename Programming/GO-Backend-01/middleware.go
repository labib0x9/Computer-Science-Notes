package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (m *Manager) Use(middlewares ...Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
}

func (m *Manager) With(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}

	for _, middleware := range m.globalMiddlewares {
		next = middleware(next)
	}

	return next
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func Info(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
	})
}

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]any{
		"test": "passed",
	})
}

func main() {
	manager := NewManager()
	manager.Use(Logger, Info) // But which order they executes ? [Loger, info] ? [Info, Logger] ?

	mux := http.NewServeMux()
	mux.Handle(
		"GET /api/notes",
		manager.With(
			http.HandlerFunc(GetAllNotes),
			// Others middleware for this specific path
		),
	)
}
