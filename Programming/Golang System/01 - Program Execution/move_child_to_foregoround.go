//go:build linux

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {

	signal.Ignore(syscall.SIGTTOU)

	cmd1 := exec.Command("./sum")
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	cmd1.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var err error
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	if err := cmd1.Start(); err != nil {
		panic(err)
	}

	childPID := cmd1.Process.Pid
	fmt.Println("Child Process ID=", childPID)
	childPGID, _ := unix.Getpgid(childPID)

	// move child process to foreground...
	unix.IoctlSetPointerInt(
		int(tty.Fd()),
		unix.TIOCSPGRP,
		childPGID,
	)

	// check ps
	err = cmd1.Wait()

	// move parent to foreground
	parentPGID := syscall.Getpgrp()
	unix.IoctlSetPointerInt(
		int(tty.Fd()),
		unix.TIOCSPGRP,
		parentPGID,
	)

	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
	} else {
		fmt.Println("Succesfully terminated")
	}

	// When this code doesn't executes ???
	// check ps
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
