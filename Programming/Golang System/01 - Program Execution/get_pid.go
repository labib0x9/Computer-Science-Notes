//go:build linux

package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {

	mainPID := os.Getpid()
	mainPGID, _ := syscall.Getpgid(mainPID)

	fmt.Println(mainPID, mainPGID)

	time.Sleep(10 * time.Minute)
}
