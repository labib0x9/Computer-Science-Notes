package main

import (
	_"fmt"
	"io"
	"os"
	_"time"
)

func main() {
	files := []string{"a.txt", "b.txt", "c.txt"}
	// done := make(chan struct{})
	error := make(chan error)
	for _, filename := range files {
		// go backup(filename, done)
		go backupError(filename, error)
	}
	for i := 0; i < len(files); i++ {

		// <-done

		if err := <-error; err != nil {
			// Do Something
		}
	}
}

func backup(filename string, done chan struct{}) {
	in, err := os.Open(filename)
	if err != nil {
		done <- struct{}{}
		return
	}
	defer in.Close()

	// os.Create() creates a file with filename, if exist  then overwrites
	out, err := os.Create(filename + ".bak")
	if err != nil {
		done <- struct{}{}
		return
	}
	defer out.Close()

	// Ignoring error
	out.WriteString("BackupStart\n")
	io.Copy(out, in)
	out.WriteString("\nBackupEnd")

	done <- struct{}{}
}

func backupError(filename string, error chan error) {
	// fmt.Println(time.Now())
	in, err := os.Open(filename)
	if err != nil {
		error <- err
		return
	}
	defer in.Close()

	out, err := os.Create(filename + ".bak")
	if err != nil {
		error <- err
		return
	}
	defer out.Close()

	if _, err := out.WriteString("Backup\n"); err != nil {
		error <- err
		return
	}

	if _, err := io.Copy(out, in); err != nil {
		error <- err
		return
	}

	if _, err := out.WriteString("\nBackup"); err != nil {
		error <- err
		return
	}

	error <- nil
}