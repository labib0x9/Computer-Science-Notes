package main

import (
	"fmt"
	"net"
)

func main() {

	localAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("hello"))

	buf := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(buf)

	fmt.Println(buf[:n])
}
