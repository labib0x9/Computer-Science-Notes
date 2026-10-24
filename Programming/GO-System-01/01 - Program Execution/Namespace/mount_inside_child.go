//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// NOT WORKING...

func main() {

	// For OS Thread Locking
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Read parent mnt namespace
	parentNs, err := os.Open("/proc/self/ns/mnt")
	if err != nil {
		panic(err)
	}
	defer parentNs.Close()

	// cmd...
	cmd1 := exec.Command("sh")
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	cmd1.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS,
	}

	if err := cmd1.Start(); err != nil {
		panic(err)
	}

	// while child running, on host/parent, setup....
	childPID := cmd1.Process.Pid
	fmt.Println("Child PID=", childPID)

	// mnt namespace id read
	mntNSPath := fmt.Sprintf("/proc/%d/ns/mnt", childPID)
	fmt.Println("ls -l", mntNSPath)

	mntNSOut, err := os.Readlink(mntNSPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(mntNSPath, "=", mntNSOut)
	mntNSID := strings.TrimPrefix(mntNSOut, "mnt:[")
	mntNSID = strings.TrimSuffix(mntNSID, "]")
	fmt.Println("mnt =", mntNSID)

	/******
		But We have to mount inside the child..., how ??
		method-1: SETNS SYSCALL
	******/

	// We have to enter child's namespace
	childNS, err := os.Open(mntNSPath)
	if err != nil {
		panic(err)
	}
	defer childNS.Close()

	// # panic: invalid argument
	if err := unix.Setns(int(childNS.Fd()), unix.CLONE_NEWNS); err != nil {
		panic(err)
	}

	// make a directory in /tmp folder -> /tmp/tmp_mntid
	tmpPath := "/tmp/tmp_" + mntNSID
	if err := os.Mkdir(tmpPath, 0677); err != nil {
		panic(err.Error())
	}

	// now we are inside child's namespace...
	unix.Mount("tmpfs", tmpPath, "tmpfs", 0, "")
	unix.Mount("", tmpPath, "", unix.MS_PRIVATE|unix.MS_REC, "")

	// Rollback to parent namespace...
	if err := unix.Setns(int(parentNs.Fd()), unix.CLONE_NEWNS); err != nil {
		panic(err)
	}
	/**/

	err = cmd1.Wait()
	// _ = err

	// Cleanup
	fmt.Println("Cleanup")
	if err := os.RemoveAll(tmpPath); err != nil {
		panic(err)
	}
}
