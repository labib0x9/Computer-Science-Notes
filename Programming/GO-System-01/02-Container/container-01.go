// //go:build linux

// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// 	"syscall"

// 	"golang.org/x/sys/unix"
// )

// func PrintPidAndMntNS(pid int, tag string) string {
// 	fmt.Println(tag, "PID=", pid)

// 	// mnt namespace id read
// 	mntNSPath := fmt.Sprintf("/proc/%d/ns/mnt", pid)

// 	mntNSOut, err := os.Readlink(mntNSPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	mntNSID := strings.TrimPrefix(mntNSOut, "mnt:[")
// 	mntNSID = strings.TrimSuffix(mntNSID, "]")
// 	fmt.Println(mntNSPath, "=", mntNSID)

// 	return mntNSID
// }

// func pivotRoot(newRoot string) error {
// 	putold := filepath.Join(newRoot, "/.pivot_root")

// 	err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, "")
// 	if err != nil {
// 		return err
// 	}

// 	// create putold directory
// 	err = os.MkdirAll(putold, 0700)
// 	if err != nil {
// 		return err
// 	}

// 	// call pivot_root
// 	err = syscall.PivotRoot(newRoot, putold)
// 	if err != nil {
// 		return err
// 	}

// 	// ensure current working directory is set to new root
// 	err = os.Chdir("/")
// 	if err != nil {
// 		return err
// 	}

// 	//umount putold, which now lives at /.pivot_root
// 	putold = "/.pivot_root"
// 	err = syscall.Unmount(putold, syscall.MNT_DETACH)
// 	if err != nil {
// 		return err
// 	}

// 	// remove putold
// 	err = os.RemoveAll(putold)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // self-re-execution pattern
// func parent() {
// 	// PrintPidAndMntNS(os.Getpid(), "parent")

// 	cmd := exec.Command("/proc/self/exe", "child")
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Cloneflags: syscall.CLONE_NEWNS |
// 			syscall.CLONE_NEWPID |
// 			syscall.CLONE_NEWUTS,
// 		Unshareflags: syscall.CLONE_NEWNS,
// 	}

// 	if err := cmd.Run(); err != nil {
// 		panic(err)
// 	}
// }

// func child() {
// 	syscall.Sethostname([]byte("toor"))

// 	rootfsPath := "/home/labib0x9/rootfs/alpine_tmp"
// 	pivotRoot(rootfsPath)

// 	syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")
// 	syscall.Mount("temp", "/tmp", "tmpfs", 0, "")
// 	syscall.Mount("proc", "/proc", "proc", 0, "")
// 	PrintPidAndMntNS(os.Getpid(), "child")

// 	cmd := exec.Command("sh")
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		panic(err)
// 	}

// 	unix.Unmount("proc", 0)
// 	unix.Unmount("temp", 0)
// }

// // sudo go run main.go <flags>
// func main() {
// 	args := os.Args[1:]

// 	if len(args) == 0 {
// 		fmt.Println("Usage: sudo go run main.go <run|child>")
// 		return
// 	}

// 	if len(args) > 1 && args[1] == "child" {
// 		child()
// 		return
// 	}

// 	switch args[0] {
// 	case "run":
// 		parent()
// 	case "child":
// 		child()
// 	default:
// 		fmt.Println("Usage: sudo go run main.go <run|child>")
// 	}
// }
