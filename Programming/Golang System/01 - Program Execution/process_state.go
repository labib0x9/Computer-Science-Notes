//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	cmd1 := exec.Command("./sum")
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	if err := cmd1.Start(); err != nil {
		panic(err)
	}

	fmt.Println("PID=", cmd1.Process.Pid)
	fmt.Println("Process State is (Nil) =", cmd1.ProcessState == nil)

	err := cmd1.Wait()
	_ = err
	// if errExit, ok := err.(*exec.ExitError); ok {
	// 	fmt.Println("Run Error=", err.Error())
	// 	fmt.Println("Exit Error=", errExit.ExitCode())
	// } else {
	// 	fmt.Println("Succesfully terminated")
	// }

	fmt.Println("Process State is (Nil) =", cmd1.ProcessState == nil)

	pState := cmd1.ProcessState
	fmt.Println(pState.String())
	fmt.Println(pState.ExitCode())
	fmt.Println(pState.Exited())
	fmt.Println(pState.Pid())
	fmt.Println(pState.Sys())
	fmt.Println(pState.SysUsage())
	fmt.Println(pState.SystemTime())
	fmt.Println(pState.UserTime())
}
