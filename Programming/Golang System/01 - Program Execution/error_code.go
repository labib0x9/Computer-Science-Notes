package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command("sh", "-c", "pwd")

	// output via terminal / stdio
	// cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// // Parent PID
	// fmt.Println("Parent PID=", os.Getpid())

	if err := cmd.Run(); err != nil {
		fmt.Println("Run error:" + err.Error())
		return
	}

	
	// Handle error code
	err := cmd.Wait()
	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Wait error:" + err.Error())
		fmt.Println("Exit Code:", errExit.ExitCode())
		return
	}
}
