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
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// what to do?
			break
		}

		fmt.Println("Client:", r.RemoteAddr, "MSG TYPE:", msgType, ", MSG:", string(msg))

		switch msgType {
		case 1:
			_ = conn.WriteMessage(msgType, []byte(string(msg)+" -> server"))
		case 2:
			_ = conn.WriteMessage(msgType, []byte("server is working broo"))
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
