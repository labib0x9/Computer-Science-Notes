//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

// Explain what is happening here, why the execution stops ??
// With - Ignore SITTOU Signal
// Without ignore..

func main() {

	// signal.Ignore(syscall.SIGTTOU)

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

	childPID := cmd1.Process.Pid
	fmt.Println("Child Process ID=", childPID)

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	childPGID, _ := unix.Getpgid(childPID)

	if err := unix.IoctlSetPointerInt(int(tty.Fd()), unix.TIOCSPGRP, childPGID); err != nil {
		panic(err)
	}

	// _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&childPGID)))
	// if errno != 0 {
	// 	panic(fmt.Sprintf("Tcsetpgrp failed: %v", errno))
	// }

	err = cmd1.Wait()

	// parentPGID := unix.Getpgrp()
	// unix.IoctlSetPointerInt(int(tty.Fd()), unix.TIOCSPGRP, parentPGID)

	parentPGID := syscall.Getpgrp()
	// _, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&parentPGID)))
	// if errno != 0 {
	// 	panic(fmt.Sprintf("Tcsetpgrp failed: %v", errno))
	// }
	if err1 := unix.IoctlSetPointerInt(int(tty.Fd()), unix.TIOCSPGRP, parentPGID); err1 != nil {
		panic(err1)
	}

	if errExit, ok := err.(*exec.ExitError); ok {
		fmt.Println("Run Error=", err.Error())
		fmt.Println("Exit Error=", errExit.ExitCode())
	} else {
		fmt.Println("Succesfully terminated")
	}

	// _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&childPID)))
	// if errno != 0 {
	// 	panic(fmt.Sprintf("Tcsetpgrp failed: %v", errno))
	// }
}
