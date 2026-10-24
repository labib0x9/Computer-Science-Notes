package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1) // Channel to receive OS signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // Subscribe to SIGINT (Ctrl+C) and SIGTERM

	go func() {
		sig := <-sigChan
		fmt.Println("Received: ", sig)
		os.Exit(0)
	}()

	select {} // Block forever
}
