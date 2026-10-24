package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// HERE IS SOME CONFUSION>> NEED TO CLEARR....

func main() {

	// PID(./sum)= 50662
	// Parent PGID = Child PGID
	// labib@MacBookAir ~ % ps -p 50662 -o pid,ppid,pgid
	// PID  PPID  PGID
	// 50662 50661 50645
	// labib@MacBookAir ~ % ps -p 50661 -o pid,ppid,pgid
	// PID  PPID  PGID
	// 50661 50645 50645

	// killing child pid will kill parent. viceversa.

	cmd := exec.Command("./sum")
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
	} else {
		fmt.Println("Succesfully terminated")
	}

	////////////-----////////////
	// PID(./sum)= 1330
	// Parent PGID != Child PGID
	// labib@MacBookAir ~ % ps -p 1330 -o pid,ppid,pgid
	// PID  PPID  PGID
	// 1330  1323  1330
	// labib@MacBookAir ~ % ps -p 1323 -o pid,ppid,pgid
	// PID  PPID  PGID
	// 1323  1307  1307

	// Killing the child will not effect parent.

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

	err = cmd1.Wait()
	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
	} else {
		fmt.Println("Succesfully terminated")
	}
}
