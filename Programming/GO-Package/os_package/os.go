package main

import (
	"os"
	"time"
)

// Incomplete

func main() {

	os.Chmod("file.txt", 0644)                     // changes file permission
	os.Chtimes("file.txt", time.Now(), time.Now()) // change timstamp
	os.FileMode

	os.Getgroups()
	gid := os.Getegid()
	uid := os.Geteuid()
	os.Chown("file.txt", uid, gid) // change owner/group

	os.Getwd()          // get working dir
	os.Chdir("dirName") // change working dir

	// Environment variable
	os.Setenv("APP_ENV", "dev")
	os.Getenv("APP_ENV")
	os.Unsetenv("APP_ENV")
	os.Environ()
	os.Clearenv()

	os.Symlink("example.txt", "example-link.txt")
	os.Readlink("example-link.txt")

	os.UserHomeDir()
	os.UserCacheDir()

	os.Getpagesize()
	os.Hostname()

	// in-memory stream
	reader, writer, err := os.Pipe()
	p := make([]byte, 512)
	reader.Read(p)
	writer.Write(p)

	os.IsExist(err)
	os.IsNotExist(err)
	os.IsPermission(err)
	os.IsPathSeparator(uint8('/'))

	// Process
	os.Getpid()
	os.ProcAttr{}
	os.StartProcess()
	os.FindProcess()

	// File System Abstractions
	fs := os.DirFS(".")
	file, err := fs.Open("text.txt")
	file.Stat()
	file.Close()

	os.CopyFS()

	os.Interrupt.Signal()
	os.Interrupt.String()

	os.NewSyscallError()

	os.Exit(34)
}
