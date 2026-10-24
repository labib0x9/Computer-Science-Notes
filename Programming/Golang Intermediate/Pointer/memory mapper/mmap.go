package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {

	// open file
	filename := "test.txt"
	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// file status
	var stat syscall.Stat_t
	if err := syscall.Fstat(int(f.Fd()), &stat); err != nil {
		panic(err)
	}

	fInfo, err := f.Stat() // calls syscall.Fstat() internally ...
	if err != nil {
		panic(err)
	}
	fmt.Println(stat.Size, fInfo.Size())

	// map area
	data, err := syscall.Mmap(
		int(f.Fd()),
		0,
		int(stat.Size),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		panic(err)
	}

	// free area
	defer func() {
		if err := syscall.Munmap(data); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Prev: ", string(data))

	threeByte := make([]byte, 3)

	// read first three bytes
	copy(threeByte, data[0:3])
	fmt.Println(threeByte, string(threeByte))

	copy(threeByte, data[3:6])
	fmt.Println(threeByte, string(threeByte))

	// write
	copy(data[0:3], []byte("hi"))
	copy(data[3:6], []byte("val"))

	fmt.Println("Updated: ", string(data))

	// flush -> file
	_, _, errno := syscall.Syscall(
		syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(stat.Size),
		uintptr(syscall.MS_SYNC),
	)

	if errno != 0 {
		panic(errno)
	}
}