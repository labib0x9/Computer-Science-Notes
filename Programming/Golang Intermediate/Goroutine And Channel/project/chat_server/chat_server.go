package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

/**

	TCP Chat Server

	What am I trying to implement ??
	Handle multiple client,
	Broadcast one message to all client using channel
**/

type Message struct {
	msg  string
	user string
}

type Client struct {
	user string
	conn net.Conn
}

var (
	ActiveClient map[Client]struct{} // it's a shared map
	mu          sync.Mutex
)

func ClientHandler(client Client, broadMsg chan Message) {
	defer func() {
		mu.Lock()
		delete(ActiveClient, client)
		mu.Unlock()
		client.conn.Close()
	}()
	fmt.Fprintf(client.conn, "Connected as %s\n", client.user)
	reader := bufio.NewScanner(client.conn)
	for {
		var msg string
		if reader.Scan() {
			msg = reader.Text()
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		broadMsg <- Message{msg: msg, user: client.user}
	}
}

func broadcast(broadMsg chan Message) {
	for msg := range broadMsg {
		mu.Lock()
		for user := range ActiveClient {
			if user.user == msg.user {
				continue
			}
			fmt.Fprintf(user.conn, "(%s) >> %s\n", msg.user, msg.msg)
		}
		mu.Unlock()
	}
}

func main() {
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listner.Close()

	broadMsg := make(chan Message)
	go broadcast(broadMsg)

	if ActiveClient == nil {
		ActiveClient = make(map[Client]struct{})
	}

	for {
		conn, err := listner.Accept() // connect one user
		if err != nil {
			continue
		}

		client := Client{
			conn: conn,
			user: conn.RemoteAddr().String(),
		}

		mu.Lock()
		ActiveClient[client] = struct{}{}
		mu.Unlock()

		go ClientHandler(client, broadMsg) // send user to goroutine
	}
}
