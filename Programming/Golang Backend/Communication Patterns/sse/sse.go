package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/mem"
)

type Status struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Perc  float64 `json:"perc"`
}

func memEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	flusher := http.NewResponseController(w)
	counter := 0

	for {
		select {
		case <-r.Context().Done():
			fmt.Println("CLient disconnected")
			return
		case <-ticker.C:
			counter++
			m, err := mem.VirtualMemory()
			if err != nil {
				return
			}

			if counter%2 == 1 {
				data := fmt.Sprintf("Total: %d, Used: %d, Perc: %.2f", m.Total, m.Used, m.UsedPercent)
				fmt.Fprintf(w, "event:mem\ndata:%s\n\n", data)
			} else {
				status := Status{Total: m.Total, Used: m.Used, Perc: m.UsedPercent}
				data, err := json.Marshal(status)
				if err != nil {
					return
				}
				fmt.Fprintf(w, "event:json\ndata:%s\n\n", data)
			}
			
			if err := flusher.Flush(); err != nil {
				return
			}
		case <-ctx.Done():
			fmt.Fprintf(w, "event: done\ndata: timeout\n\n")
			flusher.Flush()
			fmt.Println("Timeout event operation")
			return
		}
	}
}

func main() {

	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(http.Dir("./"))) // for index.html

	handler.HandleFunc("/status/mem", memEventHandler)

	http.ListenAndServe(":8080", handler)

}
