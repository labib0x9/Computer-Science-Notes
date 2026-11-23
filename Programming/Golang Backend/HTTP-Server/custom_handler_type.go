package main

import "net/http"

/*

sharing data between handler functions

*/

type appConfig struct {
	level string
}

// implements http.Handler interface, so can be used as mux
type app struct {
	config  *appConfig
	handler func(w http.ResponseWriter, r *http.Request, config *appConfig)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

// custom handler
func apiHandler(w http.ResponseWriter, r *http.Request, config *appConfig) {
	w.Write([]byte(config.level))
}

func main() {
	config := appConfig{level: "ERROR"}

	mux := http.NewServeMux()
	mux.Handle("GET /api", &app{config: &config, handler: apiHandler})

	http.ListenAndServe(":8080", mux)
}
