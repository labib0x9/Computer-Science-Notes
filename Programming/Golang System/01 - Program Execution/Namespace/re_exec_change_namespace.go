//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func PrintPidAndMntNS(pid int, tag string) {
	fmt.Println(tag, "PID=", pid)

	// mnt namespace id read
	mntNSPath := fmt.Sprintf("/proc/%d/ns/mnt", pid)

	mntNSOut, err := os.Readlink(mntNSPath)
	if err != nil {
		panic(err)
	}

	mntNSID := strings.TrimPrefix(mntNSOut, "mnt:[")
	mntNSID = strings.TrimSuffix(mntNSID, "]")
	fmt.Println(mntNSPath, "=", mntNSID)
}

// self-re-execution pattern
func parent() {
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS,
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// parentPID := cmd.Process.Pid
	PrintPidAndMntNS(cmd.Process.Pid, "Parent")

	cmd.Wait()
}

// i have a namespace, which is different from my parent...
func child() {
	cmd := exec.Command("sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	PrintPidAndMntNS(cmd.Process.Pid, "Child")

	cmd.Wait()
}

func main() {

	args := os.Args
	// fmt.Println(args)
	if len(args) > 1 && args[1] == "child" {
		// fmt.Println(args[1])
		child()
		return
	}

	pid := os.Getpid()
	PrintPidAndMntNS(pid, "Main")

	parent()
}
