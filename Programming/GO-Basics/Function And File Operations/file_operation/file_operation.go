package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

const FilePerm = 0644 // File permission: -rw-r--r--
const DirPerm = 0755  // Directory permission: drwxr-xr-x

// Check if directory exists
// also works for files
func isDirExists(path string) bool {
	info, err := os.Stat(path)  // follows symlink
	// info, err := os.Lstat(path) // doesn't follow symlinks
	if os.IsNotExist(err) {
		return false
	}

	// // File Info
	// info.Name()
	// info.Size()	// bytes
	// info.ModTime()

	// // Detect symlink with os.Lstat()
	// if info.Mode() & os.ModeSymlink != 0 {
	// 	// Symlink
	// }

	return info.IsDir()
}

// Create directory; `nested=true` creates parent directories too
func createDir(path string, nested bool) {
	var err error
	if nested {
		err = os.MkdirAll(path, DirPerm)
	} else {
		err = os.Mkdir(path, DirPerm)
	}
	if err != nil {
		panic(err)
	}
}

// Get sub-directories from the given path
func getDir(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var dirs []string
	for _, entry := range files {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
		// // File Info
		// info, err := entry.Info()
		// if err != nil {
		// 	panic(err)
		// }
		// info.Name()
		// info.Size()
		// info.ModTime()
	}
	return dirs
}

// Remove directory; `nested=true` removes recursively
func removeDir(path string, nested bool) {
	var err error
	if nested {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		panic(err)
	}
}

// Rename directory or file
func renameDir(oldPath string, newPath string) {
	if err := os.Rename(oldPath, newPath); err != nil {
		panic(err)
	}
}

func randomInt() int {
	return rand.Int()
}

// Save data to a file using a temporary file to avoid corruption
func SaveDataToFileUsingTempFile(path string, data []byte) error {
	// temp := fmt.Sprintf("%s.temp.%d", path, randomInt())
	// file, err := os.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, FilePerm) // Creates a file if not exits
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	temp := path + "-tempfile-*.temp"
	file, err := os.CreateTemp(path, temp)

	_, err = file.Write(data)
	if err != nil {
		os.Remove(temp)
		return err
	}

	return os.Rename(temp, path) // Replace original with new
}

// Create empty file (overwrites if exists)
func createFile(filename string) {
	file, err := os.Create(filename) // O_RDWR -> Read, Write
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

func fileFunctions(file *os.File) {

	// reset the file pointer from where read/write starts
	file.Seek(0, io.SeekStart)               // Jump to start of file
	file.Seek(0, io.SeekEnd)                 //Jump to end of file
	file.Seek(10, io.SeekCurrent)            // forward 10bytes from current
	file.Seek(-5, io.SeekEnd)                // backward 5bytes from end
	pos, err := file.Seek(0, io.SeekCurrent) // check current position

	if err != nil || pos != 0 {
		// Do Something
	}

	file.Sync() // forces to flush the in-memory buffer to disk

	// Random access reads
	buf := make([]byte, 5)
	info, _ := os.Stat(file.Name())
	info, _ = file.Stat()
	file.ReadAt(buf, info.Size()-int64(5)) // Read last 5 bytes

	// Start writing from the current position
	// Advances the file pointer
	// Sequential write
	file.Write([]byte("Content"))

	// Start writing from offset position
	// file pointer doesn't change
	file.WriteAt([]byte("hello"), 4) // Overrides len(data) bytes starting at byte 4
}

// Read entire file content using io
func readFileIo(file *os.File) string {
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}

// Read file line-by-line using bufio
func readFileBufio(file *os.File) string {
	reader := bufio.NewReader(file)
	var content strings.Builder
	for {
		str, err := reader.ReadString('\n')
		content.WriteString(str)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	return content.String()
}

// Open file and read with both methods (for demonstration)
func openAndRead(filename string) {
	file, err := os.Open(filename) // O_RDONLY -> Read
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content := readFileIo(file)

	// Reset file pointer before next read
	file.Seek(0, io.SeekStart)
	content = readFileBufio(file)

	fmt.Println(content)
}

// Append content to the end of file
func appendToFile(filename string, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, FilePerm)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		panic(err)
	}
	file.Sync()
}

// Overwrite file content from beginning (truncate)
func overwriteFromBeginning(filename string, content string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, FilePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		panic(err)
	}

	// os.WriteFile(filename, []byte(content), FilePerm)
}

// Remove a file
func removeFile(filename string) {
	if err := os.Remove(filename); err != nil {
		panic(err)
	}
}

// Truncate a file to zero size (clears content)
// shrink the content
func truncateFile(file *os.File) {
	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	// os.Truncate(file.Name(), 0)

	// Clear last n bytes
	n := 10
	info, _ := os.Stat(file.Name())
	currentFileSize := info.Size()
	if currentFileSize-int64(n) < 0 {
		// File doesn't have n bytes
		// Handle case
	}
	if err := file.Truncate(currentFileSize - int64(n)); err != nil {
		// handle error
	}
}

// Copy file from old path to new path
func copyFile(old, new string) {
	oldFile, err := os.Open(old)
	if err != nil {
		panic(err)
	}
	defer oldFile.Close()

	newFile, err := os.Create(new)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, oldFile); err != nil {
		panic(err)
	}
	// file.Sync() flushes the in-memory buffer to disk
	newFile.Sync()
}

func main() {

}
