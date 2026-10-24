package main

import (
	"fmt"
	"net"
	"time"
)

// Connect to server using netcat
// nc 127.0.0.1 8080

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func SendTime(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Connected as " + conn.RemoteAddr().String() + "\n"))
	for {
		fmt.Fprintf(conn, fmt.Sprintf("(%s) > %s\n", conn.RemoteAddr().String(), GetCurrentTime()))
		time.Sleep(1 * time.Second)
	}
}

func main() {

	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()	// connect one user
		if err != nil {
			continue
		}
		go SendTime(conn)	// send user to goroutine
	}
}
