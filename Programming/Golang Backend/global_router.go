package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type Topic struct {
	PageUrl          string `json:"page_url"`
	IconUrl          string `json:"icon_url"`
	Name             string `json:"name"`
	ShortDescription string `json:"desp"`
	AHref            string `json:"a_href"`
}

var AllTopics []Topic

func init() {
	linux := Topic{
		PageUrl:          "/linux",
		IconUrl:          "/static/linux-plain.svg",
		Name:             "Linux Environments",
		ShortDescription: "Deploy isolated, ephemeral Linux containers directly from your browser. Practice privilege escalation, system administration, and network forensics securely.",
		AHref:            "Explore Linux Labs →",
	}

	ccode := Topic{
		PageUrl:          "/cprog",
		IconUrl:          "/static/c-original.svg",
		Name:             "C Programming",
		ShortDescription: "Write, compile, and execute strict C code using our WebSocket-powered remote execution engine. Learn memory management, pointers, and safe architectural patterns.",
		AHref:            "Explore C Exercises →",
	}

	AllTopics = append(AllTopics, []Topic{linux, ccode}...)
}

func getAllTopics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.Encode(AllTopics)

	default:
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
}

func getTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(AllTopics[0])
}

func GlobalRouter(mux *http.ServeMux) http.Handler {
	handlerAllReq := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTION, GET")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		} else {
			mux.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(handlerAllReq)
}

func main() {

	// Request (/get/topics) -> GlobalRouter -> Mux -> getTopics

	mux := http.NewServeMux()
	mux.HandleFunc("/getAllTopics", getAllTopics)
	mux.Handle("GET /get/topics", http.HandlerFunc(getTopics))

	globalRouter := GlobalRouter(mux)

	slog.Info("Starting Server at http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", globalRouter))
}
