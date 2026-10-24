package main

// printf "How TCP\r\n" | nc -w 1 127.0.0.1 8080

import (
	"fmt"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}
		fmt.Println(string(buf[:n]))
	}
}
