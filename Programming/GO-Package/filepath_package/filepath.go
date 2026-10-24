package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func getFileType(filename string) {
	switch ext := filepath.Ext(filename); ext {
	case ".txt":
		//
	case ".cpp":
		//
	case ".go":
		//
	case ".py":
		//
	default:
	}
}

func main() {

	path := filepath.Join("golang", "dsa", "algo", "stack.go")
	fmt.Println(path)

	base := filepath.Base(path)
	fmt.Println(base)

	extension := filepath.Ext(base)
	fmt.Println(extension)

	dirName := filepath.Dir(path)
	fmt.Print(dirName)

	// Recursive walk
	filepath.Walk(dirName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		info.IsDir()
		info.Name()
		info.Size()
		return nil
	})

	// filepath.Glob()
}
