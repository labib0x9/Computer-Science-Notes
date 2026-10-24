package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command("./sum")

	// output via terminal / stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Parent PID
	fmt.Println("Parent PID=", os.Getpid())

	// // No way to get the child pid
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println("Run error:" + err.Error())
	// 	return
	// }

	if err := cmd.Start(); err != nil {
		fmt.Println("Start error:" + err.Error())
		return
	}

	// Child PID
	fmt.Println("Child PID=", cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait error:" + err.Error())
		return
	}
}
