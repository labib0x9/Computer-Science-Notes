package main

import "log"

func main() {

	log.Println("[Log] error")
	
	ipAddr, port := "127.0.0.1", "8080"
	log.Printf("ip %s failed to start at port %s", ipAddr, port)

}