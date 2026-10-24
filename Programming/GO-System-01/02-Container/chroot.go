//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func logger(reason string, err error) {
	if err == nil {
		return
	}
	panic(fmt.Sprintf("[%s] %v", reason, err))
}

// self-re-execution pattern
func parent() {
	exePath, err := os.Executable()
	logger("parent exe path", err)

	cmd := exec.Command(exePath, "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | // NEWUSER for rootless containersyscall.CLONE_NEWPID |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
		// Unshareflags: syscall.CLONE_NEWNS,
		Pdeathsig: syscall.SIGKILL,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
		GidMappingsEnableSetgroups: false,
	}

	logger("self execute", cmd.Run())
}

func child() {

	// if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
	// 	panic("I DONT WANNA FORWARD::" + err.Error())
	// }

	// make mount namespace private
	logger("mount private", syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""))

	rootfsPath := "/home/labib0x9/rootfs/alpine_tmp"

	logger("sethostname", syscall.Sethostname([]byte("toor")))

	// CHROOT - ROOT FOLDER
	logger("chroot", syscall.Chroot(rootfsPath))
	logger("chdir", syscall.Chdir("/"))

	logger("/tmp mount", syscall.Mount("tmpfs", "/tmp", "tmpfs", 0, ""))
	logger("/proc mount", syscall.Mount("proc", "/proc", "proc", 0, ""))

	defer logger("/proc umonunt", syscall.Unmount("/proc", 0))
	defer logger("/tmp unmount", syscall.Unmount("/temp", 0))

	cmd := exec.Command("sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger("", cmd.Run())
}

// go run main.go <flags>
func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: go run main.go <run|child>")
		return
	}

	switch args[0] {
	case "run":
		parent()
	case "child":
		child()
	default:
		fmt.Println("Usage: go run main.go <run|child>")
	}
}
