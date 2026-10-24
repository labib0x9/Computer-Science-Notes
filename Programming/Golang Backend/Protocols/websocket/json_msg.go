package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type A struct {
	Id int `json:"req_id"`
}

func (a *A) Increment() {
	a.Id++
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, http.Header{})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 40*time.Second)
	defer cancel()

	// cacellation - Method A
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetReadDeadline(deadline)
	}

	for {
		var ID A
		if err := conn.ReadJSON(&ID); err != nil {
			break
		}

		ID.Increment()

		if err := conn.WriteJSON(&ID); err != nil {
			break
		}
	}
}

func main() {

	mux := http.NewServeMux()
	mux.Handle(
		"GET /echo",
		http.HandlerFunc(handleConn),
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println(server.ListenAndServe())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	fmt.Println(server.Shutdown(ctx))
}
