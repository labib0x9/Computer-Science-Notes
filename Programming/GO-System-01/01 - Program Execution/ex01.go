//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Explain what is happening here, why the exeecution hangs ??

func main() {

	cmd1 := exec.Command("./sum")
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	cmd1.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd1.Start(); err != nil {
		panic(err)
	}

	fmt.Println("PID=", cmd1.Process.Pid)

	err := cmd1.Wait()
	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
	} else {
		fmt.Println("Succesfully terminated")
	}
}
