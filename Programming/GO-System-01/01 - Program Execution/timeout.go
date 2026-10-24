package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {

	// set timeout for 10sec
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	cmd := exec.CommandContext(ctx, "./sum")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
		return
	}

	fmt.Println("Succesfully terminated")
}
