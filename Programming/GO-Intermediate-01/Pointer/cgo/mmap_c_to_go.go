package main

/*
#include<stdio.h>
#include<sys/mman.h>
#include<fcntl.h>
#include<unistd.h>
#include<stdint.h>

struct Header {
    uint32_t MagicByte;    // Magic-byte = 0xa1b2c3d4 = 2712847316
    uint16_t VersionMajor; // File format version
	uint16_t VersionMinor; // File format version
	int32_t   TimeZone;     // TimeZone offset, normally 0
	uint32_t TimeStamp;    // TimStamp accuracy, normally 0
	uint32_t SnapLen;      // Max capture length per packet
	uint32_t LinkType;     // Link-Layer type
};

int HeaderSize() {
    return sizeof(struct Header);
}

char* OpenFile(int fd, int length) {
    char *addr;
    addr = mmap(
        NULL,
        length,
        PROT_READ,
        MAP_PRIVATE,
        fd,
        0
    );

    if (addr == MAP_FAILED) {
        return NULL;
    }

	printf("C addr: %p\n", &addr[0]);

    return addr;
}

int CloseFile(char *addr, int length) {
    return munmap(addr, length);
}

struct Header* GetPCAPHeader(char* addr) {
    struct Header *hdr = (struct Header*)addr;
    return hdr;
}
*/
import "C"

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

func main() {

	// filename := "http.cap"
	filename := "dns.cap"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// fmt.Println(unsafe.Sizeof(pcapHeader{}))
	x := C.HeaderSize()
	fmt.Println(int(x))

	addr := C.OpenFile(C.int(file.Fd()), C.int(info.Size()))
	if addr == nil {
		panic("Mapped failed")
	}

	hdr := C.GetPCAPHeader(addr)
	_ = hdr

	fmt.Println(hdr.MagicByte)
	fmt.Println("Go addr: ", &hdr)
	fmt.Println(unsafe.Pointer(&hdr))

	fmt.Println(reflect.TypeOf(hdr))

	if errno := C.CloseFile(addr, C.int(info.Size())); errno == -1 {
		panic("Unmapped failed")
	}
}
