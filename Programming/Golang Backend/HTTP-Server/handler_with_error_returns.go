package main

import (
	"fmt"
	"net/http"
)

/*

sharing data between handler functions

*/

type appConfig struct {
	level string
}

// implements http.Handler interface, so can be used as mux
type app struct {
	config  *appConfig
	handler func(w http.ResponseWriter, r *http.Request, config *appConfig) error
}

// its just works like a middleware where handler returns a error
func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := a.handler(w, r, a.config)
	if err != nil {
		w.Write([]byte("Found" + err.Error()))
	} else {
		w.Write([]byte("No Errorr"))
	}
}

// custom handler, see it returns a error now :)
func apiHandler(w http.ResponseWriter, r *http.Request, config *appConfig) error {
	return fmt.Errorf("Error found: Log Level = %v", config.level)
}

func main() {
	config := appConfig{level: "ERROR"}

	mux := http.NewServeMux()
	mux.Handle("GET /api", &app{config: &config, handler: apiHandler})

	http.ListenAndServe(":8080", mux)
}
