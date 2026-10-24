package main

import (
	"fmt"
	"os"
	"os/exec"
)

/*
	ps -p 49233 -o ppid,stat,comm,rss,vsz,%cpu,%mem

	lsof -p 49233
*/

func main() {

	// // set timeout for 10sec
	// ctx, cancel := context.WithTimeout(
	// 	context.Background(),
	// 	10*time.Second,
	// )
	// defer cancel()

	cmd := exec.Command("./server")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	fmt.Println("PID=", cmd.Process.Pid)

	err := cmd.Wait()
	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
		return
	}

	fmt.Println("Succesfully terminated")
}
